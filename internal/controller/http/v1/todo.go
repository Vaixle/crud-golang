package v1

import (
	"fmt"
	"github.com/Vaixle/crud-golang/internal/entity"
	"github.com/Vaixle/crud-golang/pkg/httpquery"
	"github.com/Vaixle/crud-golang/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type todoController struct {
	l       logger.Interface
	useCase entity.TodoUseCase
}

func newTODORoutes(handler *gin.RouterGroup, useCase entity.TodoUseCase, l logger.Interface) {
	r := &todoController{l: l, useCase: useCase}

	h := handler.Group("/todo")
	{
		h.GET("", r.getTodoTasks)
		h.GET("/:id", r.getTaskById)
		h.POST("", r.createTask)
	}
}

// @Summary      Get todo tasks
// @Description  Get todo tasks
// @Tags         todo
// @Produce      json
// @Param        filedName1    query     string  false  "greater than" example(gt:1)
// @Param        filedName2   query     string  false  "lower than" example(lt:1)
// @Param        filedName3    query     string  false  "greater and equal than" example(ge:1)
// @Param        filedName4    query     string  false  "lower and equal than" example(le:1)
// @Param        filedName5    query     string  false  "equal" example(eq:1)
// @Param        filedName6    query     string  false  "not equal than" example(ne:1)
// @Param        filedName7    query     string  false  "greater than" example(order_by:1)
// @Param        filedName8    query     string  false  "like something" example(like:1)
// @Param        page    query     string  false  "page" example(2)
// @Param        limit    query     string  false  "limit" example(3)
// @Success      200  {array}   entity.Todo
// @Failure 400 {string} string "{"error": "some error message"}"
// @Router       /todo [get]
func (t *todoController) getTodoTasks(gc *gin.Context) {
	filterOptions, pagination, err := httpquery.ParseQueryParams(gc.Request.URL.Query())
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error query params",
		})
		return
	}

	tasks, err := t.useCase.GetTasks(filterOptions, pagination)
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error get tasks",
		})
		return
	}

	gc.JSON(http.StatusOK, &tasks)
}

// @Summary      Get todo task
// @Description  Get todo task by id
// @Tags         todo
// @Produce      json
// @Param        id    path      int  true "Todo task ID"
// @Success      200  {object}   entity.Todo
// @Failure 400 {string} string "{"error": "some error message"}"
// @Router       /todo/{id} [get]
func (t *todoController) getTaskById(gc *gin.Context) {
	idStr := gc.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		t.l.Error(err, "http - v1 - get task")
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "error id param",
		})
	}

	task, err := t.useCase.GetTaskById(uint(id))
	if err != nil {
		t.l.Error(err, "http - v1 - get task")
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("error get task by id %d", id),
		})
		return
	}

	gc.JSON(http.StatusOK, &task)
}

// @Summary      Create todo task
// @Description  Create todo task
// @Tags         todo
// @Accept       json
// @Produce      json
// @Success      200  {object}   entity.Todo
// @Failure 400 {string} string "{"error": "some error message"}"
// @Router       /todo [post]
func (t *todoController) createTask(gc *gin.Context) {
	var task entity.Todo
	if err := gc.ShouldBindJSON(&task); err != nil {
		t.l.Error(err, "http - v1 - create task")
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := t.useCase.SaveTask(&task); err != nil {
		t.l.Error(err, "http - v1 - save task")
		gc.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "create task",
		})
		return
	}

	gc.JSON(http.StatusOK, &task)
}
