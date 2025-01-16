package util

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type MyTask struct {
	name string
}

func (myTask *MyTask) Execute() error {
	time.Sleep(time.Second)
	fmt.Println(myTask.name)
	return nil
}

func (myTask *MyTask) Stop() error {
	return nil
}

func (myTask *MyTask) Equal(t interface{}) bool {
	return myTask.name == t.(*MyTask).name
}

func TestPool(t *testing.T) {
	pool := NewPool(2, -1)
	pool.Run()
	for i := 0; i < 100; i += 1 {
		pool.Add(&MyTask{name: "wyg = " + strconv.Itoa(i)})
		time.Sleep(time.Millisecond * 20)
	}
	pool.WaitAllDone()
	pool.Close()
}

func TestPool_Remove(t *testing.T) {
	pool := NewPool(20, 100)
	task := &MyTask{name: "wyg = 0"}
	pool.Add(task)
	pool.Remove(task)
	pool.WaitAllDone()
	pool.Close()
}
