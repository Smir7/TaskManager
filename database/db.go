package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgx.Conn, error) {
	// Загружаем переменные из .env файла
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbURL := "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") + " password=" +
		os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME")

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу, если она не существует
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT,
        status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
        created_at TIMESTAMP DEFAULT now(),
        updated_at TIMESTAMP DEFAULT now()
    );`

	_, err = db.Exec(context.Background(), createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
