package logger

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
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

// MaskCard masks card number and cvv if exists
func MaskCard(s string) string {
	m := map[string]interface{}{}

	// Check if string is JSON and mask card
	if err := json.Unmarshal([]byte(s), &m); err == nil {
		for k, v := range m {
			switch strings.ToLower(k) {
			case "cvv", "securitycode":
				if value, ok := v.(string); ok {
					m[k] = mask(value, len(value))
				}
			case "number", "cardnumber", "cardnum", "cardno", "accountnumber", "card_no":
				if value, ok := v.(string); ok {
					m[k] = mask(value, 4)
				}
			case "card":
				if value, ok := v.(map[string]interface{}); ok {
					var (
						b, _ = json.Marshal(value)
						m2   = map[string]interface{}{}
					)

					_ = json.Unmarshal([]byte(MaskCard(string(b))), &m2)
					m[k] = m2
				}
			}
		}

		b, _ := json.Marshal(m)

		return string(b)
	}

	// Check if string is URL encoded and does not contain `<`, and mask card
	if values, err := url.ParseQuery(s); !strings.Contains(s, "<") && strings.Contains(s, "=") && err == nil {
		newValues := url.Values{}
		for k, v := range values {
			switch strings.ToLower(k) {
			case "cvv", "securitycode":
				newValues[k] = []string{mask(v[0], len(v[0]))}
			case "number", "cardnumber", "cardnum", "cardno", "accountnumber", "card_no":
				newValues[k] = []string{mask(v[0], 4)}
			default:
				newValues[k] = v
			}
		}

		var (
			buf  strings.Builder
			keys = make([]string, 0, len(newValues))
		)
		for k := range newValues {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			vs := newValues[k]
			for _, v := range vs {
				if buf.Len() > 0 {
					buf.WriteByte('&')
				}
				buf.WriteString(k)
				buf.WriteByte('=')
				buf.WriteString(v)
			}
		}

		return buf.String()
	}

	// Check if string is XML and mask card
	r := regexp.MustCompile(`(?i)<(number|cardnumber|cardnum|cardno|accountnumber)>(\d+)<\/(number|cardnumber|cardnum|cardno|accountnumber)>`)
	if m := r.FindStringSubmatch(s); len(m) == 4 {
		s = r.ReplaceAllString(s, fmt.Sprintf("<%s>%s</%s>", m[1], mask(m[2], 4), m[3]))
	}

	r = regexp.MustCompile(`(?i)<(cvv|securitycode|cvNumber)>(\d{3,4})<\/(cvv|securitycode|cvNumber)>`)
	if m := r.FindStringSubmatch(s); len(m) == 4 {
		s = r.ReplaceAllString(s, fmt.Sprintf("<%s>%s</%s>", m[1], mask(m[2], len(m[2])), m[3]))
	}

	return s
}

func mask(str string, size int) string {
	if len(str) == size {
		return strings.Repeat("*", size)
	}

	return fmt.Sprintf("%s%s", strings.Repeat("*", len(str)-size), str[len(str)-size:])
}
