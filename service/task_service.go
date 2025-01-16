package service

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	"com.banxiaoxiao.server/model"
	"com.banxiaoxiao.server/repo"
	"com.banxiaoxiao.server/util"
	"fmt"
	"sync"
	"time"
)

type TaskService struct {
}

var taskService *TaskService
var once sync.Once

func NewTaskService() *TaskService {
	once.Do(func() {
		if taskService == nil {
			taskService = &TaskService{}
		}
	})
	return taskService
}

func GetTaskService() *TaskService {
	return taskService
}

func (x *TaskService) ProcessTask() {
	datasource := repo.GetRepository()
	var taskList []*model.Task
	//db := datasource.Raw("SELECT * from task where task_status = 1 and features = 'tencentTag,hs-axr'")
	db := datasource.Raw("SELECT * from task where task_status = 1 and features = 'tencentTag,hs-axr' and created_time <=1729353600")
	db.Find(&taskList)
	fmt.Println(fmt.Sprintf("processTask length:%d", len(taskList)))
	for _, ts := range taskList {
		//var asr []*model.Asr
		//datasource.Raw("select * from huisheng_asr where task_id = ?", ts.TaskId).Find(&asr)
		//if len(asr) > 0 {
		//	continue
		//} else {
		url := fmt.Sprintf("%s/%s%d", config.Cfg.AiAddress, "api/task/stop?taskId=", ts.TaskId)
		logger.Log().Infof("stopUrl=%s", url)
		util.HttpPost(url, "")
		time.Sleep(1 * time.Second)
		//url = fmt.Sprintf("%s/%s%d", config.Cfg.AiAddress, "api/task/start?taskId=", ts.TaskId)
		//logger.Log().Infof("startUrl=%s", url)
		//util.HttpPost(url, "")
		//}
	}
}
