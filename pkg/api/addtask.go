package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/BulgakovDanil/go_final_project/pkg/constants"
	"github.com/BulgakovDanil/go_final_project/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHendler(w, r)
	default:
		writeJsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// addTaskHandler обработчик для добавленной задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	//Читает тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, "Error reading body", http.StatusBadRequest)
		return
	}

	//Парсим JSON в структуру task
	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Заголовок задачи не должен быть пустым
	if task.Title == "" {
		writeJsonError(w, "an empty parameter 'Title'", http.StatusBadRequest)
		return
	}

	//Проверям и корректируем дату задачи
	err = checkDate(&task)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Сохраняем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Возвращаем ответ
	writeJson(w, map[string]interface{}{"id": id}, http.StatusCreated)

}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, task, http.StatusOK)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {

	//Читает тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, "Error reading body", http.StatusBadRequest)
		return
	}

	//Парсим JSON в структуру task
	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Заголовок задачи не должен быть пустым
	if task.Title == "" {
		writeJsonError(w, "an empty parameter 'Title'", http.StatusBadRequest)
		return
	}

	//Проверям и корректируем дату задачи
	err = checkDate(&task)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Сохраняем задачу в базу данных
	err = db.UpdateTask(&task)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Возвращаем ответ
	writeEmptyJson(w, http.StatusOK)
}

func deleteTaskHendler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	err := db.DeletTask(id)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeEmptyJson(w, http.StatusOK)

}

// checDate функция проверяет и корректирует дату
func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(constants.DateFormat)

	//Если дата не указана, то устанавливаем текущую
	if task.Date == "" {
		task.Date = today
		return nil
	}

	//Парсим дату, проверяем формат
	t, err := time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("invalid 'Date' format: %v", err)
	}

	var next string
	//Если указано правило вычисляем следующую дату
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("incorrect rule: %v", err)
		}
	}

	//Проверямем что дата больше текущей даты
	if afterNow(now, t) {
		if task.Repeat == "" {
			//Если нету правила, то устанавливаем текущую
			task.Date = today
		} else {
			//Иначе берем следующую вычисленную дату
			task.Date = next
		}
	}
	return nil
}

// writeJson Вспомогательная функция для отправки JSON тветов
func writeJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeJsonError Вспомогательная функция для отправки JSON ошибок
func writeJsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// writeEmptyJson Вспомогательная функция для отправки пустых Json ответов
func writeEmptyJson(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, err := w.Write([]byte("{}"))
	if err != nil {
		log.Printf("Error writing empty JSON: %v", err)
	}
}
