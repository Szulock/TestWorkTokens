package middleware

import (
	"TokenTestWork/internal/loggerer"
	"TokenTestWork/internal/storage"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/smtp"

	"golang.org/x/crypto/bcrypt"
)

var (
	Logger = loggerer.NewLogger()
	db     = storage.GetDb(Logger)
) //почему я использую глобальные переменные - для простоты. Да, по best practices лучше так не делать, но
// тут нет большого и сложного функционала, поэтому так

func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Info("Сканирую id из поиска")
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		Logger.Info("id пользователя не заполнен")
		http.Error(w, "id пользователя не заполнен, чтобы заполнить, добавь в поисковую строку ?user_id=12 (число может быть произвольным)", http.StatusBadRequest)
		return
	}
	Logger.Info("Получаю ip")
	ip := r.RemoteAddr

	Logger.Info("Генерирую токен")
	accessToken, refreshToken, err := generateTokens(userID, ip, db) //в generateTokens также происходит сохранение в бд ip
	if err != nil {
		Logger.Printf("Не получилось создать токен: %v", err)
		http.Error(w, "Не получилось создать токен", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	Logger.Info("Создаю джсон")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		Logger.Printf("Проблема кодировании JSON: %v", err)
		http.Error(w, "Проблема при формировании ответа", http.StatusInternalServerError)
	}
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	Logger.Info("Расшифровываю джсон")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Logger.Printf("Не удалось получить данные: %v", err)
		http.Error(w, "Не удалось получить данные", http.StatusBadRequest)
		return
	}
	Logger.Info("Получаю ip")
	ip := r.RemoteAddr
	userID := ""
	Logger.Info("Проверяю токен в базе данных")

	rows, err := db.Query("SELECT user_id, token, ip_address FROM refresh_tokens")
	if err != nil {
		Logger.Printf("Ошибка при поиске токена в базе данных: %v", err)
		http.Error(w, "Ошибка при поиске токена в базе данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	Logger.Info("Токены найдены")
	for rows.Next() {
		var uid, hashedToken, oldIP string

		if err := rows.Scan(&uid, &hashedToken, &oldIP); err != nil {
			Logger.Printf("Ошибка при поиске токена: %v", err)
			http.Error(w, "Ошибка при поиске токена", http.StatusInternalServerError)
			return
		}

		rawRefreshToken := base64.StdEncoding.EncodeToString([]byte(uid + ":" + oldIP))
		err = bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(rawRefreshToken))
		if err == nil {
			userID = uid
			Logger.Info("Токен успешно проверен")

			Logger.Info("Проверяю ip")
			if oldIP != ip {
				Logger.Warn("ip изменился, отправляю предупреждение")
				sendEmailWarning(oldIP, ip)
			}
			continue
		}
	}

	if userID == "" {
		Logger.Info("Некорректный refresh токен")
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	Logger.Info("Генерирую новый access токен")
	accessToken, _, err := generateTokens(userID, ip, db)
	if err != nil {
		Logger.Printf("Не удалось сгенерировать новый access токен: %v", err)
		http.Error(w, "Не удалось сгенерировать новый access токен", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		Logger.Printf("Ошибка при кодировании джсона: %v", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}

func sendEmailWarning(oldIP, newIP string) {
	Logger.Info("Отправляю email")
	from := "MedodsmockMail@gmail.com"
	to := "userEmailMock@gmail.com"
	subject := "Ваш ip изменился"
	body := "Ваш ip адрес изменился с " + oldIP + " на " + newIP

	message := []byte("Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	smtp.SendMail("smtp.MockMedodsDomain.ru:111", smtp.PlainAuth("", "MedodsmockMail@gmail.com", "emailPass", "smtp.MockMedodsDomain.ru:111"), from, []string{to}, message)
	Logger.Info("email отправлен")
}
