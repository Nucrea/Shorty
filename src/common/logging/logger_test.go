package logging

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func TestLogger(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "log_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	logger, err := NewLogger(WithFile(tmpFile.Name()))
	if err != nil {
		t.Fatal(err)
	}

	linesToWrite := 10000
	counter := atomic.Int32{}
	threads := runtime.NumCPU()
	wg := sync.WaitGroup{}

	for range threads {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if counter.Add(1) > int32(linesToWrite) {
					break
				}
				logger.Info().Msg("test")
			}
		}()
	}
	wg.Wait()

	time.Sleep(time.Second)

	info, err := tmpFile.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if info.Size() <= 0 {
		t.Fatal("empty tmpFile")
	}

	lines, err := lineCounter(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if lines != linesToWrite {
		t.Fatal("lines count not equal")
	}
}

func BenchmarkLogger(b *testing.B) {
	b.StopTimer()

	tmpFile, err := os.CreateTemp("", "log_bench_*")
	if err != nil {
		b.Fatal(err)
	}
	defer tmpFile.Close()

	logger, err := NewLogger(WithFile(tmpFile.Name()))
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			logger.Log().Msg("test")
		}
	})
}
