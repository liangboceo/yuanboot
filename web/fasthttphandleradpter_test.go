package web

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestNewFastHTTPHandlerStreamsFlushedResponse(t *testing.T) {
	firstChunkWritten := make(chan struct{})
	finishHandler := make(chan struct{})
	handler := NewFastHTTPHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/event-stream")
		_, _ = writer.Write([]byte("data: first\n\n"))
		writer.(http.Flusher).Flush()
		close(firstChunkWritten)
		<-finishHandler
		_, _ = writer.Write([]byte("data: second\n\n"))
		writer.(http.Flusher).Flush()
	}))

	listener := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{Handler: handler}
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() {
		_ = server.Shutdown()
		_ = listener.Close()
	})

	connection, err := listener.Dial()
	require.NoError(t, err)
	defer connection.Close()
	require.NoError(t, connection.SetDeadline(time.Now().Add(3*time.Second)))
	_, err = fmt.Fprint(connection, "GET /stream HTTP/1.1\r\nHost: test\r\nConnection: close\r\n\r\n")
	require.NoError(t, err)

	reader := bufio.NewReader(connection)
	response, err := http.ReadResponse(reader, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "text/event-stream", response.Header.Get("Content-Type"))

	select {
	case <-firstChunkWritten:
	case <-time.After(time.Second):
		t.Fatal("handler did not flush first chunk")
	}
	line, err := bufio.NewReader(response.Body).ReadString('\n')
	require.NoError(t, err)
	require.Equal(t, "data: first\n", line)
	close(finishHandler)
}

func TestNewFastHTTPHandlerKeepsBufferedResponses(t *testing.T) {
	handler := NewFastHTTPHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		_, _ = writer.Write([]byte(`{"ok":true}`))
	}))

	listener := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{Handler: handler}
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() {
		_ = server.Shutdown()
		_ = listener.Close()
	})

	connection, err := listener.Dial()
	require.NoError(t, err)
	defer connection.Close()
	require.NoError(t, connection.SetDeadline(time.Now().Add(3*time.Second)))
	_, err = fmt.Fprint(connection, "GET /resource HTTP/1.1\r\nHost: test\r\nConnection: close\r\n\r\n")
	require.NoError(t, err)

	response, err := http.ReadResponse(bufio.NewReader(connection), nil)
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, response.StatusCode)
	require.True(t, strings.HasPrefix(response.Header.Get("Content-Type"), "application/json"))
	require.JSONEq(t, `{"ok":true}`, string(body))
}

func TestNewFastHTTPHandlerSupportsWebSocketUpgrade(t *testing.T) {
	upgrader := websocket.FastHTTPUpgrader{CheckOrigin: func(*fasthttp.RequestCtx) bool { return true }}
	handler := NewFastHTTPHandler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		responseWriter := writer.(*NetHTTPResponseWriter)
		err := upgrader.Upgrade(responseWriter.Ctx, func(connection *websocket.Conn) {
			messageType, message, readErr := connection.ReadMessage()
			if readErr == nil {
				_ = connection.WriteMessage(messageType, message)
			}
			_ = connection.Close()
		})
		if err == nil {
			responseWriter.MarkHijacked()
		}
	}))

	listener := fasthttputil.NewInmemoryListener()
	server := &fasthttp.Server{Handler: handler}
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(func() {
		_ = server.Shutdown()
		_ = listener.Close()
	})

	dialer := websocket.Dialer{
		NetDial: func(network, addr string) (net.Conn, error) {
			return listener.Dial()
		},
		HandshakeTimeout: 3 * time.Second,
	}
	connection, response, err := dialer.Dial("ws://test/socket", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusSwitchingProtocols, response.StatusCode)
	defer connection.Close()
	require.NoError(t, connection.WriteMessage(websocket.TextMessage, []byte("ping")))
	messageType, message, err := connection.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, messageType)
	require.Equal(t, []byte("ping"), message)
}
