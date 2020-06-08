package services

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pavelgein/exambot/models"
)

type AssignmentService interface {
	GetAssignment(user *models.User) *models.Assignment
	Assign(assignment *models.Assignment, timestamp *time.Time)
}

type SimpleAssignmentService struct {
}

func (SimpleAssignmentService) GetAssignment(user *models.User) *models.Assignment {
	return &models.Assignment{Task: models.Task{Content: user.Name}}
}

func (SimpleAssignmentService) Assign(assignment *models.Assignment, timestamp *time.Time) {

}

type FullAssignmentService struct {
	DB *gorm.DB
}

func (service FullAssignmentService) GetAssignment(user *models.User) *models.Assignment {
	var assignment models.Assignment
	res := service.DB.Set("gorm:auto_preload", true).Model(&user).Related(&assignment)
	if res.Error != nil {
		log.Printf("nothing is found, %s", res.Error.Error())
		return nil
	}
	return &assignment
}

func (service FullAssignmentService) Assign(assignment *models.Assignment, timestamp *time.Time) {
	assignment.AssignedAt = timestamp
	service.DB.Save(&assignment)
}
