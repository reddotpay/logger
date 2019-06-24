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
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return s
	}

	for k, v := range m {
		switch strings.ToLower(k) {
		case "cvv", "securitycode":
			if value, ok := v.(string); ok {
				m[k] = mask(value, len(value))
			}
		case "number", "cardnumber", "cardnum":
			if value, ok := v.(string); ok {
				m[k] = mask(value, 4)
			}
		case "card":
			if value, ok := v.(map[string]interface{}); ok {
				if number, ok := value["number"].(string); ok {
					value["number"] = mask(number, 4)
					m[k] = value
				}

				if securityCode, ok := value["securityCode"].(string); ok {
					value["securityCode"] = mask(securityCode, len(securityCode))
					m[k] = value
				}
			}
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
