package logger_test

import (
	"testing"

	"github.com/reddotpay/logger"
	"github.com/stretchr/testify/assert"
)

func TestLogger_MaskCard(t *testing.T) {
	s := `{"number":"4111111111111111"}`
	assert.Equal(t, "{\"number\":\"************1111\"}", logger.MaskCard(s))
}

func TestLogger_MaskCard_EmptyString(t *testing.T) {
	s := ""
	assert.Equal(t, "", logger.MaskCard(s))
}
