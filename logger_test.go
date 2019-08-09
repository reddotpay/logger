package logger_test

import (
	"testing"

	"github.com/reddotpay/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLogger_New(t *testing.T) {
	l := logger.New()
	assert := assert.New(t)
	assert.IsType(&zap.Logger{}, l)
}

func TestLogger_NoMatch(t *testing.T) {
	s := `helloworld`
	assert.Equal(t, "helloworld", logger.MaskCard(s))
}

func TestLogger_JSON_MaskNumber(t *testing.T) {
	s := `{"number":"4111111111111111"}`
	assert.Equal(t, "{\"number\":\"************1111\"}", logger.MaskCard(s))
}

func TestLogger_JSON_MaskSecurityCode(t *testing.T) {
	s := `{"securityCode":"123"}`
	assert.Equal(t, "{\"securityCode\":\"***\"}", logger.MaskCard(s))
}

func TestLogger_JSON_MaskCVV(t *testing.T) {
	s := `{"cvv":"123"}`
	assert.Equal(t, "{\"cvv\":\"***\"}", logger.MaskCard(s))
}

func TestLogger_JSON_MaskNestedCard(t *testing.T) {
	s := `{"card":{"number":"4111111111111111","securityCode":"123"}}`
	assert.Equal(t, "{\"card\":{\"number\":\"************1111\",\"securityCode\":\"***\"}}", logger.MaskCard(s))
}

func TestLogger_JSON_MaskCard(t *testing.T) {
	s := `{"number":"4111111111111111","cvv":"123","card":{"number":"4111111111111111"}}`
	assert.Equal(t, "{\"card\":{\"number\":\"************1111\"},\"cvv\":\"***\",\"number\":\"************1111\"}", logger.MaskCard(s))
}

func TestLogger_JSON_MaskCard_EmptyString(t *testing.T) {
	s := ""
	assert.Equal(t, "", logger.MaskCard(s))
}

func TestLogger_JSON_MaskCard_NoMatch(t *testing.T) {
	s := `{"hello":"world","foo":"bar"}`
	assert.Equal(t, "{\"foo\":\"bar\",\"hello\":\"world\"}", logger.MaskCard(s))
}

func TestLogger_URL_MaskNumber(t *testing.T) {
	s := `number=4111111111111111`
	assert.Equal(t, "number=************1111", logger.MaskCard(s))
}

func TestLogger_URL_MaskSecurityCode(t *testing.T) {
	s := `securityCode=123`
	assert.Equal(t, "securityCode=***", logger.MaskCard(s))
}

func TestLogger_URL_MaskCVV(t *testing.T) {
	s := `cvv=123`
	assert.Equal(t, "cvv=***", logger.MaskCard(s))
}

func TestLogger_URL_MaskCard(t *testing.T) {
	s := `number=4111111111111111&cvv=123`
	assert.Equal(t, "cvv=***&number=************1111", logger.MaskCard(s))
}

func TestLogger_URL_MaskCard_EmptyString(t *testing.T) {
	s := ""
	assert.Equal(t, "", logger.MaskCard(s))
}

func TestLogger_URL_MaskCard_NoMatch(t *testing.T) {
	s := `hello=world&foo=bar`
	assert.Equal(t, "foo=bar&hello=world", logger.MaskCard(s))
}

func TestLogger_XML_MaskCard(t *testing.T) {
	s := `<Message><CardNumber>4111111111111111</CardNumber><Hello>World</Hello><SecurityCode>123</SecurityCode></Message>`
	assert.Equal(t, "<Message><CardNumber>************1111</CardNumber><Hello>World</Hello><SecurityCode>***</SecurityCode></Message>", logger.MaskCard(s))
}

func TestLogger_XML_MaskCardNumber(t *testing.T) {
	s := `<Message><Number>4111111111111111</Number></Message>`
	assert.Equal(t, "<Message><Number>************1111</Number></Message>", logger.MaskCard(s))
}

func TestLogger_XML_MaskCVV(t *testing.T) {
	s := `<Message><CVV>123</CVV></Message>`
	assert.Equal(t, "<Message><CVV>***</CVV></Message>", logger.MaskCard(s))
}
