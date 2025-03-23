package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/storage/postgre"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecretKey = "your_secret_key"
)

type Handlers struct {
	logger  *slog.Logger
	storage *postgre.Storage
}

func NewHandlers(logger *slog.Logger, storage *postgre.Storage) *Handlers {
	return &Handlers{logger: logger, storage: storage}
}

func (h *Handlers) RegisterRoutes(router chi.Router) {
	router.Post("/register", h.RegisterHandler())
	router.Post("/login", h.LoginHandler())
	router.Get("/session/validate", h.SessionValidateHandler())
	router.Get("/user/{id}", h.GetUserProfileHandler())
	router.Put("/user/{id}", h.UpdateUserProfileHandler())
}

func (h *Handlers) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds postgre.Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка хэширования пароля", http.StatusInternalServerError)
			return
		}

		creds.Password = string(hashedPassword)
		userID, err := h.storage.CreateUser(creds)
		if err != nil {
			http.Error(w, "Ошибка сохранения пользователя", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": userID})
	}
}

func (h *Handlers) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds postgre.Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
			return
		}

		user, err := h.storage.GetUserByEmail(creds.Email)
		if err != nil {
			http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
			http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
			return
		}

		token, err := h.generateToken(user.ID)
		if err != nil {
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		if err := h.storage.CreateSession(user.ID, token, time.Now().Add(24*time.Hour)); err != nil {
			http.Error(w, "Ошибка сохранения сессии", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(postgre.Token{Token: token})
	}
}

func (h *Handlers) SessionValidateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		if tokenString == "" {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		exists, err := h.storage.IsSessionValid(tokenString)
		if err != nil {
			http.Error(w, "Ошибка поиска сессии", http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, "Сессия не найдена, или она истекла: ", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) GetUserProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}

		profile, err := h.storage.GetUserProfile(userID)
		if err != nil {
			http.Error(w, "Ошибка получения профиля", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(profile)
	}
}

func (h *Handlers) UpdateUserProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Некорректный ID", http.StatusBadRequest)
			return
		}

		var profile postgre.UserProfile
		if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
			http.Error(w, "Ошибка парсинга JSON", http.StatusBadRequest)
			return
		}

		profile.UserID = userID
		if err := h.storage.UpdateUserProfile(profile); err != nil {
			http.Error(w, "Ошибка обновления данных", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}
