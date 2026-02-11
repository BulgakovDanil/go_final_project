package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/signin", SigninHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", auth(TaskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))
}
