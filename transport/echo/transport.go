package echo

import (
	"github.com/labstack/echo/v4"

	"github.com/go-mosaic/runtime/transport"
)

// EchoToMiddleware преобразует middleware echo в универсальный Middleware
func EchoToMiddleware(echoMiddleware echo.MiddlewareFunc) transport.Middleware {
	return func(next transport.Handler) transport.Handler {
		return func(req transport.Request, resp transport.Response) error {
			echoHandler := func(c echo.Context) error {
				return next(req, resp)
			}

			wrappedHandler := echoMiddleware(echoHandler)

			return wrappedHandler(req.(*EchoRequest).ctx)
		}
	}
}

type EchoTransport struct {
	router  *echo.Echo
	adapter *EchoAdapter
}

func NewEchoTransport(router *echo.Echo) *EchoTransport {
	return &EchoTransport{
		router: router,
		adapter: &EchoAdapter{
			writeResponse: transport.DefaultWriteResponse,
			readData:      transport.DefaultReadData,
		},
	}
}

func (t *EchoTransport) AddRoute(method, path string, handler transport.Handler, middlewares ...transport.Middleware) {
	wrappedHandler := handler
	for _, mw := range middlewares {
		wrappedHandler = mw(wrappedHandler)
	}

	t.router.Add(method, path, t.adapter.AdaptHandler(wrappedHandler))
}

func (t *EchoTransport) Use(middlewares ...transport.Middleware) {
	for _, mw := range middlewares {
		t.router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				wrappedHandler := mw(func(req transport.Request, resp transport.Response) error {
					return next(c)
				})

				if err := wrappedHandler(&EchoRequest{ctx: c}, &EchoResponse{ctx: c}); err != nil {
					return err
				}

				return nil
			}
		})
	}
}
