package logger

type Logger interface {
	Debug(message string, fields ...map[string]interface{})
	Error(message string, fields ...map[string]interface{})
	Info(message string, fields ...map[string]interface{})
	Trace(message string, fields ...map[string]interface{})
	Warn(message string, fields ...map[string]interface{})
}
