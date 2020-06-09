package httpapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/pavelgein/exambot/models"
)

type HttpApi struct {
	DB *gorm.DB
}

func (api *HttpApi) PingHanlder(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Add("Content-Type", "plain/text")

	writer.Write([]byte("ok"))
	// writer.WriteHeader(http.StatusOK)
}

func sendJson(v interface{}, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
	response, err := json.Marshal(v)
	if err != nil {
		log.Printf("error %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(response)
}

func (api *HttpApi) listModel(model interface{}) error {
	return api.DB.Find(model).Error
}

func (api *HttpApi) ListUsers(writer http.ResponseWriter, request *http.Request) {
	users := []models.User{}
	if err := api.listModel(&users); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendJson(users, writer)
}

func (api *HttpApi) ListTasks(writer http.ResponseWriter, request *http.Request) {
	tasks := []models.Task{}
	if err := api.listModel(&tasks); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendJson(tasks, writer)
}

func (api *HttpApi) ListTelegramUsers(writer http.ResponseWriter, request *http.Request) {
	users := []models.TelegramUser{}
	if err := api.listModel(&users); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendJson(users, writer)
}

func (api *HttpApi) ListAssignments(writer http.ResponseWriter, request *http.Request) {
	assignments := []models.Assignment{}
	if err := api.DB.Find(&assignments).Error; err != nil {
		log.Printf("error %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sendJson(assignments, writer)
}
