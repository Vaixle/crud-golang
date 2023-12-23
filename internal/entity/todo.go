package entity

import (
	"github.com/Vaixle/crud-golang/pkg/httpquery"
	"gorm.io/gorm"
)

//go:generate mockgen -source=todo.go -destination=./mocks/mock.go

type Todo struct {
	gorm.Model
	Description string `json:"description" binding:"required" example:"some text"`
	Status      string `json:"status" binding:"required,oneof=open close" example:"open/close"`
}

type TodoRepository interface {
	GetTaskById(id uint) (*Todo, error)
	GetTasks(filters []httpquery.FilterOption, pagination httpquery.Pagination) ([]Todo, error)
	SaveTask(task *Todo) error
}

type TodoUseCase interface {
	GetTaskById(id uint) (*Todo, error)
	GetTasks(filters []httpquery.FilterOption, pagination httpquery.Pagination) ([]Todo, error)
	SaveTask(task *Todo) error
}
