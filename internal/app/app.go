package app

import (
	"fmt"
	"github.com/Vaixle/crud-golang/internal/controller/http/v1"
	"github.com/Vaixle/crud-golang/internal/repository"
	"github.com/Vaixle/crud-golang/internal/usecase"
	"github.com/Vaixle/crud-golang/migration"
	"github.com/Vaixle/crud-golang/pkg/db/postgres"
	"github.com/Vaixle/crud-golang/pkg/httpserver"
	"github.com/Vaixle/crud-golang/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	l := logger.New(viper.GetString("logger.log_level"))

	// Database
	pg, err := postgres.New()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	//Create task table
	migration.InitTodoTable(pg.DB)

	// Repository
	repo := repository.NewTodoRepository(pg.DB)

	// Use case
	translationUseCase := usecase.NewTodoUseCase(repo, l)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, translationUseCase, l)
	httpServer := httpserver.New(handler, httpserver.Port(viper.GetString("http.port")))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
