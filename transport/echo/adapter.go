package echo

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/go-mosaic/runtime/transport"
)

// EchoRequest адаптер для echo.Context
type EchoRequest struct {
	ctx      echo.Context
	readData transport.ReadData
}

func (r *EchoRequest) WithContext(ctx context.Context) transport.Request {
	newCtx := r.ctx.Echo().NewContext(
		r.ctx.Request().WithContext(ctx),
		r.ctx.Response(),
	)

	return &EchoRequest{ctx: newCtx, readData: r.readData}
}

func (r *EchoRequest) Context() context.Context {
	return r.ctx.Request().Context()
}

func (r *EchoRequest) Method() string {
	return r.ctx.Request().Method
}

func (r *EchoRequest) Path() string {
	return r.ctx.Request().URL.Path
}

func (r *EchoRequest) Body() io.ReadCloser {
	return r.ctx.Request().Body
}

func (r *EchoRequest) Header(key string) string {
	return r.ctx.Request().Header.Get(key)
}

func (r *EchoRequest) Queries() url.Values {
	return r.ctx.Request().URL.Query()
}

func (r *EchoRequest) PathValue(name string) string {
	return r.ctx.Param(name)
}

func (r *EchoRequest) MultipartForm(maxMemory int64) (transport.Form, error) {
	err := r.ctx.Request().ParseMultipartForm(maxMemory)

	return transport.MultipartFormWrap(r.ctx.Request().MultipartForm), err
}

func (r *EchoRequest) URLEncodedForm() (url.Values, error) {
	return r.ctx.FormParams()
}

func (r *EchoRequest) ReadData(data any) error {
	return r.readData(r, data)
}

func (r *EchoRequest) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	r.ctx.SetCookie(&http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (r *EchoRequest) Cookie(name string) (string, error) {
	c, err := r.ctx.Cookie(name)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

// EchoResponse адаптер для echo.Context
type EchoResponse struct {
	ctx           echo.Context
	writeResponse transport.WriteResponse
}

func (r *EchoResponse) SetStatusCode(code int) {
	r.ctx.Response().WriteHeader(code)
}

func (r *EchoResponse) SetHeader(key, value string) {
	r.ctx.Response().Header().Set(key, value)
}

func (r *EchoResponse) WriteData(req transport.Request, data any) {
	r.writeResponse(req, r, data)
}

func (r *EchoResponse) Write(body []byte) (int, error) {
	return r.ctx.Response().Write(body)
}

func (r *EchoResponse) SetBody(body []byte, statusCode int) int {
	r.WriteHeader(statusCode)
	n, _ := r.Write(body)

	return n
}

func (r *EchoResponse) WriteHeader(statusCode int) {
	r.ctx.Response().Writer.WriteHeader(statusCode)
}

// EchoAdapter адаптер для echo
type EchoAdapter struct {
	writeResponse transport.WriteResponse
	readData      transport.ReadData
}

func (a *EchoAdapter) AdaptHandler(handler transport.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &EchoRequest{ctx: c, readData: a.readData}
		resp := &EchoResponse{ctx: c, writeResponse: a.writeResponse}
		if err := handler(req, resp); err != nil {
			a.writeResponse(req, resp, err)
			return nil
		}
		return nil
	}
}
