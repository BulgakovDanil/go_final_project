package api

import (
	"net/http"
	"os"
)

var password string

func Init() {

	password = os.Getenv("TODO_PASSWORD")

	http.HandleFunc("/api/signin", SigninHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", auth(TaskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
}
