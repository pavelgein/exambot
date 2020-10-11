package services

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pavelgein/exambot/internal/models"
)

func TestTaskSerivce(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.Task{}, &models.Assignment{}, &models.Course{}, &models.TaskSet{})

	user := &models.User{
		Name:    "Name",
		Surname: "Surname",
		Group:   "Group",
	}

	db.NewRecord(user)
	db.Create(user)

	taskSet := &models.TaskSet{
		Name: "TaskSet",
	}

	db.NewRecord(taskSet)
	db.Create(taskSet)

	task := &models.Task{
		Number:  1,
		Content: "content",
		TaskSet: *taskSet,
	}

	db.NewRecord(task)
	db.Create(task)

	course := &models.Course{
		Name: "Course",
	}

	db.NewRecord(course)
	db.Create(course)

	one := time.Unix(1, 0)
	three := time.Unix(3, 0)

	assignment := &models.Assignment{
		Course:  *course,
		Task:    *task,
		User:    *user,
		OpenAt:  &one,
		CloseAt: &three,
	}

	db.NewRecord(assignment)
	db.Create(assignment)

	taskService := FullAssignmentService{DB: db}
	t.Run("in time", func(t *testing.T) {
		now := time.Unix(2, 0)
		assignments := taskService.GetAllAssignments(user, &now)
		if len(assignments) != 1 {
			t.Errorf("Expected len 2, actual %d", len(assignments))
		}
	})

	t.Run("before time", func(t *testing.T) {
		now := time.Unix(0, 0)
		assignments := taskService.GetAllAssignments(user, &now)
		if len(assignments) != 0 {
			t.Errorf("Expected len 0, actual %d", len(assignments))
		}
	})

	t.Run("after time", func(t *testing.T) {
		now := time.Unix(5, 0)
		assignments := taskService.GetAllAssignments(user, &now)
		if len(assignments) != 0 {
			t.Errorf("Expected len 0, actual %d", len(assignments))
		}
	})
}
