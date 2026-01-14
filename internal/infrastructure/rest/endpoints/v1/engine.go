package v1

import (
	"accounter/config"
	"accounter/internal/domain/event"
	"accounter/internal/domain/shared"
	"accounter/internal/domain/task"
	"accounter/internal/domain/user"
	"accounter/pkg/logger"
	"accounter/pkg/tools"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// API engine version 1
type v1Engine struct {
	cfg    config.Config
	logger logger.Logger

	taskService      *task.TaskService
	userService      *user.UserService
	eventService     *event.EventService
	authService      authService
	websocketService websocketService
}

// Websocket service
type websocketService interface {
	AcceptConnection(w http.ResponseWriter, r *http.Request, responseHeader http.Header) error
}

// Authorization service
type authService interface {
	LoginByToken(ctx context.Context, token string, cfg config.Config) (result user.CurrentUser, err error)
	LoginByCredentials(ctx context.Context, login, password string, cfg config.Config) (result user.CurrentUser, err error)
}

// Creates new v1Engine
func NewEngine(params EngineParams) *v1Engine {
	return &v1Engine{
		cfg:              params.Config,
		logger:           params.Logger,
		authService:      params.AuthService,
		taskService:      params.TaskService,
		userService:      params.UserService,
		eventService:     params.EventService,
		websocketService: params.WebsocketService,
	}
}

// Engine params
type EngineParams struct {
	Config           config.Config
	Logger           logger.Logger
	TaskService      *task.TaskService
	UserService      *user.UserService
	EventService     *event.EventService
	AuthService      authService
	WebsocketService websocketService
}

// getCurrentUser return CurrentUser from request context (Middlewares)
func (e *v1Engine) getCurrentUser(c *gin.Context) user.CurrentUser {
	user := c.MustGet("user").(user.CurrentUser)

	return user
}

// parseID parse ID from url query
func (e *v1Engine) parseID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

// RegisterRoutes register V1 routes
func (e *v1Engine) RegisterRoutes(s *gin.Engine) {
	v1 := s.Group("/api/v1")

	v1.GET("websocket", func(ctx *gin.Context) {
		e.websocketService.AcceptConnection(ctx.Writer, ctx.Request, nil)
	})

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

	v1.Group("/events").
		GET("list", e.getEventsList)
}

// Write success response
func (e *v1Engine) writeOk(c *gin.Context, data any) {
	resp := shared.Response[any]{
		Data:    data,
		Success: true,
		Status:  http.StatusOK,
	}

	c.JSON(http.StatusOK, resp)
}

// Write error response
func (e *v1Engine) writeErr(c *gin.Context, code int, err error) {
	resp := shared.Response[any]{
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
