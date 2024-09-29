package rest

import (
	"context"
	"net/http"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type TaskHandler struct {
	router  *gin.Engine
	service TaskServise
	log     *zerolog.Logger
}

type TaskServise interface {
	Create(ctx context.Context, title string, description string) (models.Task, error)
	Get(ctx context.Context, id uuid.UUID) (models.Task, error)
	Update(ctx context.Context, req models.Task) (models.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
}

func NewTaskHandler(svc TaskServise, log *zerolog.Logger) *TaskHandler {
	router := gin.Default()
	return &TaskHandler{
		router:  router,
		service: svc,
		log:     log,
	}
}

func (h *TaskHandler) registerRoutes() {
	tasks := h.router.Group("/tasks")
	{
		tasks.POST("/", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", h.DeleteTask)
		tasks.GET("/", h.ListTasks)
	}
}

func Run(server *TaskHandler, serverPort string) {
	server.registerRoutes()
	server.router.Run(":" + serverPort)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Create(c.Request.Context(), req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	task, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedTask, err := h.service.Update(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.service.List(c.Request.Context(), models.TaskFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
