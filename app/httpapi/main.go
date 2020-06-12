package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/pavelgein/exambot/internal/db"
	"github.com/pavelgein/exambot/internal/oauth"
	"github.com/pavelgein/exambot/packages/httpapi"
)

type ProtectedPages struct {
	Checker *oauth.OAuthMultiPageChecker
}

func (pages *ProtectedPages) MakeHandler(pageName string, handler http.HandlerFunc) http.HandlerFunc {
	checker, err := oauth.CreatePageChecker(pageName, pages.Checker)
	if err != nil {
		panic(err)
	}

	m := httpapi.OAuthMiddleware{
		Checker: checker,
	}

	return m.Wrap(handler)
}

func main() {
	config := CreateConfigFromEnv()

	db, err := db.InitWithMigrations(&config.DB)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	checker := oauth.OAuthMultiPageChecker{
		DB:   db,
		Salt: config.Salt,
	}

	pages := ProtectedPages{Checker: &checker}

	db = db.Set("gorm:auto_preload", true)

	if os.Getenv("DEBUG") != "" {
		db = db.Debug()
	}

	api := httpapi.HttpApi{DB: db}

	router := mux.NewRouter()
	router.HandleFunc("/ping", api.PingHanlder)
	router.HandleFunc("/list/users", pages.MakeHandler("list/users", api.ListUsers))
	router.HandleFunc("/list/tasks", pages.MakeHandler("list/takss", api.ListTasks))
	router.HandleFunc("/list/tgusers", pages.MakeHandler("list/tgusrs", api.ListTelegramUsers))
	router.HandleFunc("/list/assignments", pages.MakeHandler("list/assignments", api.ListAssignments))

	router.HandleFunc("/create/tasks", pages.MakeHandler("create/tasks", api.InputTask))
	router.HandleFunc("/create/assignments", pages.MakeHandler("create/assignments", api.InputAssignments))

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Server.Address,
		WriteTimeout: config.Server.WriteTimeout,
		ReadTimeout:  config.Server.ReadTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
