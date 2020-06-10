package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pavelgein/exambot/internal/config"
)

func CreateDBFromEnvironment() (*gorm.DB, error) {
	return gorm.Open(config.GetEnvWithDefault("EXAMBOT_DB_DIALECT", "sqlite3"), config.GetEnvWithDefault("EXAMBOT_CONN_PARAMS", "test.db"))
}
