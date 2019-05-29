package services

import (
	"gt-monitor/common/zap"
	"gt-monitor/services/dao"
	"gt-monitor/services/delayed"
)

type ConvWorker struct {
	handler      *ConvHandler
	delayedTaskQ *delayed.DelayedTaskQueue
	dao          *dao.MongoDao
}

func NewConvWorker(convHandler *ConvHandler, delayedTaskQ *delayed.DelayedTaskQueue, dao *dao.MongoDao) *ConvWorker {
	return &ConvWorker{convHandler, delayedTaskQ, dao,}
}

func (d *ConvWorker) Run() {
	go func() {
		zap.Get().Info("start delay task queue worker")
		for task := range d.delayedTaskQ.EventChannel {
			d.handleTask(task)
		}
	}()
}

func (d *ConvWorker) handleTask(task delayed.Task) {
	if task.Cnt > 60 {
		zap.Get().Error("")
		_ = d.delayedTaskQ.Delete(task.Params.ImpId)
		return
	}
	params := task.Params
	clickEntity, err := d.dao.FindClickByImpId(params.ImpId, params.ClkTs)
	if err != nil {
		_ = d.delayedTaskQ.Put(delayed.Task{Cnt: task.Cnt + 1, Params: params})
		return
	}

	d.handler.DoHandle(params, clickEntity)
	_ = d.delayedTaskQ.Delete(params.ImpId)
}
