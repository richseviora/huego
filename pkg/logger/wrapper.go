package logger

import (
	"github.com/hashicorp/go-hclog"
)

type LogWrapper struct {
	Logger hclog.Logger
}

func (l LogWrapper) Debug(message string, fields ...map[string]interface{}) {
	slice := mapFieldsToSlice(fields)
	l.Logger.Debug(message, slice...)
}

func (l LogWrapper) Error(message string, fields ...map[string]interface{}) {
	l.Logger.Error(message, mapFieldsToSlice(fields)...)
}

func (l LogWrapper) Info(message string, fields ...map[string]interface{}) {
	slice := mapFieldsToSlice(fields)
	l.Logger.Info(message, slice...)
}

func (l LogWrapper) Trace(message string, fields ...map[string]interface{}) {
	l.Logger.Trace(message, mapFieldsToSlice(fields)...)
}

func (l LogWrapper) Warn(message string, fields ...map[string]interface{}) {
	l.Logger.Warn(message, mapFieldsToSlice(fields)...)
}

var (
	_ Logger = &LogWrapper{}
)

func mapFieldsToSlice(s []map[string]interface{}) []interface{} {
	values := make([]interface{}, 0)
	for _, v := range s {
		values = append(values, mapToKeyValueSlice(v)...)
	}
	return values
}

func mapToKeyValueSlice(m map[string]interface{}) []interface{} {
	values := make([]interface{}, 0)
	for k, v := range m {
		values = append(values, k, v)
	}
	return values
}

func NewLogger() Logger {
	return &LogWrapper{
		Logger: hclog.Default(),
	}
}
