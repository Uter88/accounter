package v1

import (
	"accounter/internal/domain/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getUsersList(c *gin.Context) {
	if result, err := e.userService.GetUsersList(c); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, result)
	}
}

func (e *v1Engine) saveUser(c *gin.Context) {
	var form user.User

	if err := c.ShouldBind(&form); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err := e.userService.SaveUser(e.getCurrentUser(c), &form); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, form)
	}
}

func (e *v1Engine) deleteUser(c *gin.Context) {
	if id, err := e.parseID(c); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if user, err := e.userService.GetUser(c, id); err != nil {
		e.writeErr(c, http.StatusNotFound, err)

	} else if err = e.userService.DeleteUser(e.getCurrentUser(c), &user); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, "OK")
	}
}

func (e *v1Engine) isUserExists(c *gin.Context) {
	if login, ok := c.GetQuery("login"); !ok {
		e.writeErr(c, http.StatusBadRequest, errors.New("required login param"))

	} else if exists, err := e.userService.CheckUniqueLogin(c, login); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, exists)
	}
}
