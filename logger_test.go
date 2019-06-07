package logger_test

import (
	"testing"

	"github.com/reddotpay/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger_MaskNumber(t *testing.T) {
	s := `{"number":"4111111111111111"}`
	assert.Equal(t, "{\"number\":\"************1111\"}", logger.MaskCard(s))
}

func TestLogger_MaskSecurityCode(t *testing.T) {
	s := `{"securityCode":"123"}`
	assert.Equal(t, "{\"securityCode\":\"***\"}", logger.MaskCard(s))
}

func TestLogger_MaskCVV(t *testing.T) {
	s := `{"cvv":"123"}`
	assert.Equal(t, "{\"cvv\":\"***\"}", logger.MaskCard(s))
}

func TestLogger_MaskCard(t *testing.T) {
	s := `{"number":"4111111111111111","cvv":"123"}`
	assert.Equal(t, "{\"cvv\":\"***\",\"number\":\"************1111\"}", logger.MaskCard(s))
}

func TestLogger_MaskCard_EmptyString(t *testing.T) {
	s := ""
	assert.Equal(t, "", logger.MaskCard(s))
}

func TestLogger_MaskCard_NoMatch(t *testing.T) {
	s := `{"hello":"world","foo":"bar"}`
	assert.Equal(t, "{\"foo\":\"bar\",\"hello\":\"world\"}", logger.MaskCard(s))
}
