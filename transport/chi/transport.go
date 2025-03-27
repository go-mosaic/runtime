package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-mosaic/runtime/transport"
)

// ChiToMiddleware преобразует middleware chi в transport.Middleware
func ChiToMiddleware(chiMiddleware func(http.Handler) http.Handler) transport.Middleware {
	return func(next transport.Handler) transport.Handler {
		return func(req transport.Request, resp transport.Response) error {
			chiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = next(req, resp)
			})

			wrappedHandler := chiMiddleware(chiHandler)
			wrappedHandler.ServeHTTP(resp.(*ChiResponse).w, req.(*ChiRequest).req)

			return nil
		}
	}
}

type ChiTransport struct {
	router chi.Router
	adaper *ChiAdapter
}

func NewChiTransport(router chi.Router) *ChiTransport {
	return &ChiTransport{
		router: router,
		adaper: &ChiAdapter{
			writeResponse: transport.DefaultWriteResponse,
			readData:      transport.DefaultReadData,
		},
	}
}

func (t *ChiTransport) AddRoute(method, path string, handler transport.Handler, middlewares ...transport.Middleware) {
	wrappedHandler := handler
	for _, mw := range middlewares {
		wrappedHandler = mw(wrappedHandler)
	}
	t.router.MethodFunc(method, path, t.adaper.AdaptHandler(wrappedHandler))
}

func (t *ChiTransport) Use(middlewares ...transport.Middleware) {
	for _, mw := range middlewares {
		t.router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wrappedHandler := mw(func(req transport.Request, resp transport.Response) error {
					next.ServeHTTP(w, r)
					return nil
				})

				_ = wrappedHandler(&ChiRequest{req: r}, &ChiResponse{w: w})
			})
		})
	}
}
