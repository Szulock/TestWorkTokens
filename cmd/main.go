package main

import (
	"TokenTestWork/internal/loggerer"
	"TokenTestWork/internal/middleware"

	"net/http"
)

func main() {
	Logger := loggerer.NewLogger()
	Logger.Info("Логер создан")

	http.HandleFunc("/token", middleware.GenerateTokenHandler)

	http.HandleFunc("/refresh", middleware.RefreshTokenHandler)

	Logger.Println("Server started at :8080")
	Logger.Fatal(http.ListenAndServe(":8080", nil))

}
