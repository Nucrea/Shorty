package cache

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func randStr(r *rand.Rand) string {
	sb := strings.Builder{}

	for range 16 {
		rNum := r.Intn(22)
		randChar := 'a' + byte(rNum)
		sb.WriteByte(randChar)
	}

	return sb.String()
}

func TestInmem(t *testing.T) {
	testValues := map[string]string{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	for range 256 {
		testValues[randStr(r)] = randStr(r)
	}

	inmem := NewInmem[string]()
	for k, v := range testValues {
		inmem.SetEx(k, v, time.Second)
	}

	for k, v := range testValues {
		if val, ok := inmem.Get(k); !ok || val != v {
			t.Fatal("wrong value")
		}
	}

	time.Sleep(time.Second)
	for k, v := range testValues {
		if val, ok := inmem.Get(k); ok || val == v {
			t.Fatal("key did not expire")
		}
	}
}
