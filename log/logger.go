package log

// Logger интерфейс для абстракции логирования
type Logger interface {
	Debug(msg string, fields map[string]any)
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

// LoggableError интерфейс для ошибок с дополнительной информацией
type LoggableError interface {
	error
	Level() string
	Fields() map[string]any
}
