package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pavelgein/exambot/internal/config"
	"github.com/pavelgein/exambot/internal/models"
)

func Init(config *config.DBConfig) (*gorm.DB, error) {
	return gorm.Open(config.Dialect, config.ConnectionParams)
}

func InitWithMigrations(config *config.DBConfig) (*gorm.DB, error) {
	db, err := Init(config)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.ApiUser{}, &models.Page{}, &models.Role{})
	db.AutoMigrate(&models.Task{}, &models.Assignment{}, &models.Course{}, &models.User{}, &models.TelegramUser{}, &models.TaskSet{})
	return db, nil
}

func CreateConfigFromEnvironment() config.DBConfig {
	return config.DBConfig{
		Dialect:          config.GetEnvWithDefault("EXAMBOT_DB_DIALECT", "sqlite3"),
		ConnectionParams: config.GetEnvWithDefault("EXAMBOT_CONN_PARAMS", "test.db"),
	}
}

func CreateDBFromEnvironment() (*gorm.DB, error) {
	config := CreateConfigFromEnvironment()
	return Init(&config)
}
