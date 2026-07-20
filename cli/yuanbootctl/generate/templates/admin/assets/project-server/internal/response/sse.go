package response

import "net/http"

type SSEResult struct {
	Handler func(http.ResponseWriter) error
}

func (r SSEResult) Render(w http.ResponseWriter) error {
	return r.Handler(w)
}

func (r SSEResult) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "text/event-stream; charset=utf-8")
	}
}
