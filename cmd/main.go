package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BulgakovDanil/go_final_project/pkg/db"
)

func main() {

	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	http.Handle("/", http.FileServer(http.Dir("../web")))

	fmt.Println("Сервер запущен")
	http.ListenAndServe(":"+port, nil)
}
