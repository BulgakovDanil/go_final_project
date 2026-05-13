package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BulgakovDanil/go_final_project/pkg/api"
	"github.com/BulgakovDanil/go_final_project/pkg/db"
)

func main() {

	//Инициализация базы данных
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.DB.Close()

	//Инициализация API маршрутов
	api.Init()

	//Получаем порт из переменной окружения TODO_PORT
	port := os.Getenv("TODO_PORT")
	//Порт по уполчанию, если переменная не указана
	if port == "" {
		port = "7540"
	}

	//Регистрируем обработчик для фронт файлов
	http.Handle("/", http.FileServer(http.Dir("./web")))

	//Запускаем сервер
	fmt.Printf("Сервер запущен на порту %s", port)
	http.ListenAndServe(":"+port, nil)
}
