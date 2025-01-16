package controller

import (
	"com.banxiaoxiao.server/response"
	"com.banxiaoxiao.server/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TaskController struct {
}

var taskController *TaskController

func NewTaskController() *TaskController {
	if taskController == nil {
		taskController = &TaskController{}
	}
	return taskController
}

func (x *TaskController) ProcessTask(c echo.Context) error {
	taskService := service.GetTaskService()

	taskService.ProcessTask()
	return c.JSON(http.StatusOK, response.ReturnSuccessNoResultResponse())
}
