package v1

import (
	"accounter/adapters/auth"
	"accounter/config"
	"accounter/internal/domain/task"
	"accounter/internal/domain/user"
	"accounter/pkg/tools"
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API engine version 1
type v1Engine struct {
	cfg    config.Config
	logger config.Logger

	TaskService task.TaskService
	UserService user.UserService
	AuthService AuthService
}

// AuthService authorization
type AuthService interface {
	LoginByToken(ctx context.Context, token string, cfg config.Config) (result auth.CurrentUser, err error)
	LoginByCredentials(ctx context.Context, login, password string, cfg config.Config) (result auth.CurrentUser, err error)
}

// Creates new v1Engine
func NewEngine(cfg config.Config, logger config.Logger) v1Engine {
	return v1Engine{cfg: cfg, logger: logger}
}

// RegisterRoutes register V1 routes
func (e v1Engine) RegisterRoutes(s *gin.Engine) {
	v1 := s.Group("/api/v1")

	v1.Group("/login").
		GET("", e.loginByToken).
		POST("", e.loginByCredentials)

	v1.Use(e.userAuthentication(), e.logging())

	v1.Group("/users").
		GET("/list", e.getUsersList).
		POST("/save", e.saveUser).
		DELETE("/delete/:id", e.deleteUser).
		GET("/is_exists", e.isUserExists)

	v1.Group("/tasks").
		GET("/list", e.getTasksList).
		POST("/save", e.saveTask).
		DELETE("/delete/:id", e.deleteTask).
		GET("/export/:format", e.exportTasks)
}

// Write success response
func (e *v1Engine) writeOk(c *gin.Context, data any) {
	resp := Response[any]{
		Data:    data,
		Success: true,
		Status:  http.StatusOK,
	}

	c.JSON(http.StatusOK, resp)
}

// Write error response
func (e *v1Engine) writeErr(c *gin.Context, code int, err error) {
	resp := Response[any]{
		Status: code,
		Error:  err.Error(),
	}

	c.AbortWithStatusJSON(code, resp)
}

// Write blob response
func (e *v1Engine) writeBlob(c *gin.Context, format tools.FileFormat, content *bytes.Buffer) {
	switch format {
	case tools.FileFormatXLSX:
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	case tools.FileFormatCSV:
		c.Header("Content-Type", "text/csv")

	case tools.FileFormatJSON:
		c.Header("Content-Type", "application/json")

	case tools.FileFormatHTML:
		c.Header("Content-Type", "text/html")

	case tools.FileFormatPDF:
		c.Header("Content-Type", "application/pdf")

	case tools.FileFormatDocX:
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="tasks.%s"`, format))
	c.Header("Content-Length", fmt.Sprintf("%d", content.Len()))

	content.WriteTo(c.Writer)
}

// Response struct
type Response[T any] struct {
	// Response success status
	Success bool `json:"success"`

	// Response HTTP status code
	Status int `json:"status"`

	// Response payload
	Data T `json:"data"`

	// Response error
	Error string `json:"error"`

	// Total rows counter (for pagination)
	TotalRows int `json:"total_rows"`
}
