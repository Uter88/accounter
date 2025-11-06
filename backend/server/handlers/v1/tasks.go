package v1

import (
	"accounter/domain/task"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getTasksList(c *gin.Context) {
	service := e.getTaskService(c)

	if result, err := service.GetTaskList(); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, result)
	}
}

func (e *v1Engine) saveTask(c *gin.Context) {
	service := e.getTaskService(c)

	var form task.Task

	if err := c.ShouldBind(&form); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err := service.SaveTask(&form); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, form)
	}
}

func (e *v1Engine) deleteTask(c *gin.Context) {
	service := e.getTaskService(c)

	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err = service.DeleteTask(id); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, "OK")
	}
}
