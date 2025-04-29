package log

import (
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func FromSlog(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{logger: logger}
}

func (s *SlogLogger) Debug(msg string, fields map[string]any) {
	s.logger.Debug(msg, toSlogArgs(fields)...)
}

func (s *SlogLogger) Info(msg string, fields map[string]any) {
	s.logger.Info(msg, toSlogArgs(fields)...)
}

func (s *SlogLogger) Warn(msg string, fields map[string]any) {
	s.logger.Warn(msg, toSlogArgs(fields)...)
}

func (s *SlogLogger) Error(msg string, fields map[string]any) {
	s.logger.Error(msg, toSlogArgs(fields)...)
}

func toSlogArgs(fields map[string]any) []any {
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}
