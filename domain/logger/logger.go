package logger

import (
	"bytes"
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const key = "logger"

func Middleware(log *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqID := c.Request().Header.Get("X-Request-ID")
			l := log.With(zap.String("req-id", reqID))
			c.Set(key, l)

			// Read the Body content
			var bodyBytes []byte
			if c.Request().Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
			}

			// Restore the io.ReadCloser to its original state
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			// fmt.Printf("request body: %v\n", c.Request().Body)

			return next(c)
		}
	}
}

func Unwrap(c echo.Context) *zap.Logger {
	val := c.Get(key)
	if log, ok := val.(*zap.Logger); ok {
		return log
	}

	return zap.NewExample()
}
