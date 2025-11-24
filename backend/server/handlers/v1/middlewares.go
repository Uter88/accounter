package v1

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Try to authorize CurrentUser and put it to context
func (v1 *v1Engine) userAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		service := v1.getAuthService(c)

		if user, err := service.LoginByToken(c, v1.Config); err != nil {
			v1.writeErr(c, http.StatusUnauthorized, err)

		} else {
			c.Set("user", user)
		}

		c.Next()
	}
}

// Log request
func (v1 *v1Engine) logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Set("requestStartTime", t)
		c.Next()
		latency := time.Since(t)
		status := c.Writer.Status()
		err := c.Err()

		if v1.Config.DebugMode {
			msg := fmt.Sprintf(
				"Path: %s, method: %s, latency: %s, status: %d, address: %s",
				c.Request.URL.String(), c.Request.Method, latency, status, c.Request.RemoteAddr,
			)

			if err != nil {
				msg += fmt.Sprintf(", error: %s", err)
			} else if l := len(c.Errors); l > 0 {
				errs := make([]string, l)

				for i, e := range c.Errors {
					errs[i] = e.Error()
				}

				msg += fmt.Sprintf(", errors: %s", strings.Join(errs, ","))
			}

			v1.Logger.Info(msg)
		}
	}
}
