package v1

import (
	"accounter/internal/domain/common"
	"accounter/internal/domain/task"
	"accounter/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getTasksList(c *gin.Context) {
	user := e.getCurrentUser(c)
	params := common.NewRequestParams(time.Now())

	if err := c.ShouldBind(&params); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if result, err := e.taskService.GetTaskList(user, params); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, result)
	}
}

func (e *v1Engine) saveTask(c *gin.Context) {
	user := e.getCurrentUser(c)
	var form task.Task

	if err := c.ShouldBind(&form); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err := e.taskService.SaveTask(user, &form); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, form)
	}
}

func (e *v1Engine) deleteTask(c *gin.Context) {
	user := e.getCurrentUser(c)

	if id, err := e.parseID(c); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err = e.taskService.DeleteTask(user, id); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, "OK")
	}
}

func (e *v1Engine) exportTasks(c *gin.Context) {
	user := e.getCurrentUser(c)
	format := utils.FileFormat(c.Param("format"))
	params := common.NewRequestParams(time.Now())

	if err := c.ShouldBind(&params); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if result, err := e.taskService.GetTaskList(user, params); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else if result, err := e.taskService.ExportTasks(result, format); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeBlob(c, format, result)
	}
}
