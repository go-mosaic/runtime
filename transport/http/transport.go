package http

import (
	"net/http"

	"github.com/go-mosaic/runtime/transport"
)

// HTTPToMiddleware преобразует middleware net/http в универсальный Middleware
func HTTPToMiddleware(httpMiddleware func(http.Handler) http.Handler) transport.Middleware {
	return func(next transport.Handler) transport.Handler {
		return func(req transport.Request, resp transport.Response) error {
			httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_ = next(req, resp)
			})

			wrappedHandler := httpMiddleware(httpHandler)
			wrappedHandler.ServeHTTP(resp.(*HTTPResponse).w, req.(*HTTPRequest).req)

			return nil
		}
	}
}

type HTTPTransport struct {
	router      *http.ServeMux
	middlewares []transport.Middleware
	adapter     *HTTPAdapter
}

func NewHTTPTransport() *HTTPTransport {
	return &HTTPTransport{
		router: http.NewServeMux(),
		adapter: &HTTPAdapter{
			writeResponse: transport.DefaultWriteResponse,
			readData:      transport.DefaultReadData,
		},
	}
}

func (t *HTTPTransport) AddRoute(method, path string, handler transport.Handler, middlewares ...transport.Middleware) {
	wrappedHandler := handler

	for _, mw := range append(t.middlewares, middlewares...) {
		wrappedHandler = mw(wrappedHandler)
	}

	t.router.HandleFunc(path, t.adapter.AdaptHandler(wrappedHandler))
}

func (t *HTTPTransport) Use(middlewares ...transport.Middleware) {
	t.middlewares = append(t.middlewares, middlewares...)
}
