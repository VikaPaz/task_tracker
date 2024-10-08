package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Task struct {
	ID          string `bun:"column:pk,type:uuid,default:uuid_generate_v4()"`
	Title       string `bun:"column:notnull"`
	Description string
	CreatedAt   time.Time `bun:"column:notnull,default:current_timestamp"`
	UpdatedAt   time.Time `bun:"column:nullzero,default:current_timestamp"`
	Status      string    `bun:"column:notnull"`
	OwnerID     string    `bun:"column:notnull,type:uuid"`
}

func modelsTask(task Task) models.Task {
	res := models.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Created:     task.CreatedAt,
		Updated:     task.UpdatedAt,
		Status:      models.TaskStatus(task.Status),
		OwnerID:     task.OwnerID,
	}
	return res
}

func repoTask(task models.Task) Task {
	res := Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.Created,
		UpdatedAt:   task.Updated,
		Status:      task.Status.String(),
		OwnerID:     task.OwnerID,
	}
	return res
}

func (r *TaskRepository) Create(ctx context.Context, task models.Task) (models.Task, error) {
	repoTask := repoTask(task)
	_, err := r.conn.NewInsert().Model(&repoTask).Returning("*").Exec(ctx)
	if err != nil {
		r.log.Error().Err(err).Msgf("can't creating: %v", task)
		return models.Task{}, err
	}
	r.log.Debug().Msgf("maked struct %v", task)

	res := modelsTask(repoTask)
	return res, nil
}

func (r *TaskRepository) Get(ctx context.Context, id uuid.UUID) (models.Task, error) {
	var repoTask Task
	err := r.conn.NewSelect().Model(&repoTask).Where("id = ?", id).Scan(ctx)
	if err != nil {
		r.log.Error().Err(err).Msgf("can't receiving: %v", repoTask)
		if err == sql.ErrNoRows {
			return models.Task{}, models.ErrTaskNotFound
		}
		return models.Task{}, err
	}
	r.log.Debug().Msgf("received struct %v", repoTask)

	res := modelsTask(repoTask)
	return res, nil
}

func (r *TaskRepository) Update(ctx context.Context, req models.Task) (models.Task, error) {
	repoTask := repoTask(req)

	query := r.conn.NewUpdate().
		Model(&repoTask).
		WherePK("id").
		ExcludeColumn("created_at").
		Returning("*")

	if repoTask.Title == "" {
		query.ExcludeColumn("title")
	}
	if repoTask.Description == "" {
		query.ExcludeColumn("description")
	}
	if repoTask.Status == "" {
		query.ExcludeColumn("status")
	}

	res, err := query.Exec(ctx)

	if err != nil {
		r.log.Error().Err(err).Msgf("can't updating: %v", repoTask)
		if err == sql.ErrNoRows {
			return models.Task{}, models.ErrTaskNotFound
		}
		return models.Task{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		r.log.Error().Err(err).Msgf("can't update: %v", repoTask)
		return models.Task{}, err
	}

	if affected != 1 {
		r.log.Error().Err(err).Msgf("can't update: %v", repoTask)
		return models.Task{}, models.ErrTaskNotFound
	}

	resp := modelsTask(repoTask)
	return resp, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	task := &Task{ID: id}
	res, err := r.conn.NewDelete().Model(task).Where("id = ?", id).Exec(ctx)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to delete a task.")
	}

	affected, err := res.RowsAffected()
	if err != nil {
		r.log.Error().Err(err).Msgf("can't delete: %v", task)
		return err
	}
	if affected != 1 {
		r.log.Error().Err(err).Msgf("can't delete: %v", task)
		return models.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) List(ctx context.Context, filter models.TaskFilter) ([]models.Task, error) {
	query := r.conn.NewSelect().Model(&Task{})

	if len(filter.ID) > 0 {
		query = query.Where("id IN (?)", bun.In(filter.ID))
	}

	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}

	if filter.Description != "" {
		query = query.Where("description ILIKE ?", "%"+filter.Description+"%")
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.OwnerID != "" {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}

	var tasks []Task
	err := query.Scan(ctx, &tasks)
	if err != nil {
		r.log.Error().Err(err).Msg("failed to list tasks")
		if err == sql.ErrNoRows {
			return []models.Task{}, models.ErrTaskNotFound
		}
		return nil, err
	}

	res := make([]models.Task, 0, len(tasks))
	for _, val := range tasks {
		res = append(res, modelsTask(val))
	}
	return res, nil
}
