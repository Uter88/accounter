package v1

import (
	"accounter/internal/domain/task"
	"accounter/pkg/tools"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getTasksList(c *gin.Context) {
	params := task.NewTaskParams()

	if err := c.ShouldBind(params); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if result, err := e.taskService.GetTaskList(c, params); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, result)
	}
}

func (e *v1Engine) saveTask(c *gin.Context) {

	var form task.Task

	if err := c.ShouldBind(&form); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err := e.taskService.SaveTask(c, &form); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, form)
	}
}

func (e *v1Engine) deleteTask(c *gin.Context) {
	if id, err := strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if err = e.taskService.DeleteTask(c, id); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeOk(c, "OK")
	}
}

func (e *v1Engine) exportTasks(c *gin.Context) {
	format := tools.FileFormat(c.Param("format"))
	params := task.NewTaskParams()

	if err := c.ShouldBind(params); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else if result, err := e.taskService.GetTaskList(c, params); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else if result, err := e.taskService.ExportTasks(result, format); err != nil {
		e.writeErr(c, http.StatusBadRequest, err)

	} else {
		e.writeBlob(c, format, result)
	}
}
