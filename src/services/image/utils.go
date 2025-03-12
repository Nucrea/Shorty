package image

import (
	"math/rand"
	"strings"
	"time"
)

func NewShortId(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	randFuncs := []func(r *rand.Rand) byte{
		func(r *rand.Rand) byte {
			min, max := int('a'), int('z')
			return byte(min + r.Intn(max-min))
		},
		func(r *rand.Rand) byte {
			min, max := int('A'), int('Z')
			return byte(min + r.Intn(max-min))
		},
		func(r *rand.Rand) byte {
			min, max := int('0'), int('9')
			return byte(min + r.Intn(max-min))
		},
	}

	for range size {
		randFuncNum := r.Intn(len(randFuncs))
		randChar := randFuncs[randFuncNum](r)
		sb.WriteByte(randChar)
	}

	return sb.String()
}
