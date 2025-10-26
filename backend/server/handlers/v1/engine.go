package v1

import (
	"accounter/backend/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type v1Engine struct {
	*core.App
}

func NewEngine(app *core.App) *v1Engine {
	return &v1Engine{App: app}
}

func (e *v1Engine) RegisterRoutes(s *gin.Engine) {
	v1 := s.Group("/api/v1")

	//v1.Use(e.userAuthentication())
	v1.Use(e.logger())

	v1.Group("/users").
		GET("/list", e.getUsersList).
		POST("/save", e.saveUser).
		DELETE("/delete/:id", e.deleteUser)

	v1.Group("/tasks").
		GET("/list", e.getTasksList).
		POST("/save", e.saveTask).
		DELETE("/delete/:id", e.deleteTask)
}

func (e *v1Engine) writeOk(c *gin.Context, data any) {
	resp := Response[any]{
		Data:    data,
		Success: true,
		Status:  http.StatusOK,
	}

	c.JSON(http.StatusOK, resp)
}

func (e *v1Engine) writeErr(c *gin.Context, code int, err error) {
	resp := Response[any]{
		Status: code,
		Error:  err.Error(),
	}

	c.JSON(code, resp)
}

type Response[T any] struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Data    T      `json:"data"`
	Error   string `json:"error"`

	TotalRows int `json:"total_rows"`
}
