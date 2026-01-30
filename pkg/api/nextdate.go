package api

import (
	"net/http"
	"time"

	"github.com/BulgakovDanil/go_final_project/pkg/service"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	//Проверяем метот запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Извлекаем параметры из запроса
	//Формат запроса: "/api/nextdate?now=<20060102>&date=<20060102>&repeat=<правило>"
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	//Валидация параметров
	if date == "" {
		http.Error(w, "missing 'date'", http.StatusBadRequest)
		return
	}
	if repeat == "" {
		http.Error(w, "missing 'repeat'", http.StatusBadRequest)
		return
	}

	var now time.Time

	//Обработка параметра 'now'
	if nowStr != "" {
		//Если передан, парсим его
		parsedNow, err := time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "invalid 'now' format ", http.StatusBadRequest)
			return
		}
		now = parsedNow
	} else {
		//Если 'now' не указан, устанавливаем текущую дату
		now = time.Now()
	}

	//Используем функцию NextDate из пакета service для расчета следующей даты
	nextDate, err := service.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Возвращаем дату
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}
