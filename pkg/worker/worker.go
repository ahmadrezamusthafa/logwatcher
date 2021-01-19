package worker

import (
	"github.com/ahmadrezamusthafa/logwatcher/common/logger"
	"github.com/ahmadrezamusthafa/logwatcher/config"
	"github.com/gammazero/workerpool"
)

type EngineWorker struct {
	Config     config.Config `inject:"config"`
	WorkerPool *workerpool.WorkerPool
}

func (mw *EngineWorker) StartUp() {
	logger.Info("Init worker... ")
	mw.WorkerPool = workerpool.New(mw.Config.MaxWorkerPool)
}

func (mw *EngineWorker) Shutdown() {
	logger.Info("Stopping worker... ")
	mw.WorkerPool.StopWait()
}
