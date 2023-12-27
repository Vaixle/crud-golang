package v1

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Vaixle/crud-golang/internal/entity"
	mock_entity "github.com/Vaixle/crud-golang/internal/entity/mocks"
	"github.com/Vaixle/crud-golang/pkg/httpquery"
	"github.com/Vaixle/crud-golang/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http/httptest"
	"testing"
	"time"
)

func TestController_CreateTask(t *testing.T) {

	testTable := []struct {
		name               string
		inputBody          string
		inputTodoTask      *entity.Todo
		mockBehavior       func(u *mock_entity.MockTodoUseCase, task *entity.Todo)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "OK",
			inputBody: `{"description":"new Task3","status":"open"}`,
			inputTodoTask: &entity.Todo{
				Description: "new Task3",
				Status:      "open",
			},
			mockBehavior: func(u *mock_entity.MockTodoUseCase, task *entity.Todo) {
				u.EXPECT().SaveTask(task).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"ID":0,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"description":"new Task3","status":"open"}`,
		},
		{
			name:      "ERROR",
			inputBody: `{"description":"new Task3","status":"open"}`,
			inputTodoTask: &entity.Todo{
				Description: "new Task3",
				Status:      "open",
			},
			mockBehavior: func(u *mock_entity.MockTodoUseCase, task *entity.Todo) {
				u.EXPECT().SaveTask(task).Return(errors.New("error save to data base"))
			},
			expectedStatusCode: 400,
			expectedBody:       `{"error":"create task"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoUseCase := mock_entity.NewMockTodoUseCase(ctrl)

			testCase.mockBehavior(todoUseCase, testCase.inputTodoTask)

			r := gin.New()
			l := logger.New("info")
			c := &todoController{l: l, useCase: todoUseCase}

			r.POST("/api/v1/todo", c.createTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/todo", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestController_GetTaskById(t *testing.T) {

	testTable := []struct {
		name               string
		inputId            uint
		mockBehavior       func(u *mock_entity.MockTodoUseCase, id uint)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:    "OK",
			inputId: 3,
			mockBehavior: func(u *mock_entity.MockTodoUseCase, id uint) {
				u.EXPECT().GetTaskById(id).Return(&entity.Todo{
					Model: gorm.Model{
						ID:        3,
						CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
					Description: "new Task3",
					Status:      "open",
				}, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"ID":3,"CreatedAt":"2023-01-01T12:00:00Z","UpdatedAt":"2023-01-01T12:00:00Z","DeletedAt":null,"description":"new Task3","status":"open"}`,
		},
		{
			name:    "ERROR",
			inputId: 1,
			mockBehavior: func(u *mock_entity.MockTodoUseCase, id uint) {
				u.EXPECT().GetTaskById(id).Return(nil, errors.New("id not found"))
			},
			expectedStatusCode: 400,
			expectedBody:       `{"error":"error get task by id 1"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoUseCase := mock_entity.NewMockTodoUseCase(ctrl)

			testCase.mockBehavior(todoUseCase, testCase.inputId)

			r := gin.New()
			l := logger.New("info")
			c := &todoController{l: l, useCase: todoUseCase}

			r.GET("/api/v1/todo/:id", c.getTaskById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/todo/%d", testCase.inputId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestController_GetTasks(t *testing.T) {

	testTable := []struct {
		name                  string
		query                 string
		inputFilterOptions    []httpquery.FilterOption
		inputFilterPagination httpquery.Pagination
		mockBehavior          func(u *mock_entity.MockTodoUseCase, inputFilterOptions []httpquery.FilterOption, inputFilterPagination httpquery.Pagination)
		expectedStatusCode    int
		expectedBody          string
	}{
		{
			name:                  "OK",
			query:                 "?id=gt:0",
			inputFilterOptions:    []httpquery.FilterOption{{Operator: "gt", Field: "id", Value: "0"}},
			inputFilterPagination: httpquery.Pagination{Limit: 100, Page: 0},
			mockBehavior: func(u *mock_entity.MockTodoUseCase, filters []httpquery.FilterOption, pagination httpquery.Pagination) {
				u.EXPECT().GetTasks(filters, pagination).Return([]entity.Todo{{
					Model: gorm.Model{
						ID:        3,
						CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					},
					Description: "new Task3",
					Status:      "open",
				}}, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `[{"ID":3,"CreatedAt":"2023-01-01T12:00:00Z","UpdatedAt":"2023-01-01T12:00:00Z","DeletedAt":null,"description":"new Task3","status":"open"}]`,
		},
		{
			name:                  "ERROR",
			query:                 "?id=gt:0",
			inputFilterOptions:    []httpquery.FilterOption{{Operator: "gt", Field: "id", Value: "0"}},
			inputFilterPagination: httpquery.Pagination{Limit: 100, Page: 0},
			mockBehavior: func(u *mock_entity.MockTodoUseCase, filters []httpquery.FilterOption, pagination httpquery.Pagination) {
				u.EXPECT().GetTasks(filters, pagination).Return(nil, errors.New("some error message"))
			},
			expectedStatusCode: 400,
			expectedBody:       `{"error":"error get tasks"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			todoUseCase := mock_entity.NewMockTodoUseCase(ctrl)

			testCase.mockBehavior(todoUseCase, testCase.inputFilterOptions, testCase.inputFilterPagination)

			r := gin.New()
			l := logger.New("info")
			c := &todoController{l: l, useCase: todoUseCase}

			r.GET("/api/v1/todo", c.getTodoTasks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/todo%s", testCase.query), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedBody, w.Body.String())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
