package middleware

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/logger"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/valyala/fasttemplate"
)

func InitLoggerMiddleware(e *echo.Echo) {
	e.Use(RequestLoggerMiddleware())
	e.Use(ActionLoggerMiddleware())
	if len(config.Cfg.QuerySecret.Value) > 0 {
		e.Use(QuerySecretKeyMiddleware())
	}
}

func QuerySecretKeyMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + config.Cfg.QuerySecret.Key,
		AuthScheme: "",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == config.Cfg.QuerySecret.Value, nil
		},
	})
}

func RequestLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			if err := next(c); err != nil {
				c.Error(err)
			}

			template := fasttemplate.New(config.Cfg.Log.RequestLogFormat, "${", "}")
			logStr := template.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "remote_ip":
					return w.Write([]byte(c.RealIP()))
				case "uri":
					return w.Write([]byte(req.RequestURI))
				case "method":
					return w.Write([]byte(req.Method))
				case "status":
					return w.Write([]byte(strconv.Itoa(res.Status)))
				case "content_length":
					return w.Write([]byte(strconv.FormatInt(req.ContentLength, 10)))
				case "response_length":
					return w.Write([]byte(strconv.FormatInt(res.Size, 10)))
				default:
					return w.Write([]byte(""))
				}
			})
			logger.Log().Infof(logStr)
			return nil
		}
	}
}

func ActionLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			st := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			et := time.Now()
			delta := et.UnixMilli() - st.UnixMilli()
			if delta > 20 {
				logger.Log().Infof("%s: %d ms", c.Path(), delta)
			}
			return nil
		}
	}
}
