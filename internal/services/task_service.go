package services

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pavelgein/exambot/internal/models"
)

type AssignmentService interface {
	GetAssignment(user *models.User) *models.Assignment
	Assign(assignment *models.Assignment, timestamp *time.Time)
	GetAllAssignments(user *models.User) []models.Assignment
	AssignMany(assignments []models.Assignment, timestamp *time.Time)
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

func (service FullAssignmentService) GetAllAssignments(user *models.User) []models.Assignment {
	assignments := []models.Assignment{}
	if err := service.DB.Set("gorm:auto_preload", true).Model(user).Related(&assignments).Error; err != nil {
		log.Printf("error: %s", err.Error())
	}

	return assignments
}

func (service FullAssignmentService) Assign(assignment *models.Assignment, timestamp *time.Time) {
	assignment.AssignedAt = timestamp
	service.DB.Save(&assignment)
}

func (service FullAssignmentService) AssignMany(assignments []models.Assignment, timestamp *time.Time) {
	for i := 0; i != len(assignments); i++ {
		service.DB.Model(&assignments[i]).Update("AssignedAt", timestamp)
	}
}
