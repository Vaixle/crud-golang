package usecase

import (
	"github.com/Vaixle/crud-golang/internal/entity"
	"github.com/Vaixle/crud-golang/pkg/httpquery"
	"github.com/Vaixle/crud-golang/pkg/logger"
)

var _ entity.TodoUseCase = (*TodoUseCase)(nil)

type TodoUseCase struct {
	repo entity.TodoRepository
	l    logger.Interface
}

func NewTodoUseCase(repo entity.TodoRepository, l logger.Interface) entity.TodoUseCase {
	return &TodoUseCase{repo: repo, l: l}
}

func (t TodoUseCase) GetTaskById(id uint) (*entity.Todo, error) {
	task, err := t.repo.GetTaskById(id)
	if err != nil {
		return task, err
	}

	t.l.Info("returning task by id")
	return task, nil
}

func (t TodoUseCase) GetTasks(filters []httpquery.FilterOption, pagination httpquery.Pagination) ([]entity.Todo, error) {
	tasks, err := t.repo.GetTasks(filters, pagination)
	if err != nil {
		return tasks, err
	}

	t.l.Info("returning tasks")
	return tasks, nil
}

func (t TodoUseCase) SaveTask(task *entity.Todo) error {
	if err := t.repo.SaveTask(task); err != nil {
		return err
	}

	t.l.Info("success creating task")
	return nil
}
