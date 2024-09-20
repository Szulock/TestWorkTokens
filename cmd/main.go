package main

import (
	"TokenTestWork/internal/handlers"
	"TokenTestWork/internal/loggerer"

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
	Logger := loggerer.NewLogger()
	Logger.Info("Логер создан")

	http.HandleFunc("/token", handlers.GenerateTokenHandler)

	http.HandleFunc("/refresh", handlers.RefreshTokenHandler)

	Logger.Println("Server started at :8080")
	Logger.Fatal(http.ListenAndServe(":8080", nil))

}
