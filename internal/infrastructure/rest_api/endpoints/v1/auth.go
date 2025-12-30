package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e v1Engine) loginByCredentials(c *gin.Context) {
	var form struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBind(&form); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if result, err := e.authService.LoginByCredentials(c, form.Login, form.Password, e.cfg); err != nil {
		e.writeErr(c, http.StatusUnauthorized, err)

	} else {
		e.writeOk(c, result)
	}
}

func (e v1Engine) loginByToken(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if result, err := e.authService.LoginByToken(c, token, e.cfg); err != nil {
		e.writeErr(c, http.StatusUnauthorized, err)

	} else {
		e.writeOk(c, result)
	}
}
