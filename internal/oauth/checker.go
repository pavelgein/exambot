package oauth

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pavelgein/exambot/models"
)

type OAuthMultiPageChecker struct {
	DB   *gorm.DB
	Salt string
}

func (checker OAuthMultiPageChecker) Check(page *models.Page, token string) bool {
	encrypted := checker.Encrypt(token)

	user := models.ApiUser{}

	if err := checker.DB.Where("Token = ?", encrypted).Take(&user).Error; err != nil {
		return false
	}

	pageRoles := []models.Role{}
	err := checker.DB.Set("gorm:auto_preload", true).Model(&user).Related(&pageRoles, "UserID").Error
	if err != nil {
		log.Printf("error: %s", err.Error())
		return false
	}

	for _, role := range pageRoles {
		if role.Page.ID == page.ID {
			return true
		}
	}
	return false
}

func (checker OAuthMultiPageChecker) Encrypt(token string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(token+checker.Salt)))
}

func (checker OAuthMultiPageChecker) CreateUser(name string, token string) models.ApiUser {
	apiUser := models.ApiUser{
		Name:  name,
		Token: checker.Encrypt(token),
	}

	checker.DB.NewRecord(&apiUser)
	checker.DB.Create(&apiUser)
	return apiUser
}

func (checker OAuthMultiPageChecker) GetUser(name string) (models.ApiUser, error) {
	apiUser := models.ApiUser{}
	return apiUser, checker.DB.Where("Name = ?", name).Take(&apiUser).Error
}

func (checker OAuthMultiPageChecker) GrantPermission(user *models.ApiUser, page *models.Page) error {
	userRoles := []models.Role{}
	err := checker.DB.Set("gorm:auto_preload", true).Model(&user).Related(&userRoles, "UserID").Error
	if err != nil {
		return err
	}

	for _, role := range userRoles {
		if role.Page.ID == page.ID {
			return nil
		}
	}

	role := models.Role{
		Page: *page,
		User: *user,
	}

	checker.DB.NewRecord(&role)
	checker.DB.Create(&role)
	return nil
}

func (checker OAuthMultiPageChecker) GetPage(pageName string) (models.Page, error) {
	page := models.Page{
		Name: pageName,
	}

	return page, checker.DB.Where(&page).Take(&page).Error
}

type OAuthPageChecker struct {
	Parent *OAuthMultiPageChecker
	Page   *models.Page
}

func (checker OAuthPageChecker) Check(token string) bool {
	return checker.Parent.Check(checker.Page, token)
}

func CreatePageChecker(pageName string, parent *OAuthMultiPageChecker) (*OAuthPageChecker, error) {
	page, err := parent.GetPage(pageName)
	if err != nil {
		return nil, err
	}
	return &OAuthPageChecker{
		Parent: parent,
		Page:   &page,
	}, nil
}
