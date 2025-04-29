package log

import (
	"context"
	"time"
)

// MetricSpan расширяет базовый Span для сбора метрик
type MetricSpan struct {
	*Span
	collector MetricsCollector
	operation string
	start     time.Time
}

// StartMetricSpan начинает новый span с метриками
func StartMetricSpan(
	ctx context.Context,
	logger Logger,
	operation string,
	collector MetricsCollector,
	opts ...Option,
) *MetricSpan {
	if collector != nil {
		collector.RecordCall(operation)
	}

	return &MetricSpan{
		Span:      StartLogSpan(ctx, logger, operation, opts...),
		collector: collector,
		operation: operation,
		start:     time.Now(),
	}
}

// Finish завершает span и записывает метрики успешного выполнения
func (s *MetricSpan) Finish(results ...any) {
	duration := time.Since(s.start)
	if s.collector != nil {
		s.collector.RecordSuccess(s.operation, duration)
	}

	s.Span.Finish(results...)
}

// FinishWithError завершает span и записывает метрики ошибки
func (s *MetricSpan) FinishWithError(err error, results ...any) {
	duration := time.Since(s.start)
	if s.collector != nil {
		s.collector.RecordError(s.operation, duration)
	}

	s.Span.FinishWithError(err, results...)
}
