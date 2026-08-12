package web

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/valyala/fasthttp"
)

const (
	responseModeDone = iota + 1
	responseModeFlushed
	responseModeHijacked
	responseModePanicked
)

var responseStreamBufferPool = sync.Pool{
	New: func() any {
		buffer := make([]byte, 32*1024)
		return &buffer
	},
}

func NewFastHTTPHandler(h http.Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		request, err := convertFastHTTPRequest(ctx)
		if err != nil {
			ctx.Logger().Printf("cannot parse requestURI %q: %s", request.RequestURI, err)
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			return
		}

		writer := newNetHTTPResponseWriter(ctx)
		go func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					ctx.Logger().Printf("panic in net/http handler: %v", recovered)
					writer.signalMode(responseModePanicked)
				} else {
					writer.signalMode(responseModeDone)
				}
				_ = writer.Close()
			}()
			h.ServeHTTP(writer, request.WithContext(ctx))
		}()

		switch <-writer.modeCh {
		case responseModeDone:
			copyResponseMetadata(ctx, writer, false)
			if body := writer.consumeBufferedBody(); len(body) > 0 {
				ctx.Response.SetBody(body)
			}
			releaseNetHTTPResponseWriter(writer)
		case responseModeFlushed:
			copyResponseMetadata(ctx, writer, true)
			ctx.SetBodyStreamWriter(func(bufferedWriter *bufio.Writer) {
				defer releaseNetHTTPResponseWriter(writer)
				if body := writer.consumeBufferedBody(); len(body) > 0 {
					if _, err := bufferedWriter.Write(body); err != nil {
						return
					}
					if err := bufferedWriter.Flush(); err != nil {
						return
					}
				}

				buffer := responseStreamBufferPool.Get().(*[]byte)
				defer responseStreamBufferPool.Put(buffer)
				for {
					n, err := writer.pipeReader.Read(*buffer)
					if n > 0 {
						if _, writeErr := bufferedWriter.Write((*buffer)[:n]); writeErr != nil {
							return
						}
						if flushErr := bufferedWriter.Flush(); flushErr != nil {
							return
						}
					}
					if err != nil {
						return
					}
				}
			})
			close(writer.streamReady)
		case responseModeHijacked:
			return
		case responseModePanicked:
			panic("net/http handler panicked")
		}
	}
}

func convertFastHTTPRequest(ctx *fasthttp.RequestCtx) (*http.Request, error) {
	body := ctx.PostBody()
	request := &http.Request{
		Method:        string(ctx.Method()),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		RequestURI:    string(ctx.RequestURI()),
		ContentLength: int64(len(body)),
		Host:          string(ctx.Host()),
		RemoteAddr:    ctx.RemoteAddr().String(),
		Header:        make(http.Header),
		Body:          &netHTTPBody{body},
	}
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		name := string(key)
		if name == "Transfer-Encoding" {
			request.TransferEncoding = append(request.TransferEncoding, string(value))
			return
		}
		request.Header.Set(name, string(value))
	})
	requestURL, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		return request, err
	}
	request.URL = requestURL
	return request, nil
}

func copyResponseMetadata(ctx *fasthttp.RequestCtx, writer *NetHTTPResponseWriter, streaming bool) {
	ctx.SetStatusCode(writer.StatusCode())
	haveContentType := false
	for key, values := range writer.Header() {
		if streaming && key == fasthttp.HeaderContentLength {
			continue
		}
		if key == fasthttp.HeaderContentType {
			haveContentType = true
		}
		for _, value := range values {
			ctx.Response.Header.Add(key, value)
		}
	}
	if !haveContentType {
		writer.mu.Lock()
		defer writer.mu.Unlock()
		if len(writer.body) > 0 {
			length := min(len(writer.body), 512)
			ctx.Response.Header.Set(fasthttp.HeaderContentType, http.DetectContentType(writer.body[:length]))
		}
	}
}

type netHTTPBody struct {
	b []byte
}

func (r *netHTTPBody) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}

func (r *netHTTPBody) Close() error {
	r.b = r.b[:0]
	return nil
}

type NetHTTPResponseWriter struct {
	Ctx *fasthttp.RequestCtx

	statusCode atomic.Int64
	h          http.Header
	mu         sync.Mutex
	body       []byte
	bodyPool   *[]byte

	pipeReader  *io.PipeReader
	pipeWriter  *io.PipeWriter
	modeCh      chan int
	streamReady chan struct{}
	flushOnce   sync.Once
	closeOnce   sync.Once
	hijackOnce  sync.Once
}

func newNetHTTPResponseWriter(ctx *fasthttp.RequestCtx) *NetHTTPResponseWriter {
	pipeReader, pipeWriter := io.Pipe()
	return &NetHTTPResponseWriter{
		Ctx:         ctx,
		h:           make(http.Header),
		pipeReader:  pipeReader,
		pipeWriter:  pipeWriter,
		modeCh:      make(chan int, 1),
		streamReady: make(chan struct{}),
	}
}

func releaseNetHTTPResponseWriter(writer *NetHTTPResponseWriter) {
	_ = writer.Close()
	if writer.bodyPool != nil {
		responseStreamBufferPool.Put(writer.bodyPool)
		writer.bodyPool = nil
	}
}

func (w *NetHTTPResponseWriter) StatusCode() int {
	if statusCode := int(w.statusCode.Load()); statusCode != 0 {
		return statusCode
	}
	return http.StatusOK
}

func (w *NetHTTPResponseWriter) Header() http.Header {
	return w.h
}

func (w *NetHTTPResponseWriter) WriteHeader(statusCode int) {
	if statusCode < 100 || statusCode > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", statusCode))
	}
	w.statusCode.CompareAndSwap(0, int64(statusCode))
}

func (w *NetHTTPResponseWriter) Write(p []byte) (int, error) {
	select {
	case <-w.streamReady:
		return w.pipeWriter.Write(p)
	default:
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.body == nil {
		w.bodyPool = responseStreamBufferPool.Get().(*[]byte)
		w.body = (*w.bodyPool)[:0]
	}
	w.body = append(w.body, p...)
	return len(p), nil
}

func (w *NetHTTPResponseWriter) Flush() {
	w.flushOnce.Do(func() {
		w.signalMode(responseModeFlushed)
	})
	<-w.streamReady
}

func (w *NetHTTPResponseWriter) MarkHijacked() {
	w.hijackOnce.Do(func() {
		w.signalMode(responseModeHijacked)
	})
}

func (w *NetHTTPResponseWriter) Close() error {
	w.closeOnce.Do(func() {
		_ = w.pipeWriter.Close()
		_ = w.pipeReader.Close()
	})
	return nil
}

func (w *NetHTTPResponseWriter) signalMode(mode int) {
	select {
	case w.modeCh <- mode:
	default:
	}
}

func (w *NetHTTPResponseWriter) consumeBufferedBody() []byte {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.body) == 0 {
		return nil
	}
	body := w.body
	w.body = nil
	return body
}
