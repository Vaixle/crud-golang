package repository

import (
	"fmt"
	"github.com/Vaixle/crud-golang/internal/entity"
	"github.com/Vaixle/crud-golang/pkg/httpquery"
	"gorm.io/gorm"
)

var _ entity.TodoRepository = (*TodoRepository)(nil)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) entity.TodoRepository {
	return &TodoRepository{db: db}
}

func (t *TodoRepository) GetTaskById(id uint) (*entity.Todo, error) {
	var todoTask entity.Todo
	if err := t.db.First(&todoTask, id).Error; err != nil {
		return nil, err
	}
	return &todoTask, nil
}

func (t *TodoRepository) GetTasks(filters []httpquery.FilterOption, pagination httpquery.Pagination) ([]entity.Todo, error) {
	var todoTasks []entity.Todo

	query := t.db.Model(&entity.Todo{})

	for _, filter := range filters {
		switch filter.Operator {
		case "gt":
			query = query.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
		case "lt":
			query = query.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
		case "ge":
			query = query.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
		case "le":
			query = query.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
		case "eq":
			query = query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
		case "like":
			query = query.Where(fmt.Sprintf("%s LIKE ?", filter.Field), "%"+filter.Value+"%")
		case "order_by":
			query = query.Order(fmt.Sprintf("%s %s", filter.Field, filter.Value))
		default:
			query = query.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
		}
	}

	query.Offset(pagination.GetOffset()).Limit(pagination.Limit)

	if err := query.Find(&todoTasks).Error; err != nil {
		return nil, err
	}

	return todoTasks, nil
}

func (t *TodoRepository) SaveTask(task *entity.Todo) error {
	if err := t.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}
