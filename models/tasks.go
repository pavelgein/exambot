package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Number  uint64
	Content string
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
