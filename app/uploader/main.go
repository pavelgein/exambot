package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pavelgein/exambot/internal/config"
	"github.com/pavelgein/exambot/internal/db_helpers"
	"github.com/pavelgein/exambot/models"
)

type InputTask struct {
	Task    uint64
	TaskSet string
	Content string
}

type InputTasks []InputTask

type InputItem struct {
	Name    string
	Surname string
	Login   string
	Group   string
	Task    InputTask
	Course  string
}

type InputItems []InputItem

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

func GetTask(db *gorm.DB, item *InputTask) models.Task {

	taskSet := models.TaskSet{
		Name: item.TaskSet,
	}

	if err := db.Where(&taskSet).First(&taskSet).Error; err != nil {
		log.Panicf("can not find taskset with Name = %s", item.TaskSet)
	}

	task := models.Task{
		Number: item.Task,
	}

	if err := db.Model(&taskSet).Where("Number = ?", item.Task).Related(&task).Error; err != nil {
		log.Panicf("can not find task with Number = %d in taskset %s", item.Task, item.TaskSet)
	}

	return task
}

func GetUser(db *gorm.DB, item *InputItem) models.User {
	user := models.User{
		Name:    item.Name,
		Surname: item.Surname,
		Group:   item.Group,
	}

	created, err := db_helpers.GetOrCreate(db, &user)
	if err != nil {
		log.Panicf("error in getting user %s", err.Error())
	}

	if created {
		log.Printf("created user %s %s %s", user.Name, user.Surname, user.Group)
	}

	tguser := models.TelegramUser{
		Login: item.Login,
		User:  user,
	}

	if err := db.Model(&user).Where("Login = ?", item.Login).Related(&tguser).Error; err != nil {
		db.NewRecord(&tguser)
		db.Create(&tguser)
		log.Printf("created tguser %s", tguser.Login)
	}

	return user
}

func GetCourse(db *gorm.DB, item *InputItem) models.Course {
	course := models.Course{
		Name: item.Course,
	}

	created, err := db_helpers.GetOrCreate(db, &course)
	if err != nil {
		log.Panicf("can nont fetch course %s", err.Error())
	}

	if created {
		log.Printf("created course %s", item.Course)
	}

	return course
}

func InsertItem(db *gorm.DB, item *InputItem) {
	task := GetTask(db, &item.Task)
	user := GetUser(db, item)
	course := GetCourse(db, item)
	assignment := models.Assignment{
		Task:   task,
		Course: course,
		User:   user,
	}

	db.NewRecord(&assignment)
	db.Save(&assignment)
}

func Insert(db *gorm.DB, items InputItems) {
	for _, item := range items {
		InsertItem(db, &item)
	}
}

func InsertAssignments(db *gorm.DB, sourceFile string) {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Panicf("can not read from file %s", err.Error())
	}

	var items InputItems

	err = json.Unmarshal(data, &items)
	if err != nil {
		log.Panicf("can not parse source file %s", err.Error())
	}
	log.Printf("need to insert %d items", len(items))
	Insert(db, items)
}

func InsertTask(db *gorm.DB, task *InputTask) {
	taskSet := models.TaskSet{
		Name: task.TaskSet,
	}

	db_helpers.GetOrCreate(db, &taskSet)

	taskModel := models.Task{
		TaskSet: taskSet,
		Number:  task.Task,
		Content: task.Content,
	}

	db.NewRecord(&taskModel)
	db.Create(&taskModel)
}

func InsertTasks(db *gorm.DB, sourceFile string) {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Panicf("can not read from file %s", err.Error())
	}

	var tasks InputTasks

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Panicf("can not parse source file %s", err.Error())
	}
	log.Printf("need to insert %d tasks", len(tasks))
	for _, task := range tasks {
		InsertTask(db, &task)
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
		InsertAssignments(db, config.SourceFile)
	} else {
		InsertTasks(db, config.SourceFile)
	}

}
