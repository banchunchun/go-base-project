package router

import (
	"com.banxiaoxiao.server/controller"
	"github.com/labstack/echo/v4"
)

func initTask(e *echo.Echo) {
	t := controller.NewTaskController()

	e.POST("/api/task/process", func(c echo.Context) error {
		return t.ProcessTask(c)
	})
}
