package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pavelgein/exambot/packages/httpapi"
)

func main() {
	config := CreateConfigFromEnv()

	db, err := gorm.Open(config.DB.Dialect, config.DB.ConnectionParams)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db = db.Set("gorm:auto_preload", true)

	if os.Getenv("DEBUG") != "" {
		db = db.Debug()
	}

	api := httpapi.HttpApi{DB: db}

	router := mux.NewRouter()
	router.HandleFunc("/ping", api.PingHanlder)
	router.HandleFunc("/list/users", api.ListUsers)
	router.HandleFunc("/list/tasks", api.ListTasks)
	router.HandleFunc("/list/tgusers", api.ListTelegramUsers)
	router.HandleFunc("/list/assignments", api.ListAssignments)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Server.Address,
		WriteTimeout: config.Server.WriteTimeout,
		ReadTimeout:  config.Server.ReadTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
