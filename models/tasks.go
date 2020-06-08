package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TaskSet struct {
	gorm.Model
	Name string
}

type Task struct {
	gorm.Model
	Number    uint64
	Content   string
	TaskSetID int
	TaskSet   TaskSet
}

type Course struct {
	gorm.Model
	Name string
}

type Assignment struct {
	gorm.Model
	CourseID   int
	Course     Course
	UserID     int
	User       User
	TaskID     int
	Task       Task
	AssignedAt *time.Time
}
