package services

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pavelgein/exambot/internal/models"
)

type TestCase struct {
	Name           string
	Time           time.Time
	ExpectedLength int
}

type TestCases []TestCase

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
	five := time.Unix(5, 0)

	assignment := &models.Assignment{
		Course:  *course,
		Task:    *task,
		User:    *user,
		OpenAt:  &one,
		CloseAt: &three,
	}

	db.NewRecord(assignment)
	db.Create(assignment)

	secondAssignment := &models.Assignment{
		Course:  *course,
		Task:    *task,
		User:    *user,
		OpenAt:  &one,
		CloseAt: &five,
	}

	db.NewRecord(secondAssignment)
	db.Create(secondAssignment)

	testCases := TestCases{
		TestCase{
			Name:           "in time",
			Time:           time.Unix(2, 0),
			ExpectedLength: 2,
		},

		TestCase{
			Name:           "in time2",
			Time:           time.Unix(4, 0),
			ExpectedLength: 1,
		},

		TestCase{
			Name:           "before",
			Time:           time.Unix(0, 0),
			ExpectedLength: 0,
		},

		TestCase{
			Name:           "after",
			Time:           time.Unix(10, 0),
			ExpectedLength: 0,
		},
	}

	taskService := FullAssignmentService{DB: db}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			assignments := taskService.GetAllAssignments(user, &testCase.Time)
			if len(assignments) != testCase.ExpectedLength {
				t.Errorf("expected length %d, actual %d", testCase.ExpectedLength, len(assignments))
			}
		})
	}
}
