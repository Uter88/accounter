package v1

import (
	"accounter/domain/user"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getUsersList(c *gin.Context) {
	service := e.getUserService(c)

	if result, err := service.GetUsersList(); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, result)
	}
}

func (e *v1Engine) saveUser(c *gin.Context) {
	service := e.getUserService(c)

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
	service := e.getUserService(c)

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err = service.DeleteUser(id); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, "OK")
	}
}

func (e *v1Engine) isUserExists(c *gin.Context) {
	service := e.getUserService(c)

	if login, ok := c.GetQuery("login"); !ok {
		e.writeErr(c, http.StatusBadRequest, errors.New("required login param"))

	} else if exists, err := service.CheckUniqueLogin(login); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, exists)
	}
}
