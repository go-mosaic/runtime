package fiber

import (
	"context"
	"io"
	"net/url"

	"github.com/gofiber/fiber/v3"

	"github.com/go-mosaic/runtime/transport"
)

type fiberReadCloser struct {
	data  []byte
	index int
}

func (f *fiberReadCloser) Close() error {
	return nil
}

func (f *fiberReadCloser) Read(p []byte) (n int, err error) {
	if f.index >= len(f.data) {
		return 0, io.EOF
	}
	n = copy(p, f.data[f.index:])
	f.index += n

	return n, nil
}

// FiberRequest адаптер для fiber.Ctx
type FiberRequest struct {
	ctx      fiber.Ctx
	readData transport.ReadData
}

func (r *FiberRequest) WithContext(ctx context.Context) transport.Request {
	// TODO: необходимо найти реализацию
	return r
}

func (r *FiberRequest) Context() context.Context {
	return r.ctx.RequestCtx()
}

func (r *FiberRequest) Method() string {
	return r.ctx.Method()
}

func (r *FiberRequest) Path() string {
	return r.ctx.Path()
}

func (r *FiberRequest) Body() io.ReadCloser {
	return &fiberReadCloser{data: r.ctx.Body()}
}

func (r *FiberRequest) Header(key string) string {
	return r.ctx.Get(key)
}

func (r *FiberRequest) Queries() url.Values {
	m := make(url.Values, r.ctx.RequestCtx().QueryArgs().Len())
	r.ctx.RequestCtx().QueryArgs().VisitAll(func(key, value []byte) {
		m[string(key)] = []string{string(value)}
	})

	return m
}

func (r *FiberRequest) PathValue(name string) string {
	return r.ctx.Params(name)
}

func (r *FiberRequest) MultipartForm(maxMemory int64) (transport.Form, error) {
	form, err := r.ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	return transport.MultipartFormWrap(form), nil
}

func (r *FiberRequest) URLEncodedForm() (url.Values, error) {
	return url.Values{}, nil
}

func (r *FiberRequest) ReadData(data any) error {
	return r.readData(r, data)
}

// FiberResponse адаптер для fiber.Ctx
type FiberResponse struct {
	ctx           fiber.Ctx
	writeResponse transport.WriteResponse
}

func (r *FiberResponse) SetStatusCode(code int) {
	r.ctx.Status(code)
}

func (r *FiberResponse) SetHeader(key, value string) {
	r.ctx.Set(key, value)
}

func (r *FiberResponse) WriteData(req transport.Request, data any) {
	r.writeResponse(req, r, data)
}

func (r *FiberResponse) Write(body []byte) (int, error) {
	return r.ctx.Write(body)
}

func (r *FiberResponse) SetBody(body []byte, statusCode int) int {
	r.WriteHeader(statusCode)
	n, _ := r.Write(body)

	return n
}

func (r *FiberResponse) WriteHeader(statusCode int) {
	r.ctx.Status(statusCode)
}

// FiberAdapter адаптер для fiber
type FiberAdapter struct {
	writeResponse transport.WriteResponse
	readData      transport.ReadData
}

func (a *FiberAdapter) AdaptHandler(handler transport.Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := &FiberRequest{ctx: c, readData: a.readData}
		resp := &FiberResponse{ctx: c, writeResponse: a.writeResponse}
		if err := handler(req, resp); err != nil {
			a.writeResponse(req, resp, err)
		}
		return nil
	}
}
