package util

import (
	"com.banxiaoxiao.server/logger"
	"container/list"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

func defaultPanicHandler(p interface{}) {
	logger.Log().Errorf("Worker exits from a panic: %v\nStack trace: %s\n", p, string(debug.Stack()))
}

type Job interface {
	Execute() error
	Stop() error
	Equal(t interface{}) bool
}

type Pool struct {
	lock             sync.RWMutex
	pendingQueue     list.List
	runningQueue     list.List
	maxWorkerNum     int32
	maxWaitingNum    int32
	waitingTaskCount int32
	runningTaskCount int32
	stopSignal       chan bool
	taskPoolWait     *sync.WaitGroup
	panicHandler     func(interface{})
}

func NewPool(maxWorkerNum int32, maxWaitingNum int32) *Pool {
	pool := &Pool{
		lock:          sync.RWMutex{},
		maxWorkerNum:  maxWorkerNum,
		maxWaitingNum: maxWaitingNum,
		stopSignal:    make(chan bool, 1),
		taskPoolWait:  &sync.WaitGroup{},
		panicHandler:  defaultPanicHandler,
	}
	return pool
}

func (p *Pool) SetMaxWorkerNumber(maxWorkerNum int32) {
	atomic.StoreInt32(&p.maxWorkerNum, maxWorkerNum)
}

func (p *Pool) WaitingTasks() int32 {
	return atomic.LoadInt32(&p.waitingTaskCount)
}

func (p *Pool) RunningTasks() int32 {
	return atomic.LoadInt32(&p.runningTaskCount)
}

func (p *Pool) Add(t Job) {
	defer p.lock.Unlock()
	p.lock.Lock()
	if p.maxWaitingNum > 0 && p.waitingTaskCount > p.maxWaitingNum {
		t := p.pendingQueue.Front()
		p.pendingQueue.Remove(t)
		atomic.AddInt32(&p.waitingTaskCount, -1)
	}
	p.pendingQueue.PushBack(t)
	atomic.AddInt32(&p.waitingTaskCount, 1)
}

func (p *Pool) find(t Job, l list.List) *list.Element {
	for e := l.Front(); e != nil; e = e.Next() {
		if t.Equal(e.Value.(Job)) {
			return e
		}
	}
	return nil
}

func (p *Pool) Exist(t Job) bool {
	defer p.lock.Unlock()
	p.lock.Lock()
	found := false
	if e := p.find(t, p.runningQueue); e != nil {
		found = true
	}
	if !found {
		if e := p.find(t, p.pendingQueue); e != nil {
			found = true
		}
	}
	return found
}

func (p *Pool) ExistPending(t Job) bool {
	defer p.lock.Unlock()
	p.lock.Lock()
	found := false
	if e := p.find(t, p.pendingQueue); e != nil {
		found = true
	}
	return found
}

func (p *Pool) Remove(t Job) bool {
	defer p.lock.Unlock()
	p.lock.Lock()
	found := false
	if e := p.find(t, p.runningQueue); e != nil {
		go func() {
			e.Value.(Job).Stop()
		}()
		p.runningQueue.Remove(e)
		atomic.AddInt32(&p.runningTaskCount, -1)
		found = true
	}
	if !found {
		if e := p.find(t, p.pendingQueue); e != nil {
			go func() {
				e.Value.(Job).Stop()
			}()
			p.pendingQueue.Remove(e)
			atomic.AddInt32(&p.waitingTaskCount, -1)
			found = true
		}
	}
	logger.Log().Infof("task waiting task = %d, running task = %d", atomic.LoadInt32(&p.waitingTaskCount), atomic.LoadInt32(&p.runningTaskCount))
	return found
}

func (p *Pool) Run() {
	p.taskPoolWait.Add(1)
	go func() {
		p.loop()
	}()
}

func (p *Pool) Close() {
	p.stopSignal <- true
	p.taskPoolWait.Wait()
}

func (p *Pool) WaitAllDone() {
	for {
		if p.WaitingTasks() == 0 {
			break
		}
		select {
		case <-time.After(time.Second):
			logger.Log().Debug("pool loop timeout")
		}
	}
}

func (p *Pool) loop() {
	logger.Log().Info("enter pool loop")
	exitFlag := false
	for !exitFlag {
		if p.WaitingTasks() > 0 {
			if p.RunningTasks() < atomic.LoadInt32(&p.maxWorkerNum) {
				p.startTask()
				continue
			}
		}
		select {
		case <-time.After(time.Second):
			//logger.Log().Debug("pool loop timeout")
		case <-p.stopSignal:
			logger.Log().Info("pool loop stop signal")
			exitFlag = true
		}
	}
	p.taskPoolWait.Done()
	logger.Log().Info("leave pool loop")
}

func (p *Pool) startTask() {
	defer p.lock.Unlock()
	p.lock.Lock()
	t := p.pendingQueue.Front()
	if t != nil {
		p.runningQueue.PushBack(t.Value)
		p.pendingQueue.Remove(t)
		atomic.AddInt32(&p.waitingTaskCount, -1)
		atomic.AddInt32(&p.runningTaskCount, 1)
		logger.Log().Infof("enter start task waiting task = %d, running task = %d", atomic.LoadInt32(&p.waitingTaskCount), atomic.LoadInt32(&p.runningTaskCount))
		go func() {
			p.runTask(t.Value.(Job))
		}()
	}
}

func (p *Pool) runTask(t Job) {
	defer func() {
		if _p := recover(); _p != nil {
			p.panicHandler(_p)
		}
		defer p.lock.Unlock()
		p.lock.Lock()
		if e := p.find(t, p.runningQueue); e != nil {
			p.runningQueue.Remove(e)
			atomic.AddInt32(&p.runningTaskCount, -1)
		}
		logger.Log().Infof("leave running task waiting = %d, running task = %d", atomic.LoadInt32(&p.waitingTaskCount), atomic.LoadInt32(&p.runningTaskCount))
	}()
	t.Execute()
}
