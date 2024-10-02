package service

import (
	"context"
	"time"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Repo interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	Get(ctx context.Context, id uuid.UUID) (models.Task, error)
	Update(ctx context.Context, req models.Task) (models.Task, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error)
}

type TaskService struct {
	repo Repo
	log  *zerolog.Logger
}

func NewTaskService(repo Repo, log *zerolog.Logger) *TaskService {
	return &TaskService{
		repo: repo,
		log:  log,
	}
}

func (s *TaskService) Create(ctx context.Context, task models.Task) (models.Task, error) {
	s.log.Debug().Msgf("Creating task: %v", task)

	task.Created, task.Updated = time.Now().UTC(), time.Now().UTC()

	task, err := s.repo.Create(ctx, task)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create task")
		return models.Task{}, err
	}
	s.log.Debug().Msg("created new task")

	return task, nil
}

func (s *TaskService) Get(ctx context.Context, id uuid.UUID) (models.Task, error) {
	s.log.Debug().Msgf("Fetching task with ID: %s", id.String())

	task, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Msgf("Error fetching task with ID: %s", id.String())
		return models.Task{}, err
	}
	s.log.Debug().Msg("received new task")

	return task, nil
}

func (s *TaskService) Update(ctx context.Context, req models.Task) (models.Task, error) {
	s.log.Info().Msgf("Updating task with ID: %s", req.ID)

	req.Updated = time.Now()

	task, err := s.repo.Update(ctx, req)
	if err != nil {
		s.log.Error().Err(err).Msgf("Error updating task with ID: %s", req.ID)
		return models.Task{}, err
	}

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	s.log.Info().Msgf("Deleting task with ID: %s", id.String())

	err := s.repo.Delete(ctx, id.String())
	if err != nil {
		s.log.Error().Err(err).Msgf("Error deleting task with ID: %s", id.String())
		return err
	}

	return nil
}

func (s *TaskService) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	s.log.Info().Msg("Listing tasks with filter")

	tasks, err := s.repo.List(ctx, filter)
	if err != nil {
		s.log.Error().Err(err).Msg("Error listing tasks")
		return nil, err
	}

	return tasks, nil
}
