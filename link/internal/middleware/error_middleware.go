package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"link/internal/utils"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			_, file, line, _ := runtime.Caller(1)

			switch e := err.(type) {
			case *utils.AppError:
				return c.JSON(e.StatusCode, map[string]interface{}{
					"error":       e.Message,
					"status_code": e.StatusCode,
					"file":        file,
					"line":        line,
				})
			case *echo.HTTPError:
				return c.JSON(e.Code, map[string]interface{}{
					"error":       e.Message,
					"status_code": e.Code,
					"file":        file,
					"line":        line,
				})
			default:
				// Log the full error for server-side debugging
				fmt.Printf("Unexpected error: %v\n", err)
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error":       "Internal server error",
					"status_code": http.StatusInternalServerError,
					"file":        file,
					"line":        line,
					"details":     fmt.Sprintf("%v", err),
				})
			}
		}
		return nil
	}
}
