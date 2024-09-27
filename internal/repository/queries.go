package repository

import (
	"context"
	"time"

	"github.com/VikaPaz/task_tracker/internal/models"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Task struct {
	ID          uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	Title       string    `bun:",notnull"`
	Description string
	Created     time.Time `bun:",notnull,default:current_timestamp"`
	Updated     time.Time `bun:",nullzero,default:current_timestamp"`
	Status      string    `bun:",notnull" validate:"oneof=in_progress done"`
	OwnerID     string    `bun:",notnull" validate:"uuid4"`
}

func (r *TaskRepository) Create(ctx context.Context, title string, description string) (models.Task, error) {
	task := Task{
		Title:       title,
		Description: description,
		Created:     time.Now(),
		Status:      "in_progress",
	}
	_, err := r.conn.NewInsert().Model(&task).Returning("*").Exec(ctx)
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating a task")
		return models.Task{}, err
	}

	res := models.Task{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Created:     task.Created,
		Updated:     task.Updated,
		Status:      task.Status,
		OwnerID:     task.OwnerID,
	}
	return res, nil
}

func (r *TaskRepository) Get(ctx context.Context, id int64) (models.Task, error) {
	var task Task
	err := r.conn.NewSelect().Model(&task).Where("id = ?", id).Scan(ctx)
	if err != nil {
		r.log.Error().Err(err).Msg("Error receiving a task.")
		return models.Task{}, err
	}

	res := models.Task{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Created:     task.Created,
		Updated:     task.Updated,
		Status:      task.Status,
		OwnerID:     task.OwnerID,
	}
	return res, nil
}

func (r *TaskRepository) Update(ctx context.Context, req models.Task) (models.Task, error) {
	task := Task{
		ID:          uuid.MustParse(req.ID),
		Title:       req.Title,
		Description: req.Description,
		Updated:     time.Now(),
		Status:      req.Status,
	}
	_, err := r.conn.NewUpdate().Model(req).Where("id = ?", task.ID.String()).Exec(ctx)
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating a task.")
	}

	res := models.Task{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Created:     task.Created,
		Updated:     task.Updated,
		Status:      task.Status,
		OwnerID:     task.OwnerID,
	}
	return res, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	task := &Task{ID: id}
	_, err := r.conn.NewDelete().Model(task).Where("id = ?", id).Exec(ctx)
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting a task.")
	}
	return err
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

	if filter.TaskSort.Field != "" && filter.TaskSort.Order != "" {
		query = query.OrderExpr("? ?", bun.Ident(filter.TaskSort.Field), bun.Ident(filter.TaskSort.Order))
	}

	if filter.Range.Limit > 0 {
		query = query.Limit(int(filter.Range.Limit))
	}
	if filter.Range.Offset > 0 {
		query = query.Offset(int(filter.Range.Offset))
	}

	var tasks []Task
	err := query.Scan(ctx, &tasks)
	if err != nil {
		r.log.Error().Err(err).Msg("Error listing tasks")
		return nil, err
	}

	res := make([]models.Task, 0, len(tasks))
	for _, val := range tasks {
		res = append(res, models.Task{
			ID:          val.ID.String(),
			Title:       val.Title,
			Description: val.Description,
			Created:     val.Created,
			Updated:     val.Updated,
			Status:      val.Status,
			OwnerID:     val.OwnerID,
		},
		)
	}
	return res, nil
}
