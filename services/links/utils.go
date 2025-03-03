package links

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateShortId(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	min := int('a')
	max := int('z')
	for range size {
		char := min + r.Intn(max-min)
		sb.WriteByte(byte(char))
	}

	return sb.String()
}
