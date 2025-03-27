package fiber

import (
	"github.com/gofiber/fiber/v3"

	"github.com/go-mosaic/runtime/transport"
)

// FiberToMiddleware преобразует middleware fiber в универсальный Middleware
func FiberToMiddleware(fiberMiddleware fiber.Handler) transport.Middleware {
	return func(next transport.Handler) transport.Handler {
		return func(req transport.Request, resp transport.Response) error {
			if err := fiberMiddleware(req.(*FiberRequest).ctx); err != nil {
				return err
			}

			return next(req, resp)
		}
	}
}

// FiberTransport реализация Transport с использованием fiber
type FiberTransport struct {
	app     *fiber.App
	adapter *FiberAdapter
}

// NewFiberTransport создает новый экземпляр FiberTransport
func NewFiberTransport(app *fiber.App) *FiberTransport {
	return &FiberTransport{
		app: app,
		adapter: &FiberAdapter{
			writeResponse: transport.DefaultWriteResponse,
			readData:      transport.DefaultReadData,
		},
	}
}

func (t *FiberTransport) AddRoute(method, path string, handler transport.Handler, middlewares ...transport.Middleware) {
	wrappedHandler := handler
	for _, mw := range middlewares {
		wrappedHandler = mw(wrappedHandler)
	}

	t.app.Add([]string{method}, path, t.adapter.AdaptHandler(wrappedHandler))
}

func (t *FiberTransport) Use(middlewares ...transport.Middleware) {
	for _, mw := range middlewares {
		t.app.Use(func(c fiber.Ctx) error {
			wrappedHandler := mw(func(req transport.Request, resp transport.Response) error {
				return c.Next()
			})

			return wrappedHandler(&FiberRequest{ctx: c}, &FiberResponse{ctx: c})
		})
	}
}
