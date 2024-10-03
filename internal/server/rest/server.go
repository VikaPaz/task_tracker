package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/VikaPaz/task_tracker/docs"
	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type TaskHandler struct {
	router   *gin.Engine
	service  TaskServise
	validate *validator.Validate
	log      *zerolog.Logger
}

type TaskServise interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	Get(ctx context.Context, id uuid.UUID) (models.Task, error)
	Update(ctx context.Context, req models.Task) (models.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
}

func NewTaskHandler(svc TaskServise, log *zerolog.Logger) *TaskHandler {
	router := gin.Default()
	validate := validator.New()
	return &TaskHandler{
		router:   router,
		service:  svc,
		validate: validate,
		log:      log,
	}
}

func (h *TaskHandler) registerRoutes() {
	tasks := h.router.Group("/task")
	{
		tasks.POST("/", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.PUT("/", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
		tasks.GET("/", h.ListTasks)
	}
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func Run(server *TaskHandler, serverPort string) {
	server.registerRoutes()
	server.router.Run(":" + serverPort)
}

func (h *TaskHandler) Response(
	c *gin.Context,
	responseBody interface{},
	status int,
	err error,
) {
	if err != nil {
		responseBody = gin.H{"error": err.Error()}
	}
	h.logRequest(c, status, err)
	c.JSON(status, responseBody)
}

func (h *TaskHandler) logRequest(c *gin.Context, status int, err error) {
	logger := h.log.Info()

	if err != nil {
		logger = h.log.Error().Str("error", err.Error())
	}

	logger.
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.String()).
		Str("client_ip", c.ClientIP()).
		Int("status", status)
}

// TODO: add auth
func genOwner() string {
	return uuid.New().String()
}

type CreateRequest struct {
	ID          string `swaggerignore:"true" validate:"omitempty,uuid4"`
	Title       string
	Description string
	Created     time.Time         `swaggerignore:"true"`
	Updated     time.Time         `swaggerignore:"true"`
	Status      models.TaskStatus `validate:"omitempty,oneof=in_progress done"`
	OwnerID     string            `swaggerignore:"true" validate:"uuid4"`
}

// @Summary Creating a new task
// @Description Handles request to create a new task and returns the task information in JSON.
// @Tags task
// @Accept json
// @Produce json
// @Param request body CreateRequest true "New task"
// @Success 201 {object} models.Task "Created task"
// @Failure 400
// @Failure 500
// @Router /task/ [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("failed to bind request JSON: %w", err))
		return
	}

	task := models.Task(req)

	task.OwnerID = genOwner()

	if task.Status == "" {
		task.Status = models.InProgress
	}

	var err error
	err = h.validate.Struct(task)
	if err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("failed to bind request JSON: %w", err))
		return
	}
	h.log.Debug().Msg("validated new task")

	task, err = h.service.Create(c.Request.Context(), task)
	if err != nil {
		h.Response(c, nil, http.StatusInternalServerError, fmt.Errorf("failed to create task: %w", err))
		return
	}
	h.Response(c, gin.H{"task": task}, http.StatusCreated, nil)
}

// @Summary Receiving a task
// @Description Handles request to get a task and returns the task information in JSON.
// @Tags task
// @Produce json
// @Param id path string false "Task ID"
// @Success 200 {object} models.Task "task"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /task/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("invalid UUID: %w", err))
		return
	}

	task, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrTaskNotFound {
			status = http.StatusNotFound
		}
		h.Response(c, nil, status, fmt.Errorf("failed to receive task: %w", err))
		return
	}

	h.Response(c, gin.H{"task": task}, http.StatusOK, nil)
}

type UpdateRequest struct {
	ID          string `validate:"omitempty,uuid4"`
	Title       string
	Description string
	Created     time.Time         `swaggerignore:"true"`
	Updated     time.Time         `swaggerignore:"true"`
	Status      models.TaskStatus `validate:"omitempty,oneof=in_progress done"`
	OwnerID     string            `validate:"uuid4"`
}

// @Summary Updating a task
// @Description Handles request to update a task and returns the task information in JSON.
// @Tags task
// @Accept json
// @Produce json
// @Param request body UpdateRequest true "fields"
// @Success 200 {object} models.Task "Updated task"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /task/ [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req UpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("error: %w", err))
		return
	}

	task := models.Task(req)

	task.OwnerID = genOwner()
	_, err := uuid.Parse(task.ID)
	if err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("invalid UUID: %w", err))
		return
	}

	err = h.validate.Struct(task)
	if err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("failed to bind request JSON: %w", err))
		return
	}
	h.log.Debug().Msg("validated update task")

	updatedTask, err := h.service.Update(c.Request.Context(), task)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrTaskNotFound {
			status = http.StatusNotFound
		}
		h.Response(c, nil, status, fmt.Errorf("failed to update task: %w", err))
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// @Summary Deleting a task
// @Description Handles request to delete a task.
// @Tags task
// @Param id path string false "Task ID"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /task/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("invalid UUID: %w", err))
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrTaskNotFound {
			status = http.StatusNotFound
		}
		h.Response(c, nil, status, fmt.Errorf("failed to delete task: %w", err))
		return
	}

	h.Response(c, nil, http.StatusNoContent, nil)
}

// @Summary Listing a task
// @Description Handles request to get tasks and returns the list of tasks information in JSON.
// @Tags task
// @Produce json
// @Param id query string false "Task ID"
// @Param title query string false "Title"
// @Param description query string false "Description"
// @Param status query string false "Status"
// @Param owner_id query string false "Owner ID"
// @Success 200 {object} models.Task "task"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /task/ [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	var filter models.TaskFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("error: %w", err))
		return
	}

	if err := h.validate.Struct(filter); err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("failed to bind request JSON: %w", err))
		return
	}
	h.log.Debug().Msg("validated a filter")

	tasks, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrTaskNotFound {
			status = http.StatusNotFound
		}
		h.Response(c, nil, status, fmt.Errorf("failed to list tasks: %w", err))
		return
	}

	h.Response(c, gin.H{"tasks": tasks}, http.StatusOK, nil)
}
