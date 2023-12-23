package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vaixle/crud-golang/internal/entity"
	"github.com/Vaixle/crud-golang/pkg/httpquery"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"testing"
	"time"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatal(err)
	}
	return sqldb, gormdb, mock
}

func Test_GetTaskById(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()

	repo := NewTodoRepository(db)

	type args struct {
		id uint
	}

	testTable := []struct {
		name           string
		args           args
		mockBehavior   func(args args)
		expectedEntity *entity.Todo
		wantErr        bool
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "updated_at", "created_at", "deleted_at", "description"}).AddRow(
					1, time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					nil,
					"New Task1",
				)
				mock.ExpectQuery(`SELECT (.+) FROM "todos" WHERE "todos"."id" = (.+)`).
					WithArgs(args.id).WillReturnRows(rows)
			},
			expectedEntity: &entity.Todo{
				Model: gorm.Model{
					ID:        1,
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "New Task1",
			},
			wantErr: false,
		},
		{
			name: "ERROR",
			args: args{
				id: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM "todos" WHERE "todos"."id" = (.+)`).
					WithArgs(args.id).WillReturnError(errors.New("some error"))
			},
			expectedEntity: &entity.Todo{
				Model: gorm.Model{
					ID:        1,
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "New Task1",
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			task, err := repo.GetTaskById(testCase.args.id)
			assert.Nil(t, mock.ExpectationsWereMet())

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.expectedEntity, task)
			}
		})
	}
}

func TestTodoRepository_GetTasks(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewTodoRepository(db)

	type args struct {
		id uint
	}

	testTable := []struct {
		name                  string
		args                  args
		inputFilterOptions    []httpquery.FilterOption
		inputFilterPagination httpquery.Pagination
		mockBehavior          func(args args)
		expectedEntity        []entity.Todo
		wantErr               bool
	}{
		{
			name: "OK",
			args: args{
				id: 1,
			},
			inputFilterOptions:    []httpquery.FilterOption{{Operator: "gt", Field: "id", Value: "1"}},
			inputFilterPagination: httpquery.Pagination{Limit: 100, Page: 0},
			mockBehavior: func(args args) {
				rows := mock.NewRows([]string{"id", "updated_at", "created_at", "deleted_at", "description"}).AddRow(
					1, time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					nil,
					"New Task1",
				)
				mock.ExpectQuery(`SELECT (.+) FROM "todos" WHERE id > (.+)`).
					WithArgs(strconv.Itoa(int(args.id))).WillReturnRows(rows)
			},
			expectedEntity: []entity.Todo{{
				Model: gorm.Model{
					ID:        1,
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "New Task1",
			}},
			wantErr: false,
		},
		{
			name: "ERROR",
			args: args{
				id: 1,
			},
			inputFilterOptions:    []httpquery.FilterOption{{Operator: "gt", Field: "id", Value: "1"}},
			inputFilterPagination: httpquery.Pagination{Limit: 100, Page: 0},
			mockBehavior: func(args args) {
				mock.ExpectQuery(`SELECT (.+) FROM "todos" WHERE id > (.+)`).
					WithArgs(strconv.Itoa(int(args.id))).WillReturnError(errors.New("some error"))
			},
			expectedEntity: nil,
			wantErr:        true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args)

			task, err := repo.GetTasks(testCase.inputFilterOptions, testCase.inputFilterPagination)
			assert.Nil(t, mock.ExpectationsWereMet())

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.expectedEntity, task)
			}
		})
	}
}

func TestTodoRepository_SaveTask(t *testing.T) {
	sqlDB, db, mock := DbMock(t)
	defer sqlDB.Close()
	repo := NewTodoRepository(db)

	type args struct {
		id uint
	}

	testTable := []struct {
		name           string
		inputEntity    *entity.Todo
		mockBehavior   func(inputEntity *entity.Todo)
		expectedEntity *entity.Todo
		wantErr        bool
	}{
		{
			name: "OK",
			inputEntity: &entity.Todo{
				Model: gorm.Model{
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "New Task1",
			},
			mockBehavior: func(inputEntity *entity.Todo) {
				rows := mock.NewRows([]string{"id", "updated_at", "created_at", "deleted_at", "description", "status"}).AddRow(1, inputEntity.CreatedAt,
					inputEntity.UpdatedAt,
					nil,
					"New Task1",
					"open",
				)
				expectedSQL := `INSERT INTO "todos" (.+) VALUES (.+)`
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(rows).
					WithArgs(inputEntity.CreatedAt, inputEntity.UpdatedAt, inputEntity.DeletedAt, inputEntity.Description,
						inputEntity.Status)
				mock.ExpectCommit()

			},
			expectedEntity: &entity.Todo{
				Model: gorm.Model{
					ID:        1,
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "New Task1",
				Status:      "open",
			},
			wantErr: false,
		},
		{
			name: "ERROR",
			inputEntity: &entity.Todo{
				Model: gorm.Model{
					UpdatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
				},
				Description: "New Task1",
				Status:      "open",
			},
			mockBehavior: func(inputEntity *entity.Todo) {
				expectedSQL := `INSERT INTO "todos" (.+) VALUES (.+)`
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).
					WithArgs(inputEntity.CreatedAt, inputEntity.UpdatedAt, inputEntity.DeletedAt, inputEntity.Description, inputEntity.Status).
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()

			},
			expectedEntity: nil,
			wantErr:        true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.inputEntity)

			err := repo.SaveTask(testCase.inputEntity)
			assert.Nil(t, mock.ExpectationsWereMet())

			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.expectedEntity, testCase.inputEntity)
			}
		})
	}
}
