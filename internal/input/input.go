package input

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pavelgein/exambot/internal/models"
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
	OpenAt  int64
	CloseAt int64
}

type InputItems []InputItem

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
	var user models.User
	db.FirstOrCreate(&user, models.User{
		Name:    item.Name,
		Surname: item.Surname,
		Group:   item.Group,
	})

	var tguser models.TelegramUser
	if res := db.Model(&user).Related(&tguser).Take(&tguser); res.Error == nil {
		log.Printf("found login: %s", tguser.Login)
		if tguser.Login != item.Login {
			log.Panic("wrong login")
		}
	} else if res.RecordNotFound() {
		log.Printf("create user")
		tguser.Login = item.Login
		tguser.User = user
		db.NewRecord(&tguser)
		db.Create(&tguser)
	}

	return user
}

func GetCourse(db *gorm.DB, item *InputItem) models.Course {
	var course models.Course
	log.Printf("course %s", item.Course)
	db.FirstOrCreate(&course, models.Course{
		Name: item.Course,
	})

	return course
}

func InsertItem(db *gorm.DB, item *InputItem) {
	task := GetTask(db, &item.Task)
	user := GetUser(db, item)
	course := GetCourse(db, item)

	openAt := time.Unix(item.OpenAt, 0)
	closeAt := time.Unix(item.CloseAt, 0)
	assignment := models.Assignment{
		Task:    task,
		Course:  course,
		User:    user,
		OpenAt:  &openAt,
		CloseAt: &closeAt,
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
	var taskSet models.TaskSet

	db.FirstOrCreate(&taskSet, models.TaskSet{
		Name: task.TaskSet,
	})

	taskModel := models.Task{
		TaskSet: taskSet,
		Number:  task.Task,
		Content: task.Content,
	}

	db.NewRecord(&taskModel)
	db.Create(&taskModel)
}

func InsertTasks(db *gorm.DB, tasks InputTasks) {
	for _, task := range tasks {
		InsertTask(db, &task)
	}
}

func InsertTasksFromFile(db *gorm.DB, sourceFile string) {
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
