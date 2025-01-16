package util

import (
	"com.banxiaoxiao.server/logger"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type TaskQueue[T any] struct {
	name         string
	callback     func(task T) error
	maxTasks     int32
	pendingLock  sync.RWMutex
	pendingTasks []T
	runningTasks []T
	runningCount int32
	c            int32
}

func NewTaskQueue[T any](name string, callback func(task T) error, maxTasks int32) *TaskQueue[T] {
	return &TaskQueue[T]{
		name:         name,
		callback:     callback,
		maxTasks:     maxTasks,
		pendingLock:  sync.RWMutex{},
		pendingTasks: make([]T, 0),
		runningTasks: make([]T, 0),
		runningCount: 0,
	}
}

func (tq *TaskQueue[T]) Add(task T) {
	tq.pendingLock.Lock()
	defer tq.pendingLock.Unlock()

	tq.pendingTasks = append(tq.pendingTasks, task)
}

func (tq *TaskQueue[T]) processQueue() {
	tq.pendingLock.Lock()
	defer tq.pendingLock.Unlock()
	tq.runningTasks = append(tq.runningTasks, tq.pendingTasks...)
	tq.pendingTasks = make([]T, 0)
}

func (tq *TaskQueue[T]) process() {
	tq.processQueue()
	l := int32(len(tq.runningTasks))
	runningCount := atomic.LoadInt32(&tq.runningCount)
	tq.c++
	if tq.c > 1000 {
		tq.c = 0
		if runningCount > 0 || l > 0 {
			logger.Log().Infof("JYW==> %s task queue, running: %d, pending: %d", tq.name, runningCount, l)
		}
	}
	if l == 0 {
		return
	}
	if l > tq.maxTasks*2 {
		i := l - tq.maxTasks*2
		logger.Log().Infof("JYW==> %s task queue is too long, drop %d tasks", tq.name, i)
		tq.runningTasks = tq.runningTasks[i:]
	}
	l = int32(len(tq.runningTasks))
	if runningCount < tq.maxTasks {
		n := tq.maxTasks - runningCount
		n = int32(math.Min(float64(n), float64(l)))
		for i := int32(0); i < n; i++ {
			a := tq.runningTasks[i]
			atomic.AddInt32(&tq.runningCount, 1)
			go func() {
				defer atomic.AddInt32(&tq.runningCount, -1)
				tq.callback(a)
			}()
		}
		tq.runningTasks = tq.runningTasks[n:]
	}
}

func (tq *TaskQueue[T]) Start() {
	go func() {
		for {
			time.Sleep(3 * time.Millisecond)
			tq.process()
		}
	}()
}
