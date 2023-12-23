// Package postgres implements postgres connection.
package postgres

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	DB *gorm.DB
}

// New -.
func New() (*Postgres, error) {
	dbHost := viper.GetString("db.host")
	dbPort := viper.GetString("db.port")
	dbUser := viper.GetString("db.user")
	dbPass := viper.GetString("db.password")
	dbName := viper.GetString("db.name")
	dbDSN := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)
	fmt.Println(dbDSN)

	gormDB, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, err
	}
	return &Postgres{DB: gormDB}, nil
}
