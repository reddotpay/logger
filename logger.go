package logger

import (
	"encoding/json"
	"fmt"
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
	m := map[string]string{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}

	for k, v := range m {
		switch strings.ToLower(k) {
		case "cvv", "securitycode":
			m[k] = mask(v, len(v))
		case "number", "cardnumber", "cardnum":
			m[k] = mask(v, 4)
		}
	}

	b, _ := json.Marshal(m)

	return string(b)
}

func mask(str string, size int) string {
	if len(str) == size {
		return strings.Repeat("*", size)
	}

	return fmt.Sprintf("%s%s", strings.Repeat("*", len(str)-size), str[len(str)-size:])
}
