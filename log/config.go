package log

import "maps"

// Config конфигурация
type Config struct {
	logStart   bool
	logResult  bool
	logErrors  bool
	skipFields []string
	fields     map[string]any
}

// Option тип для функциональных опций
type Option func(*Config)

// WithLogStart включает/выключает логирование начала выполнения
func WithLogStart(enabled bool) Option {
	return func(c *Config) {
		c.logStart = enabled
	}
}

// WithLogResult включает/выключает логирование результатов
func WithLogResult(enabled bool) Option {
	return func(c *Config) {
		c.logResult = enabled
	}
}

// WithLogErrors включает/выключает логирование ошибок
func WithLogErrors(enabled bool) Option {
	return func(c *Config) {
		c.logErrors = enabled
	}
}

// WithSkipFields задает поля которые не нужно логировать
func WithSkipFields(fields []string) Option {
	return func(c *Config) {
		c.skipFields = fields
	}
}

// WithFields добавляет дополнительные поля в лог
func WithFields(fields map[string]any) Option {
	return func(c *Config) {
		if c.fields == nil {
			c.fields = make(map[string]any)
		}

		maps.Copy(c.fields, fields)
	}
}

// NewConfig создает новую конфигурацию с опциями
func NewConfig(opts ...Option) Config {
	config := Config{
		logStart:  true,
		logResult: true,
		logErrors: true,
	}

	for _, applyOpt := range opts {
		applyOpt(&config)
	}

	return config
}
