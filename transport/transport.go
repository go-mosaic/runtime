package transport

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/a-h/templ"
	"github.com/aohorodnyk/mimeheader"
)

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

type Cookie struct {
	Name        string
	Value       string
	Path        string    // опциональный параметр
	Domain      string    // опциональный параметр
	Expires     time.Time // опциональный параметр
	MaxAge      int
	Secure      bool
	HttpOnly    bool
	SameSite    SameSite
	Partitioned bool
}

type Form interface {
	FormValue(name string) string
	FormFile(name string) (multipart.File, *multipart.FileHeader, error)
}

// Transport интерфейс для HTTP транспорта
type Transport interface {
	AddRoute(method, path string, handler Handler, middlewares ...Middleware)
	Use(middlewares ...Middleware)
}

// Request универсальный интерфейс для HTTP-запроса
type Request interface {
	Context() context.Context
	WithContext(ctx context.Context) Request
	Method() string
	Path() string
	Body() io.ReadCloser
	Header(key string) string
	Queries() url.Values
	PathValue(name string) string
	MultipartForm(maxMemory int64) (Form, error)
	URLEncodedForm() (url.Values, error)
	ReadData(data any) error
	Cookie(name string) (string, error)
}

// Response универсальный интерфейс для HTTP-ответа
type Response interface {
	SetStatusCode(code int)
	SetHeader(key, value string)
	SetBody(body []byte, statusCode int) int
	WriteHeader(statusCode int)
	Write([]byte) (int, error)
	WriteData(req Request, data any)
}

type ByteReader interface {
	ReadBytes(mimeType string) ([]byte, error)
}

// ReadData универсальный интерфейс для чтения данных из тела HTTP-запроса
type ReadData func(req Request, data any) error

// WriteResponse универсальный интерфейс для записи данных в HTTP-ответ
type WriteResponse func(req Request, resp Response, data any)

// Handler универсальный интерфейс для обработчика запросов
type Handler func(req Request, resp Response) error

// Middleware универсальный интерфейс для middleware
type Middleware func(next Handler) Handler

type MultipartFormWrapper struct {
	form *multipart.Form
}

// StatusCoder интерфейс для получения статус кода
type StatusCoder interface {
	StatusCode() int
}

// Headerer интерфейс для получения заголовков HTTP-ответа
type Headerer interface {
	Headers() http.Header
}

// Errorer интерфейс для получения текста ошибки
type Errorer interface {
	Error() string
}

func (w *MultipartFormWrapper) FormValue(name string) string {
	if vs := w.form.Value[name]; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (w *MultipartFormWrapper) FormFile(name string) (multipart.File, *multipart.FileHeader, error) {
	if fhs := w.form.File[name]; len(fhs) > 0 {
		f, err := fhs[0].Open()
		return f, fhs[0], err
	}

	return nil, nil, http.ErrMissingFile
}

func MultipartFormWrap(form *multipart.Form) Form {
	return &MultipartFormWrapper{form: form}
}

func DefaultWriteResponse(req Request, resp Response, data any) {
	mimeType := headerAccept(req)
	resp.SetHeader("Content-Type", mimeType)

	// Если нет данных то возвращаем пустой ответ с кодом 204
	if data == nil || data == struct{}{} {
		handleNonDataResponse(resp)
	}

	// Устанавливаем статус код по умолчанию
	statusCode := determineStatusCode(data)

	if headerer, ok := data.(Headerer); ok {
		setHeaders(resp, headerer.Headers())
	}

	// Обрабатываем ответ в зависимости от MIME типа
	switch mimeType {
	case "application/json":
		handleJSONResponse(resp, data, statusCode)
	case "application/xml":
		handleXMLResponse(resp, data, statusCode)
	default:
		handleNonJSONResponse(req, resp, data, mimeType, statusCode)
	}
}

// determineStatusCode определяет HTTP статус код на основе данных
func determineStatusCode(data any) int {
	if data == nil {
		return http.StatusNoContent
	}

	if sc, ok := data.(StatusCoder); ok {
		return sc.StatusCode()
	}

	if _, ok := data.(Errorer); ok {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// setHeaders устанавливает заголовки ответа
func setHeaders(resp Response, headers map[string][]string) {
	for k, values := range headers {
		for _, v := range values {
			resp.SetHeader(k, v)
		}
	}
}

// handleJSONResponse обрабатывает JSON ответ
func handleJSONResponse(resp Response, data any, statusCode int) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		writeFailure(resp, err)
		return
	}
	resp.SetBody(dataBytes, statusCode)
}

// handleJSONResponse обрабатывает XML ответ
func handleXMLResponse(resp Response, data any, statusCode int) {
	dataBytes, err := xml.Marshal(data)
	if err != nil {
		writeFailure(resp, err)
		return
	}
	resp.SetBody(dataBytes, statusCode)
}

// handleNonJSONResponse обрабатывает не-JSON ответы
func handleNonJSONResponse(req Request, resp Response, data any, mimeType string, statusCode int) {
	switch t := data.(type) {
	case ByteReader:
		handleByteReaderResponse(resp, t, mimeType, statusCode)
	case templ.Component:
		handleTemplComponentResponse(req, resp, t, statusCode)
	case Errorer:
		handleErrorerResponse(resp, t, statusCode)
	default:
		resp.SetHeader("Content-Type", "text/plain")
		resp.WriteHeader(http.StatusNotAcceptable)
	}
}

// handleNonDataResponse обрабатываем пустой ответ
func handleNonDataResponse(resp Response) {
	resp.WriteHeader(http.StatusNoContent)
}

func handleByteReaderResponse(resp Response, br ByteReader, mimeType string, statusCode int) {
	bytes, err := br.ReadBytes(mimeType)
	if err != nil {
		writeFailure(resp, err)
		return
	}
	resp.SetHeader("Content-Type", mimeType)
	resp.SetBody(bytes, statusCode)
}

func handleTemplComponentResponse(req Request, resp Response, component templ.Component, statusCode int) {
	resp.WriteHeader(statusCode)

	if err := component.Render(req.Context(), resp); err != nil {
		writeFailure(resp, err)
	}
}

func handleErrorerResponse(resp Response, err Errorer, statusCode int) {
	resp.SetHeader("Content-Type", "text/plain")
	resp.SetBody([]byte(err.Error()), statusCode)
}

func writeFailure(resp Response, err error) {
	resp.SetHeader("content-type", "text/plain")
	resp.SetBody([]byte(err.Error()), http.StatusInternalServerError)
}

func DefaultReadData(req Request, data any) error {
	if err := json.NewDecoder(req.Body()).Decode(data); err != nil {
		return err
	}

	return nil
}

func headerAccept(req Request) string {
	acceptHeader := req.Header("Accept")
	ah := mimeheader.ParseAcceptHeader(acceptHeader)
	_, mimeType, _ := ah.Negotiate([]string{"application/json", "application/xml", "text/plain", "text/html"}, "application/json")

	return mimeType
}
