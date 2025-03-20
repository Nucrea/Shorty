package common

import (
	"bytes"
	"encoding/base64"
	"regexp"

	"github.com/dchest/captcha"
)

var captchaRegexp = regexp.MustCompile(`^\d+$`)

func NewCaptchaImageBase64(id, value string) string {
	if !captchaRegexp.MatchString(value) {
		panic("captcha value must contain only digits")
	}

	digits := make([]byte, len(value))
	for i, c := range value {
		digits[i] = byte(c - '0')
	}

	buf := bytes.NewBuffer(nil)
	image := captcha.NewImage(id, digits, 200, 80)
	image.WriteTo(buf)

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
