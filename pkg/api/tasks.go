package api

import (
	"net/http"
	"time"

	"github.com/BulgakovDanil/go_final_project/pkg/constants"
	"github.com/BulgakovDanil/go_final_project/pkg/db"
)

// TasksResp структура для JSON ответа списка задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обработчик HTTP запросов для получения списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {

	//Проверяем метод запроса
	if r.Method != http.MethodGet {
		writeJsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Получаем параметр search
	search := r.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	//Определяем тип поиска
	switch {
	case search == "":
		//Если параметр search пустой - возвращаем все задачи
		tasks, err = db.Tasks(constants.TasksLimit)

	case isDate(search):
		//Если serch соотвествует формату даты - поиск по дате
		tasks, err = db.SearchDate(search, constants.TasksLimit)

	default:
		//Если search не пустой и не дата - поиск по тексту
		tasks, err = db.SearchText(search, constants.TasksLimit)
	}

	//Если при выполнении запроса к БД произошла ошибка
	if err != nil {
		//Отправляем ошибку в JSON формате
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
	}

	//Формируем структуру отвеа
	resp := TasksResp{Tasks: tasks}

	//Отправляем ответ
	writeJson(w, resp, http.StatusOK)
}

// isDate проверяет, является ли строка датой
func isDate(search string) bool {

	//Парсим строку как дату
	_, err := time.Parse("02.01.2006", search)
	//если ошибки нет значит это дата
	return err == nil
}
