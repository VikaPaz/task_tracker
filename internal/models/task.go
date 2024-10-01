package models

import "time"

type TaskStatus string

const (
	InProgress TaskStatus = "in_progress"
	Done       TaskStatus = "done"
)

type Task struct {
	ID          string `validate:"omitempty,uuid4"`
	Title       string `validate:"required"`
	Description string
	Created     time.Time
	Updated     time.Time
	Status      TaskStatus `validate:"oneof=in_progress done"`
	OwnerID     string     `validate:"uuid4"`
}

type TaskFilter struct {
	ID          []string `validate:"omitempty,dive,uuid4"`
	Title       string
	Description string
	Status      string `validate:"omitempty,oneof=in_progress done"`
	OwnerID     string `validate:"omitempty,uuid4"`
	// Created     time.Time // after
	// Updated     time.Time // after
	TaskSort TaskSort
	Range    Range
}

type TaskSort struct {
	Field string `validate:"omitempty,oneof=title description status created updated"`
	Order string `validate:"omitempty,oneof=asc desc"`
}

type Range struct {
	Limit  uint64
	Offset uint64
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
