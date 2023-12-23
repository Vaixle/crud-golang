package migration

import (
	"github.com/Vaixle/crud-golang/internal/entity"
	"gorm.io/gorm"
	"log"
)

func InitTodoTable(db *gorm.DB) {
	if err := db.AutoMigrate(&entity.Todo{}); err != nil {
		log.Fatal(err.Error())
	}
}
