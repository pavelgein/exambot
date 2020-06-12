package main

import (
	"flag"
	"log"
	"os"

	"github.com/pavelgein/exambot/internal/config"
	"github.com/pavelgein/exambot/internal/db"
	"github.com/pavelgein/exambot/internal/input"
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
		DBConfig:   db.CreateConfigFromEnvironment(),
	}
}

func main() {
	config := MakeConfig()

	db, err := db.InitWithMigrations(&config.DBConfig)
	if err != nil {
		log.Panicf("can not open database: %s", err.Error())
	}
	defer db.Close()

	if os.Getenv("DEBUG") != "" {
		db = db.Debug()
	}

	if config.Mode == Assignments {
		input.InsertAssignments(db, config.SourceFile)
	} else {
		input.InsertTasksFromFile(db, config.SourceFile)
	}

}
