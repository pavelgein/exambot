package oauth

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pavelgein/exambot/models"
)

func TestOAuth(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.ApiUser{}, &models.Page{}, &models.Role{})

	db = db.Debug()

	page := models.Page{
		Name: "page",
	}

	db.NewRecord(&page)
	db.Create(&page)

	otherPage := models.Page{
		Name: "other page",
	}

	db.NewRecord(&otherPage)
	db.Create(&otherPage)

	checker := OAuthMultiPageChecker{DB: db, Salt: "1"}
	user := checker.CreateUser("user", "token")
	otherUser := checker.CreateUser("user2", "token2")
	checker.GrantPermission(&otherUser, &otherPage)

	t.Run("userGetting", func(t *testing.T) {
		extractedUser, err := checker.GetUser("user")
		if err != nil {
			t.Error("should not be an error")
		}
		if extractedUser.ID != user.ID {
			t.Error("should be extracted same user")
		}

		extractedUser, err = checker.GetUser("user12313")
		if err == nil {
			t.Error("user should not have been found")
		}
	})

	createPageChecker := func(page *models.Page, token string, expected bool) func(*testing.T) {
		return func(t *testing.T) {
			result := checker.Check(page, token)
			if result != expected {
				t.Errorf("Permission with token %s to page %s expected %t", token, page.Name, expected)
			}
		}
	}

	t.Run("beforeGranted", createPageChecker(&page, "token", false))

	checker.GrantPermission(&user, &page)
	t.Run("token2 and otherPage", createPageChecker(&otherPage, "token2", true))
	t.Run("token and page", createPageChecker(&page, "token", true))

	t.Run("spChecker", func(t *testing.T) {
		spChecker, err := CreatePageChecker("page", &checker)
		if err != nil {
			t.Error("can not create page checker")
		}

		if !spChecker.Check("token") {
			t.Error("should have access with token `token`")
		}

		if spChecker.Check("token2") {
			t.Error("should not have access with token `token2`")
		}

		if spChecker.Check("token3") {
			t.Error("unexpected token")
		}
	})
}
