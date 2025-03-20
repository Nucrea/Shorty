package common

import (
	"math/rand"
	"strings"
	"time"
)

const ShortIdCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewShortId(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	for range size {
		rNum := r.Intn(len(ShortIdCharset))
		randChar := ShortIdCharset[rNum]
		sb.WriteByte(randChar)
	}

	return sb.String()
}

func NewDigitsString(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	for range size {
		randNum := byte(r.Intn(10))
		randByte := byte('0') + randNum
		sb.WriteByte(randByte)
	}

	return sb.String()
}

func MaskSecret(secret string) string {
	hLen := len(secret) / 2

	sb := strings.Builder{}
	sb.WriteString(secret[0:hLen])
	for range hLen {
		sb.WriteRune('*')
	}

	return sb.String()
}
