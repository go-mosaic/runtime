package log

import (
	"context"
	"maps"
)

// ContextExtractor функция для извлечения значений из контекста
type ContextExtractor func(ctx context.Context) map[string]any

var (
	contextExtractors []ContextExtractor
)

// AddContextExtractor добавляет функцию для извлечения данных из контекста
func AddContextExtractor(extractor ContextExtractor) {
	contextExtractors = append(contextExtractors, extractor)
}

// extractFieldsFromContext извлекает все поля из контекста
func extractFieldsFromContext(ctx context.Context) map[string]any {
	fields := make(map[string]any)
	for _, extractor := range contextExtractors {
		maps.Copy(fields, extractor(ctx))
	}
	return fields
}
