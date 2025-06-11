package logger

type NoopLogger struct{}

func (n NoopLogger) Debug(message string, fields ...map[string]interface{}) {
}

func (n NoopLogger) Error(message string, fields ...map[string]interface{}) {
}

func (n NoopLogger) Info(message string, fields ...map[string]interface{}) {
}

func (n NoopLogger) Trace(message string, fields ...map[string]interface{}) {
}

func (n NoopLogger) Warn(message string, fields ...map[string]interface{}) {
}

var _ Logger = &NoopLogger{}
