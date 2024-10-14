package service

import (
	"context"
	"testing"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var logger = log.Logger
var repo = TestRepo{}

type TestRepo struct{}

func (s *TestRepo) Create(ctx context.Context, task models.Task) (models.Task, error) {
	return models.Task{ID: uuid.New().String(), Title: task.Title, Description: task.Description}, nil
}

func (s *TestRepo) Get(ctx context.Context, id uuid.UUID) (models.Task, error) {
	return models.Task{ID: id.String(), Title: "Test Task", Description: "Test Description"}, nil
}

func (s *TestRepo) Update(ctx context.Context, req models.Task) (models.Task, error) {
	return models.Task{ID: req.ID, Title: req.Title, Description: req.Description}, nil
}

func (s *TestRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *TestRepo) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	tasks := []models.Task{
		{ID: uuid.New().String(), Title: "Test Task 1", Description: "Test Description 1"},
		{ID: uuid.New().String(), Title: "Test Task 2", Description: "Test Description 2"},
	}
	return tasks, nil
}

func TestNewServise(t *testing.T) {
	svc := NewTaskService(&repo, &logger)

	if svc == nil {
		t.Fatalf("expected service to be initialized, got nil")
	}
}

func TestCreate(t *testing.T) {
	svc := NewTaskService(&repo, &logger)

	task := models.Task{
		Title:       "New Task",
		Description: "New Task Description",
	}
	createdTask, err := svc.Create(context.Background(), task)

	assert.NoError(t, err)
	assert.Equal(t, "New Task", createdTask.Title)
	assert.Equal(t, "New Task Description", createdTask.Description)
	assert.NotEqual(t, uuid.Nil, createdTask.ID)
}

func TestGet(t *testing.T) {
	svc := NewTaskService(&repo, &logger)

	taskID := uuid.New()
	task, err := svc.Get(context.Background(), taskID)

	assert.NoError(t, err)
	assert.Equal(t, taskID.String(), task.ID)
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "Test Description", task.Description)
}

func TestUpdate(t *testing.T) {
	svc := NewTaskService(&repo, &logger)

	task := models.Task{
		ID:          uuid.New().String(),
		Title:       "Updated Task",
		Description: "Updated Description",
	}
	updatedTask, err := svc.Update(context.Background(), task)

	assert.NoError(t, err)
	assert.Equal(t, task.ID, updatedTask.ID)
	assert.Equal(t, "Updated Task", updatedTask.Title)
	assert.Equal(t, "Updated Description", updatedTask.Description)
}

func TestDelete(t *testing.T) {
	svc := NewTaskService(&repo, &logger)

	err := svc.Delete(context.Background(), uuid.New())

	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	svc := NewTaskService(&repo, &logger)

	tasks, err := svc.List(context.Background(), models.TaskFilter{})

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "Test Task 1", tasks[0].Title)
	assert.Equal(t, "Test Task 2", tasks[1].Title)
}
