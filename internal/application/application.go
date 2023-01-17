package application

import (
	"context"
	"encoding/json"
	"sync"

	"git.carriot.ir/warning-detector/internal/models"
	p "git.carriot.ir/warning-detector/internal/ports"
	"go.uber.org/zap"
)

type app struct {
	queue     p.LogBroker
	logRepo   p.LogRepository
	detector  p.WarningDetector
	warnQueue p.WarningQueue
	warnRepo  p.WarningRepository
}

func New(q p.LogBroker, lr p.LogRepository, wd p.WarningDetector, wq p.WarningQueue, wr p.WarningRepository) (*app, error) {
	res := app{
		queue:     q,
		logRepo:   lr,
		detector:  wd,
		warnQueue: wq,
		warnRepo:  wr,
	}
	return &res, nil
}

func (a *app) Do(ctx context.Context) {
	a.queue.SetHandler(a.logConsumer)
	a.queue.Subscribe(ctx)
	a.warnQueue.SetHandler(a.warnQueueConsumer)
	a.warnQueue.Subscribe(ctx)
}

func (a *app) logConsumer(ctx context.Context, msg []byte) error {
	log := models.TempLog{}
	err := json.Unmarshal(msg, &log)
	if err != nil {
		zap.L().Warn(err.Error())
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logRepo.Store(ctx, log)
	}()
	data, err := json.Marshal(&log)
	if err != nil {
		return err
	}
	a.warnQueue.PublishAllLogs(ctx, string(data))
	wg.Wait()
	return nil
}

func (a *app) warnQueueConsumer(ctx context.Context, msg []byte) error {
	log := models.TempLog{}
	err := json.Unmarshal(msg, &log)
	if err != nil {
		zap.L().Error(err.Error())
		// TODO retry
	}
	return a.checkWarning(ctx, log)
}

func (a *app) checkWarning(ctx context.Context, log models.TempLog) error {
	warnLog, err := a.detector.Detect(log)
	if err != nil {
		return err
	}
	if warnLog == nil {
		return nil
	}
	return a.handleWarnLog(ctx, *warnLog)
}

func (a *app) handleWarnLog(ctx context.Context, log models.HasWarnings) error {
	var e error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e = a.warnRepo.Store(ctx, log)
	}()
	data, err := json.Marshal(&log)
	if err != nil {
		return err
	}
	err = a.warnQueue.PublishWarningLogs(ctx, string(data))
	if e != nil {
		return e
	}
	if err != nil {
		return err
	}
	return nil
}
