package chi

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/go-mosaic/runtime/transport"
)

// ChiRequest адаптер для net/http.Request
type ChiRequest struct {
	req      *http.Request
	readData transport.ReadData
}

func (r *ChiRequest) WithContext(ctx context.Context) transport.Request {
	return &ChiRequest{req: r.req.WithContext(ctx), readData: r.readData}
}

func (r *ChiRequest) Context() context.Context {
	return r.req.Context()
}

func (r *ChiRequest) Method() string {
	return r.req.Method
}

func (r *ChiRequest) Path() string {
	return r.req.URL.Path
}

func (r *ChiRequest) Body() io.ReadCloser {
	return r.req.Body
}

func (r *ChiRequest) Header(key string) string {
	return r.req.Header.Get(key)
}

func (r *ChiRequest) Queries() url.Values {
	return r.req.URL.Query()
}

func (r *ChiRequest) PathValue(name string) string {
	return chi.URLParam(r.req, name)
}

func (r *ChiRequest) MultipartForm(maxMemory int64) (transport.Form, error) {
	err := r.req.ParseMultipartForm(maxMemory)
	return transport.MultipartFormWrap(r.req.MultipartForm), err
}

func (r *ChiRequest) URLEncodedForm() (url.Values, error) {
	if err := r.req.ParseForm(); err != nil {
		return nil, err
	}

	return r.req.Form, nil
}

func (r *ChiRequest) ReadData(data any) error {
	return r.readData(r, data)
}

// ChiResponse адаптер для net/http.ResponseWriter
type ChiResponse struct {
	w             http.ResponseWriter
	writeResponse transport.WriteResponse
}

func (r *ChiResponse) SetStatusCode(code int) {
	r.w.WriteHeader(code)
}

func (r *ChiResponse) SetHeader(key, value string) {
	r.w.Header().Set(key, value)
}

func (r *ChiResponse) WriteData(req transport.Request, data any) {
	r.writeResponse(req, r, data)
}

func (r *ChiResponse) Write(body []byte) (int, error) {
	return r.w.Write(body)
}

func (r *ChiResponse) SetBody(body []byte, statusCode int) int {
	r.WriteHeader(statusCode)
	n, _ := r.Write(body)

	return n
}

func (r *ChiResponse) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

// ChiAdapter адаптер для net/http
type ChiAdapter struct {
	writeResponse transport.WriteResponse
	readData      transport.ReadData
}

func (a *ChiAdapter) AdaptHandler(handler transport.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &ChiRequest{req: r, readData: a.readData}
		resp := &ChiResponse{w: w, writeResponse: a.writeResponse}
		if err := handler(req, resp); err != nil {
			a.writeResponse(req, resp, err)
		}
	}
}
