package v1

import (
	"accounter/adapters/adapter_sql"
	"accounter/backend/core"
	"accounter/domain/task"
	"accounter/domain/user"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API engine version 1
type v1Engine struct {
	*core.App
}

// Creates new v1Engine
func NewEngine(app *core.App) *v1Engine {
	return &v1Engine{App: app}
}

// RegisterRoutes register V1 routes
func (e *v1Engine) RegisterRoutes(s *gin.Engine) {
	v1 := s.Group("/api/v1")

	v1.Group("/login").
		GET("", e.loginByToken).
		POST("", e.loginByCredentials)

	//v1.Use(e.userAuthentication())
	v1.Use(e.logger())

	v1.Group("/users").
		GET("/list", e.getUsersList).
		POST("/save", e.saveUser).
		DELETE("/delete/:id", e.deleteUser).
		GET("/is_exists", e.isUserExists)

	v1.Group("/tasks").
		GET("/list", e.getTasksList).
		POST("/save", e.saveTask).
		DELETE("/delete/:id", e.deleteTask)
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

	c.JSON(code, resp)
}

// Creates new User service
func (e *v1Engine) getUserService(ctx context.Context) user.UserService {
	return user.NewUserService(adapter_sql.NewUserRepository(ctx, e.DbClient))
}

// Creates new Task service
func (e *v1Engine) getTaskService(ctx context.Context) task.TaskService {
	return task.NewTaskService(adapter_sql.NewTaskRepository(ctx, e.DbClient))
}

// Creates new auth service
func (e *v1Engine) getAuthService(ctx context.Context) core.AuthService {
	return core.NewAuthService(adapter_sql.NewUserRepository(ctx, e.DbClient))
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
