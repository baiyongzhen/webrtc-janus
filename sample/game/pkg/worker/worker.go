package worker

import (
	"context"
	"log"

	"example.com/webrtc-game/pkg/config"
	"example.com/webrtc-game/pkg/monitoring"
)

type Worker struct {
	ctx context.Context
	cfg config.Config
	worker *Handler
	monitoringServer *monitoring.ServerMonitoring
}

func New(ctx context.Context, cfg config.Config) *Worker {
	return &Worker{
		ctx: ctx,
		cfg: cfg,

		monitoringServer: monitoring.NewServerMonitoring(cfg.MonitoringConfig),
	}
}

func (o *Worker) Run() error {
	go o.initializeWorker()
	return nil
}

func (o *Worker) Shutdown() {
	// 모니터링을 분리함.
}

// initializeWorker setup a worker
func (o *Worker) initializeWorker() {
	worker := NewHandler(o.cfg)
	defer func() {
		log.Println("Close worker")
		worker.Close()
	}()

	go worker.Run()

    //임시적으로 input keymap 테스트
	o.worker = worker
}

func (o *Worker) InputKeyboard() {
	o.worker.InputKeyboard()
}
