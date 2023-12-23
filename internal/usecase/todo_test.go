package usecase

import (
	"crud-golang/internal/entity"
	mock_entity "crud-golang/internal/entity/mocks"
	"crud-golang/pkg/httpquery"
	"crud-golang/pkg/logger"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestTodoUseCase_GetTaskById(t *testing.T) {
	testTable := []struct {
		inputId        uint
		mockBehaviour  func(r *mock_entity.MockTodoRepository, id uint)
		expectedEntity *entity.Todo
		wantErr        bool
	}{
		{
			inputId: 1,
			mockBehaviour: func(r *mock_entity.MockTodoRepository, id uint) {
				r.EXPECT().GetTaskById(id).Return(&entity.Todo{Model: gorm.Model{ID: 1}}, nil)
			},
			expectedEntity: &entity.Todo{Model: gorm.Model{ID: 1}},
			wantErr:        false,
		},
		{
			inputId: 1,
			mockBehaviour: func(r *mock_entity.MockTodoRepository, id uint) {
				r.EXPECT().GetTaskById(id).Return(nil, errors.New("some error"))
			},
			expectedEntity: nil,
			wantErr:        true,
		},
	}

	for _, testCase := range testTable {
		ctrl := gomock.NewController(t)

		repo := mock_entity.NewMockTodoRepository(ctrl)
		l := logger.New("info")

		useCase := NewTodoUseCase(repo, l)

		testCase.mockBehaviour(repo, testCase.inputId)

		task, err := useCase.GetTaskById(testCase.inputId)

		if testCase.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, task, testCase.expectedEntity)
	}
}

func TestTodoUseCase_GetTasks(t *testing.T) {
	testTable := []struct {
		name                  string
		inputFilterOptions    []httpquery.FilterOption
		inputFilterPagination httpquery.Pagination
		mockBehaviour         func(r *mock_entity.MockTodoRepository, inputFilterOptions []httpquery.FilterOption, inputFilterPagination httpquery.Pagination)
		expectedEntities      []entity.Todo
		wantErr               bool
	}{
		{
			name:                  "OK",
			inputFilterOptions:    []httpquery.FilterOption{{Operator: "gt", Field: "id", Value: "0"}},
			inputFilterPagination: httpquery.Pagination{Limit: 100, Page: 0},
			mockBehaviour: func(r *mock_entity.MockTodoRepository, filters []httpquery.FilterOption, pagination httpquery.Pagination) {
				r.EXPECT().GetTasks(filters, pagination).Return(
					[]entity.Todo{{
						Model: gorm.Model{
							ID:        3,
							CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
							UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
						},
						Description: "new Task3",
					}}, nil)
			},
			expectedEntities: []entity.Todo{{
				Model: gorm.Model{
					ID:        3,
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "new Task3",
			}},
			wantErr: false,
		},
		{
			name:                  "ERROR",
			inputFilterOptions:    []httpquery.FilterOption{{Operator: "gt", Field: "id", Value: "0"}},
			inputFilterPagination: httpquery.Pagination{Limit: 100, Page: 0},
			mockBehaviour: func(r *mock_entity.MockTodoRepository, filters []httpquery.FilterOption, pagination httpquery.Pagination) {
				r.EXPECT().GetTasks(filters, pagination).Return(nil, errors.New("some error"))
			},
			expectedEntities: nil,
			wantErr:          true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mock_entity.NewMockTodoRepository(ctrl)

			l := logger.New("info")

			useCase := NewTodoUseCase(repo, l)

			testCase.mockBehaviour(repo, testCase.inputFilterOptions, testCase.inputFilterPagination)

			tasks, err := useCase.GetTasks(testCase.inputFilterOptions, testCase.inputFilterPagination)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tasks, testCase.expectedEntities)
		})
	}
}

func TestTodoUseCase_SaveTask(t *testing.T) {
	testTable := []struct {
		name           string
		inputEntity    *entity.Todo
		mockBehavior   func(r *mock_entity.MockTodoRepository, task *entity.Todo)
		expectedEntity *entity.Todo
		wantErr        bool
	}{
		{
			name: "OK",
			inputEntity: &entity.Todo{
				Description: "new Task3",
			},
			mockBehavior: func(r *mock_entity.MockTodoRepository, task *entity.Todo) {
				r.EXPECT().SaveTask(task).Return(nil)
			},
			expectedEntity: &entity.Todo{
				Description: "new Task3",
			},
			wantErr: false,
		},
		{
			name:        "ERROR",
			inputEntity: nil,
			mockBehavior: func(r *mock_entity.MockTodoRepository, task *entity.Todo) {
				r.EXPECT().SaveTask(nil).Return(errors.New("some error"))
			},
			expectedEntity: nil,
			wantErr:        true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mock_entity.NewMockTodoRepository(ctrl)

			l := logger.New("info")

			useCase := NewTodoUseCase(repo, l)

			testCase.mockBehavior(repo, testCase.inputEntity)

			err := useCase.SaveTask(testCase.inputEntity)

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.inputEntity, testCase.expectedEntity)
		})
	}
}
