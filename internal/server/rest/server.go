package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
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

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("failed to bind request JSON: %w", err))
		return
	}

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

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		h.Response(c, nil, http.StatusBadRequest, fmt.Errorf("error: %w", err))
		return
	}
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

	fmt.Println(filter)

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
