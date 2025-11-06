package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (v1 *v1Engine) loginByCredentials(c *gin.Context) {
	service := v1.getAuthService(c)

	var form struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBind(&form); err != nil {
		v1.writeErr(c, http.StatusBadRequest, err)

	} else if result, err := service.LoginByCredentials(form.Login, form.Password, v1.Config); err != nil {
		v1.writeErr(c, http.StatusUnauthorized, err)

	} else {
		v1.writeOk(c, result)
	}
}

func (v1 *v1Engine) loginByToken(c *gin.Context) {
	service := v1.getAuthService(c)

	if result, err := service.LoginByToken(c, v1.Config); err != nil {
		v1.writeErr(c, http.StatusUnauthorized, err)

	} else {
		v1.writeOk(c, result)
	}
}
