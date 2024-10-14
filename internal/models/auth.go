package models

type User struct {
	Id       string `validate:"omitempty,uuid4"`
	Login    string `validate:"required"`
	Password string `validate:"required"`
}
