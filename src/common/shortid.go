package common

import (
	"math/rand"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewShortId(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	for range size {
		rNum := r.Intn(len(charset))
		randChar := charset[rNum]
		sb.WriteByte(randChar)
	}

	return sb.String()
}
