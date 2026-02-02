package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	//Используем функцию NextDate для расчета следующей даты
	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Возвращаем дату
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))

}

// NextDate вычисляет следующую дату выполнения задачи на основе правила
func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	//Возвращаем ошибку если правило не указано
	if repeat == "" {
		return "", fmt.Errorf("the rule is not specified")
	}

	//Парсим начальную дату
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	//Разбиваем правило на части(тип и параметры)
	rule := strings.Split(repeat, " ")

	//Используем цикл для поиска следующей даты, которая будет после текущей даты(now)
	for {
		//Получаем следующую дату по правилу
		date, err = next(rule, date)
		if err != nil {
			return "", err
		}
		//Проверям,что найденная дата позже текущей
		if afterNow(date, now) {
			break
		}
	}

	//Возвращаем дату в формате "YYYYMMDD"
	return date.Format("20060102"), nil
}

// afterNow функция проверяет, что дата позже текущей даты
func afterNow(date, now time.Time) bool {
	if date.After(now) {
		return true
	}
	return false
}

// next функция вычисляет следующую дату на основе правила
func next(repeat []string, date time.Time) (time.Time, error) {
	if len(repeat) == 0 {
		return time.Time{}, fmt.Errorf("empty repeat rule")
	}
	//Определяем тип правила
	switch repeat[0] {
	case "d": //Ежедневное правило: "d x" - каждые х дней
		if len(repeat) < 2 {
			return time.Time{}, fmt.Errorf("missing number for 'd' rule")
		}

		//Парсим количество дней
		days, err := strconv.Atoi(repeat[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid number: %v", err)
		}

		//Проверям диапозон дней
		if days < 1 || days > 400 {
			return time.Time{}, fmt.Errorf("days should be between 1 and 400")
		}

		//Добовляем дни к текущей дате
		return date.AddDate(0, 0, days), nil

	case "y": //Ежегодное правило: "y" - каждый год в ту же дату
		return date.AddDate(1, 0, 0), nil

	case "w": //Еженедельное правило "w x" - x определенные дни недели
		if len(repeat) < 2 {
			return time.Time{}, fmt.Errorf("missing number for 'w' rule")
		}
		return weekdayRule(repeat[1], date)

	case "m": //Ежемесячное правило: "m x [y]" - в опеделенные дни(х) и месяцы(у) при наличии
		return monthRule(repeat, date)

	default:
		return time.Time{}, fmt.Errorf("unknown rule")
	}
}

// weekdayRule обрабатывает еженедельное правило
func weekdayRule(days string, date time.Time) (time.Time, error) {

	//Разбиваем строрку дней на отдлельные числа
	parts := strings.Split(days, ",")
	if len(parts) == 0 {
		return time.Time{}, fmt.Errorf("no weekdays")
	}

	//Созжаем массив для хранения дней недели
	weekdays := make([]int, 0, len(parts))

	//Парсим каждый день недели
	for _, part := range parts {
		day, err := strconv.Atoi(part)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid weekday number: %v", err)
		}

		//Проверям, что день недели в допустимом диапазоне
		if day < 1 || day > 7 {
			return time.Time{}, fmt.Errorf("weekday must be 1-7")
		}

		weekdays = append(weekdays, day)
	}

	currentDate := date

	//Ищем подходящий день в ближайшие 7 дней
	for i := 0; i < 7; i++ {

		//Переходим к следующему дню
		currentDate = currentDate.AddDate(0, 0, 1)

		//Получаем день недели
		weekday := int(currentDate.Weekday())
		if weekday == 0 { //Преобразуем воскресенье из 0 в 7
			weekday = 7
		}

		//Проверяем совпадение дней недели
		for _, w := range weekdays {
			if w == weekday {
				//Возвращаем результат
				return currentDate, nil
			}
		}
	}
	//Возвращаем ошибку если в течении 7 дней не было совпадений
	return time.Time{}, fmt.Errorf("cannot find suitable weekday")
}

// monthRule обрабатывает ежемесячное правило
func monthRule(rule []string, date time.Time) (time.Time, error) {

	//Проверям структуру правила
	if len(rule) < 2 || len(rule) > 3 {
		return time.Time{}, fmt.Errorf("invalid rule")
	}

	//Разбираем дни месяца
	daysParts := strings.Split(rule[1], ",")
	if len(daysParts) == 0 {
		return time.Time{}, fmt.Errorf("no days")
	}

	//Парсим дни месяца
	days := make([]int, 0, len(daysParts))
	for _, part := range daysParts {
		day, err := strconv.Atoi(part)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid day number: %v", err)
		}

		days = append(days, day)
	}

	//Валидация дней
	for _, day := range days {
		//Допустимые значения от -2 до 31
		if day < -2 || day > 31 {
			return time.Time{}, fmt.Errorf("invalid day")
		}
	}

	//Обрабатываем необязательный третий параметр - месяцы
	var months []int
	if len(rule) == 3 {
		monthParts := strings.Split(rule[2], ",")
		if len(monthParts) == 0 {
			return time.Time{}, fmt.Errorf("no month")
		}

		//Парсим месяцы
		for _, part := range monthParts {
			month, err := strconv.Atoi(part)
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid month number: %v", err)
			}

			months = append(months, month)
		}
	}

	//Вызываем функцию поиска следующей даты по месячному правилу
	return nextMonthDate(date, days, months)
}

// nextMonthDate ищет следующую дату для месячного правила
func nextMonthDate(date time.Time, days, months []int) (time.Time, error) {

	//Начинаем поиск со следующего дня
	currentDate := date.AddDate(0, 0, 1)

	//Ищем подходящий день в ближайший год
	for i := 0; i < 365; i++ {

		//Получаем текущий месяц
		currentMonth := int(currentDate.Month())

		//Если указаны месяцы проверям подходит ли текущий
		if len(months) > 0 {
			monthOK := false
			for _, m := range months {
				if m == currentMonth {
					monthOK = true
					break //Месяц подходит, выходим из проверки
				}
			}
			if !monthOK {
				//Месяц не подходит, переходим к первому дню следующего месяца
				currentDate = time.Date(currentDate.Year(), currentDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location())
				continue
			}
		}

		//Вычисляем количество дней в текущем месяце
		daysInMonth := time.Date(currentDate.Year(), currentDate.Month()+1, 0, 0, 0, 0, 0, currentDate.Location()).Day()

		//Проверям каждый день по правилу
		for _, day := range days {
			finalDay := day
			//Если день отрицательный, вычисляем день с конца месяца
			if day < 0 {
				finalDay = daysInMonth + day + 1
			}

			//Проверяем, что день существует в этом месяце и совпадает с текущей датой
			if finalDay > 0 && finalDay <= daysInMonth && finalDay == currentDate.Day() {
				return currentDate, nil //Возвращаем результат
			}
		}

		//Переходим к следующему дню
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	//Возвращаем ошибку если в течении 365 дней не было совпадений
	return time.Time{}, fmt.Errorf("not found")

}
