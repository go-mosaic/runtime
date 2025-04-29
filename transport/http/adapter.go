package http

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/go-mosaic/runtime/transport"
)

// HTTPRequest адаптер для net/http.Request
type HTTPRequest struct {
	req      *http.Request
	readData transport.ReadData
}

func (r *HTTPRequest) WithContext(ctx context.Context) transport.Request {
	return &HTTPRequest{req: r.req.WithContext(ctx), readData: r.readData}
}

func (r *HTTPRequest) Context() context.Context {
	return r.req.Context()
}

func (r *HTTPRequest) Method() string {
	return r.req.Method
}

func (r *HTTPRequest) Path() string {
	return r.req.URL.Path
}

func (r *HTTPRequest) Body() io.ReadCloser {
	return r.req.Body
}

func (r *HTTPRequest) Header(key string) string {
	return r.req.Header.Get(key)
}

func (r *HTTPRequest) Queries() url.Values {
	return r.req.URL.Query()
}

func (r *HTTPRequest) PathValue(name string) string {
	return r.req.PathValue(name)
}

func (r *HTTPRequest) MultipartForm(maxMemory int64) (transport.Form, error) {
	err := r.req.ParseMultipartForm(maxMemory)
	return transport.MultipartFormWrap(r.req.MultipartForm), err
}

func (r *HTTPRequest) URLEncodedForm() (url.Values, error) {
	if err := r.req.ParseForm(); err != nil {
		return nil, err
	}

	return r.req.Form, nil
}

func (r *HTTPRequest) ReadData(data any) error {
	return r.readData(r, data)
}

func (r *HTTPRequest) SetCookie(c transport.Cookie) error {
	r.req.AddCookie(&http.Cookie{
		Name:        c.Name,
		Value:       c.Value,
		Path:        c.Path,
		Domain:      c.Domain,
		Expires:     c.Expires,
		MaxAge:      c.MaxAge,
		Secure:      c.Secure,
		HttpOnly:    c.HttpOnly,
		SameSite:    http.SameSite(c.SameSite),
		Partitioned: c.Partitioned,
	})

	return nil
}

func (r *HTTPRequest) Cookie(name string) (string, error) {
	c, err := r.req.Cookie(name)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

// HTTPResponse адаптер для net/http.ResponseWriter
type HTTPResponse struct {
	w             http.ResponseWriter
	writeResponse transport.WriteResponse
}

func (r *HTTPResponse) SetStatusCode(code int) {
	r.w.WriteHeader(code)
}

func (r *HTTPResponse) SetHeader(key, value string) {
	r.w.Header().Set(key, value)
}

func (r *HTTPResponse) Write(body []byte) (int, error) {
	return r.w.Write(body)
}

func (r *HTTPResponse) WriteData(req transport.Request, data any) {
	r.writeResponse(req, r, data)
}

func (r *HTTPResponse) SetBody(body []byte, statusCode int) int {
	r.WriteHeader(statusCode)
	n, _ := r.Write(body)

	return n
}

func (r *HTTPResponse) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

// HTTPAdapter адаптер для net/http
type HTTPAdapter struct {
	writeResponse transport.WriteResponse
	readData      transport.ReadData
}

func (a *HTTPAdapter) AdaptHandler(handler transport.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &HTTPRequest{req: r, readData: a.readData}
		resp := &HTTPResponse{w: w, writeResponse: a.writeResponse}
		if err := handler(req, resp); err != nil {
			a.writeResponse(req, resp, err)
		}
	}
}
