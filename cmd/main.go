package main

import (
	"TokenTestWork/internal/logger"
	"TokenTestWork/internal/middleware"
	"net/http"
)

//
//Как будет строиться сервис:
//Сначала организую файлы сервиса
//Далее я инициализирую логгер для отслеживания всех событий
//Потом я напишу обработчики маршрутов
//Потом я напишу функцию создания токена Access, потом Refresh
//Потом я пишу логику обновления токенов
//Потом я напишу логику уведомления по email
//Потом Контейнеризирую

func main() {
	logger := logger.NewLogger()
	logger.Info("Логер создан")

	http.HandleFunc("/newToken", middleware.CreateAccessToken)
	http.HandleFunc("/newToken", middleware.RefreshToken)

	http.ListenAndServe(":8080", nil)
}
