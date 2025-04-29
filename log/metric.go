package log

import "time"

// MetricsCollector интерфейс для сбора метрик
type MetricsCollector interface {
	// RecordSuccess записывает успешный вызов операции с указанием времени выполнения
	RecordSuccess(operation string, duration time.Duration)

	// RecordError записывает ошибку операции с указанием времени выполнения
	RecordError(operation string, duration time.Duration)

	// RecordCall записывает вызов операции
	RecordCall(operation string)
}
