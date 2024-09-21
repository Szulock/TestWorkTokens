package storage

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// подключаюсь к базе данных
func GetDb(logger *logrus.Logger) *sql.DB {
	connStr := "user=postgres password=szulock dbname=TestForWorkToken sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error("Не получилось подключиться к базе данных:", err)
		return nil
	}

	// Проверка подключения
	if err = db.Ping(); err != nil {
		logger.Error("Ошибка при проверке подключения к базе данных:", err)
		return nil
	}

	logger.Info("Успешное подключение к базе данных")
	return db
}
