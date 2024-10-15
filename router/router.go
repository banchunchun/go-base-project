package router

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/controller"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Init(e *echo.Echo) {
	conf := config.Cfg
	if conf.Extension.CorsEnabled {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"*"},
			AllowHeaders:     []string{"*"},
			AllowMethods: []string{
				http.MethodHead,
				http.MethodOptions,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			MaxAge: 86400,
		}))
	}

	errorHandler := controller.NewErrorController()
	e.HTTPErrorHandler = errorHandler.JSONError
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("100G"))

}
