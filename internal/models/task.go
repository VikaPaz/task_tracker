package models

import "time"

type TaskStatus string

const (
	InProgress TaskStatus = "in_progress"
	Done       TaskStatus = "done"
)

type Task struct {
	ID          string `validate:"omitempty,uuid4"`
	Title       string
	Description string
	Created     time.Time
	Updated     time.Time
	Status      TaskStatus `validate:"omitempty,oneof=in_progress done"`
	OwnerID     string     `validate:"uuid4"`
}

type TaskFilter struct {
	ID          []string `form:"id" validate:"omitempty,dive,uuid4"`
	Title       string   `form:"title"`
	Description string   `form:"description"`
	Status      string   `form:"status" validate:"omitempty,oneof=in_progress done"`
	OwnerID     string   `form:"owner_id" validate:"omitempty,uuid4"`
}

type TaskUpdate struct {
	ID          string `validate:"uuid4"`
	Title       string
	Description string
	Status      string `validate:"required,oneof=in_progress done"`
	OwnerID     string `validate:"omitempty,uuid4"`
}

func (t TaskStatus) String() string {
	return string(t)
}
