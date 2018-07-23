package logger

import (
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Client contains logger interface
var Client Logger

// Logger contains necessary methods of zap
type Logger interface {
	Info(key string, fields ...zapcore.Field)
	Error(key string, fields ...zapcore.Field)
	Sync() error
}

// New initialises a new zap logger
func New() *zap.Logger {
	zapLogger, _ := zap.NewProduction()
	return zapLogger
}

// MaskCard masks card number if exists
func MaskCard(s string) string {
	r := regexp.MustCompile(`([0-9]{16,})`)
	number := r.FindString(s)

	if number != "" {
		l := len(number)
		maskedNumber := fmt.Sprintf("%s%s", strings.Repeat("*", l-4), number[l-4:])
		s = r.ReplaceAllString(s, maskedNumber)
	}

	return s
}
