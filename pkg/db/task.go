package db

import (
	"fmt"
	"time"
)

// Структура таблицы scheduler базы данных
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Addtask добавляет новую задачу в базу данных
func AddTask(task *Task) (int64, error) {
	var id int64

	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?,?,?,?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

// Tasks возвращаетсписок всех задач с ограничением по количеству
func Tasks(limit int) ([]*Task, error) {

	//Устанавливаем значение по умолчанию для limit
	var tasks []*Task
	if limit <= 0 {
		limit = 50
	}

	//SQL запрос с сортировкой по дате
	query := `
		SELECT id, date, title, comment, repeat FROM scheduler
		ORDER BY date ASC, id ASC
		LIMIT ?
	`

	//Выполняем запрос к базе данных
	rows, err := DB.Query(query, limit)
	if err != nil {
		return tasks, fmt.Errorf("query error: %v", err)
	}

	//Отложенное закрытие rows
	defer rows.Close()

	//Интерируем по результатам запроса
	for rows.Next() {
		task := &Task{}

		//Сканируем значения каждой строки в структуру Task
		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)

		if err != nil {
			return tasks, fmt.Errorf("interation error: %v", err)
		}

		//Добавляем задачу
		tasks = append(tasks, task)
	}

	//Если tasks равен nil, возавращаем пустой слайс
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil

}

// SearchText выполняет поиск задач по тексту в заголовке или коментарию
func SearchText(search string, limit int) ([]*Task, error) {

	if limit <= 0 {
		limit = 50
	}

	var tasks []*Task

	//Ищем строку в любом месте
	searchPattern := "%" + search + "%"

	//SQL запрос для поиска по тексту
	query := `SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`

	//Выполняем запрос
	rows, err := DB.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return tasks, fmt.Errorf("query error: %v", err)
	}

	defer rows.Close()

	//Обрабатываем результаты
	for rows.Next() {

		task := &Task{}

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)

		if err != nil {
			return tasks, fmt.Errorf("interation error: %v", err)
		}

		tasks = append(tasks, task)
	}

	//Возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil

}

// SearchDate выполняет поиск задач по конктретной дате
func SearchDate(date string, limit int) ([]*Task, error) {

	if limit <= 0 {
		limit = 50
	}

	var tasks []*Task

	//Парсим дату из формата DD.MM.YYYY
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		return tasks, fmt.Errorf("invalid date format: %v", err)
	}

	//Преобразуем дату в формат YYYYMMDD
	targetDate := t.Format("20060102")

	//SQL запрос для поиска по дате
	query := `SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`

	//Выполняем запрос
	rows, err := DB.Query(query, targetDate, limit)
	if err != nil {
		return tasks, fmt.Errorf("query error: %v", err)
	}

	defer rows.Close()

	//Обрабатываем результат
	for rows.Next() {

		task := &Task{}

		err := rows.Scan(
			&task.ID,
			&task.Date,
			&task.Title,
			&task.Comment,
			&task.Repeat,
		)

		if err != nil {
			return tasks, fmt.Errorf("interation error: %v", err)
		}

		tasks = append(tasks, task)
	}

	//Возвращаем пустой слайс вместо nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil

}

// GetTask извлекает задачу из базы данных по id
func GetTask(id string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("empty id")
	}

	task := &Task{}

	err := DB.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("query error")
	}

	return task, nil
}

// UpdateTask Обнавляет существующую задачу в базе данных
func UpdateTask(task *Task) error {

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil
}

// DeletTask удаляет задачу с базы данных
func DeletTask(id string) error {

	result, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func UpdateDate(task *Task, date string) error {

	res, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	return nil

}
