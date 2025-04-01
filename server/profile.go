package server

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"shorty/src/common/metrics"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type profiler struct {
	mutex        *sync.Mutex
	file         *os.File
	stopChan     chan struct{}
	statusMetric metrics.Gauge
}

func (p *profiler) Start(ctx context.Context) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.file != nil {
		return "", fmt.Errorf("profiling already enabled")
	}

	var err error
	p.file, err = os.Create(fmt.Sprintf("./%d.pprof", time.Now().UnixMilli()))
	if err != nil {
		return "", fmt.Errorf("failed creating pprof file: %w", err)
	}
	fileName := p.file.Name()

	p.statusMetric.Set(1)
	p.stopChan = make(chan struct{})

	go func() {
		defer func() {
			pprof.StopCPUProfile()
			p.statusMetric.Set(0)
			p.file.Close()
			p.file = nil
		}()

		//TODO: add error handling
		pprof.StartCPUProfile(p.file)

		select {
		case <-ctx.Done():
		case <-p.stopChan:
		}
	}()

	return fileName, nil
}

func (p *profiler) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.file == nil {
		return fmt.Errorf("profiling not enabled")
	}

	close(p.stopChan)
	return nil
}

func (s *server) ProfileStart(c *gin.Context) {
	//TODO: change context
	fileName, err := s.profiler.Start(context.Background())
	if err != nil {
		c.Data(400, "application/json", []byte(fmt.Sprintf(`{"status": "error, "message": "%s"}`, err.Error())))
	} else {
		c.Data(200, "application/json", []byte(fmt.Sprintf(`{"status": "ok", "file": "%s"}`, fileName)))
	}
}

func (s *server) ProfileStop(c *gin.Context) {
	err := s.profiler.Stop()
	if err != nil {
		c.Data(400, "application/json", []byte(fmt.Sprintf(`{"status": "error, "message": "%s"}`, err.Error())))
	} else {
		c.Data(200, "application/json", []byte(`{"status": "ok"}`))
	}
}
