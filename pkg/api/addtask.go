package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/BulgakovDanil/go_final_project/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		tasksHandler(w, r)
	default:
		writeJsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

// addTaskHandler обработчик для добавленной задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	//Проверяем что запрос использует метот POST
	if r.Method != http.MethodPost {
		writeJsonError(w, "Medhod not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Читает тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJsonError(w, "Error reading body", http.StatusBadRequest)
		return
	}

	//Отложенное закрытие Body
	defer r.Body.Close()

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

// checDate функция проверяет и корректирует дату
func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format("20060102")

	//Если правило ежедневное, то устанавливаем текущую дату
	if task.Repeat == "d 1" {
		task.Date = today
		return nil
	}

	//Если дата не указана, то устанавливаем текущую
	if task.Date == "" {
		task.Date = today
		return nil
	}

	//Парсим дату, проверяем формат
	t, err := time.Parse("20060102", task.Date)
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
		if len(task.Repeat) == 0 {
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
