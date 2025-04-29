package log

import (
	"context"
	"maps"
	"strconv"
	"time"
)

type Span struct {
	ctx       context.Context
	logger    Logger
	operation string
	fields    map[string]any
	start     time.Time
	config    Config
}

// StartLogSpan
func StartLogSpan(ctx context.Context, logger Logger, operation string, opts ...Option) *Span {
	config := NewConfig(opts...)

	fields := make(map[string]any)
	maps.Copy(fields, config.fields)

	return &Span{
		ctx:       ctx,
		logger:    logger,
		operation: operation,
		fields:    fields,
		start:     time.Now(),
		config:    config,
	}
}

// WithField добавляет поле в span
func (s *Span) WithField(key string, value any) *Span {
	s.fields[key] = value
	return s
}

// WithFields добавляет несколько полей в span
func (s *Span) WithFields(fields map[string]any) *Span {
	maps.Copy(s.fields, fields)
	return s
}

// Finish завершает span и логирует результат
func (s *Span) Finish(results ...any) {
	s.finish(results...)

	s.logger.Info(s.operation+" completed", s.fields)
}

// FinishWithError завершает span с ошибкой
func (s *Span) FinishWithError(err error, results ...any) {
	s.finish(results...)

	s.fields["error"] = err.Error()

	if le, ok := err.(LoggableError); ok {
		maps.Copy(s.fields, le.Fields())
	}

	s.logger.Error(s.operation+" failed", s.fields)
}

func (s *Span) finish(results ...any) {
	if s.ctx != nil {
		contextFields := extractFieldsFromContext(s.ctx)
		for k, v := range contextFields {
			if _, exists := s.fields[k]; !exists {
				s.fields[k] = v
			}
		}
	}

	duration := time.Since(s.start)
	s.fields["duration"] = duration.String()

	if s.config.logResult {
		s.addResultsToFields(results)
	}
}

func (s *Span) addResultsToFields(results []any) {
	if len(results) == 0 {
		return
	}

	if len(results) == 1 {
		s.fields["result"] = results[0]
		return
	}

	for i, res := range results {
		fieldName := "result"
		if i > 0 {
			fieldName += "_" + strconv.Itoa(i)
		}
		s.fields[fieldName] = res
	}
}
