package main

import (
	"flag"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pavelgein/exambot/internal/config"
	"github.com/pavelgein/exambot/internal/input"
	"github.com/pavelgein/exambot/internal/models"
)

type Mode int

const (
	Tasks Mode = iota
	Assignments
)

func ModeFromString(s string) Mode {
	if s == "tasks" {
		return Tasks
	}

	return Assignments
}

type Config struct {
	SourceFile string
	Mode       Mode
	DBConfig   config.DBConfig
}

func MakeConfig() Config {
	sourceFile := flag.String("source-file", "", "File with items to put into database")
	mode := flag.String("mode", "assignments", "assinments | tasks")
	flag.Parse()
	if *sourceFile == "" {
		log.Panicf("-source-file should be set")
	}

	return Config{
		SourceFile: *sourceFile,
		Mode:       ModeFromString(*mode),
		DBConfig: config.DBConfig{
			Dialect:          config.GetEnvWithDefault("EXAMBOT_DB_DIALECT", "sqlite3"),
			ConnectionParams: config.GetEnvWithDefault("EXAMBOT_CONN_PARAMS", "test.db"),
		},
	}
}

func main() {
	config := MakeConfig()

	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.ConnectionParams)
	if err != nil {
		log.Panicf("can not open database: %s", err.Error())
	}

	if os.Getenv("DEBUG") != "" {
		db = db.Debug()
	}

	db.AutoMigrate(&models.Assignment{})
	db.AutoMigrate(&models.Task{})
	db.AutoMigrate(&models.TaskSet{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.TelegramUser{})
	db.AutoMigrate(&models.Course{})

	if config.Mode == Assignments {
		input.InsertAssignments(db, config.SourceFile)
	} else {
		input.InsertTasksFromFile(db, config.SourceFile)
	}

}
