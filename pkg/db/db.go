package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128) NOT NULL
);

CREATE INDEX date_id ON scheduler (date)`

var DB *sql.DB

func Init(dbFile string) error {
	//Директория к БД
	dataDir := os.Getenv("TODO_DBFILE")
	if dataDir == "" {
		dataDir = "data"
	}

	//Создаем директорию
	err := os.MkdirAll(dataDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	dataDir = filepath.Clean(dataDir)

	//Путь к файлу БД
	dbPath := filepath.Join(dataDir, dbFile)

	//Проверяем, существует ли файл БД
	_, err = os.Stat(dbPath)
	var install bool
	if err != nil {
		install = true
	}

	//Открывам БД
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	//Проверяем соединение
	err = DB.Ping()
	if err != nil {
		DB.Close()
		DB = nil
		return fmt.Errorf("no connection with db: %v", err)
	}

	//Создаем таблицу если ее нет
	if install == true {
		_, err = DB.Exec(schema)
		if err != nil {
			DB.Close()
			DB = nil
			return fmt.Errorf("error creating the table: %v", err)
		}
	}

	return nil
}
