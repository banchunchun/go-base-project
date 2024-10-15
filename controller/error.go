package controller

import (
	"com.banxiaoxiao.server/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type APIError struct {
	Code    int
	Message string
}

type ErrorController struct {
}

func NewErrorController() *ErrorController {
	return &ErrorController{}
}

func (controller *ErrorController) JSONError(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}

	var apiError APIError
	apiError.Code = code
	apiError.Message = msg

	if !c.Response().Committed {
		if resErr := c.JSON(code, apiError); resErr != nil {
			logger.Log().Errorf(resErr.Error())
		}
	}
	logger.Log().Debugf(err.Error())
}
