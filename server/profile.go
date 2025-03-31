package server

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	once sync.Once
	prof *profiler
)

func init() {
	once.Do(func() {
		prof = &profiler{
			mutex: &sync.Mutex{},
		}
	})
}

type profiler struct {
	mutex    *sync.Mutex
	file     *os.File
	stopChan chan struct{}
}

func (p *profiler) Start(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.file != nil {
		return fmt.Errorf("profiling already enabled")
	}

	var err error
	p.file, err = os.Create(fmt.Sprintf("./%d.pprof", time.Now().UnixMilli()))
	if err != nil {
		return fmt.Errorf("failed creating pprof file: %w", err)
	}

	p.stopChan = make(chan struct{})
	go func() {
		pprof.StartCPUProfile(p.file)
		defer func() {
			pprof.StopCPUProfile()
			p.file.Close()
		}()

		select {
		case <-ctx.Done():
		case <-p.stopChan:
		}

	}()

	return nil
}

func (p *profiler) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.file == nil {
		return fmt.Errorf("profiling not enabled")
	}

	close(p.stopChan)
	p.file = nil
	return nil
}

func (s *server) ProfileStart(c *gin.Context) {
	err := prof.Start(context.Background())
	if err != nil {
		c.Data(400, "application/json", []byte(fmt.Sprintf(`{"status": "error, "message": "%s"}`, err.Error())))
	} else {
		c.Data(200, "application/json", []byte(`{"status": "ok"}`))
	}
}

func (s *server) ProfileStop(c *gin.Context) {
	err := prof.Stop()
	if err != nil {
		c.Data(400, "application/json", []byte(fmt.Sprintf(`{"status": "error, "message": "%s"}`, err.Error())))
	} else {
		c.Data(200, "application/json", []byte(`{"status": "ok"}`))
	}
}
