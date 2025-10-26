package v1

import (
	"accounter/adapters"
	"accounter/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getUsersList(c *gin.Context) {
	service := user.NewUserService(adapters.NewUserRepository(c, e.App.DbClient))

	if result, err := service.GetList(); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)
	} else {
		e.writeOk(c, result)
	}
}

func (e *v1Engine) saveUser(c *gin.Context) {
	service := user.NewUserService(adapters.NewUserRepository(c, e.App.DbClient))

	var form user.User

	if err := c.ShouldBind(&form); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err := service.SaveUser(&form); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, form)
	}
}

func (e *v1Engine) deleteUser(c *gin.Context) {

}
