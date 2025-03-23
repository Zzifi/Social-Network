package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var (
	dbConnStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
)

type Users struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Sessions struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserProfile struct {
	UserID      int       `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Birthday    string    `json:"birthday"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Credentials struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Storage struct {
	logger *slog.Logger
	db     *sql.DB
}

func New(logger *slog.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Error("Ошибка подключения к БД:", err)
		return nil, err
	}

	return &Storage{logger: logger, db: db}, nil
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Storage) CreateUser(creds Credentials) (int, error) {
	query := `INSERT INTO users (username, email, password_hash, role, created_at) VALUES ($1, $2, $3, 'user', NOW()) RETURNING id`
	var userID int
	err := s.db.QueryRow(query, creds.UserName, creds.Email, creds.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *Storage) GetUserByEmail(email string) (*Users, error) {
	var user Users
	err := s.db.QueryRow(`SELECT id, email, password_hash FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) CreateSession(userID int, token string, expirationTime time.Time) error {
	_, err := s.db.Exec(`INSERT INTO sessions (user_id, token, created_at, expires_at) VALUES ($1, $2, NOW(), $3)`,
		userID, token, expirationTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) IsSessionValid(token string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1 AND expires_at > NOW())`, token).
		Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Storage) GetUserProfile(userID int) (*UserProfile, error) {
	var profile UserProfile
	err := s.db.QueryRow(`
		SELECT user_id, phone_number, first_name, last_name, birthday, bio, avatar_url, updated_at 
		FROM user_profile WHERE user_id = $1`, userID).
		Scan(&profile.UserID, &profile.PhoneNumber, &profile.FirstName, &profile.LastName,
			&profile.Birthday, &profile.Bio, &profile.AvatarURL, &profile.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *Storage) UpdateUserProfile(profile UserProfile) error {
	_, err := s.db.Exec(`
	INSERT INTO user_profile (user_id, phone_number, first_name, last_name, birthday, bio, avatar_url)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (user_id) 
	DO UPDATE SET 
		phone_number = EXCLUDED.phone_number,
		first_name = EXCLUDED.first_name,
		last_name = EXCLUDED.last_name,
		birthday = EXCLUDED.birthday,
		bio = EXCLUDED.bio,
		avatar_url = EXCLUDED.avatar_url,
		updated_at = NOW()
		`, profile.UserID, profile.PhoneNumber, profile.FirstName, profile.LastName, profile.Birthday, profile.Bio, profile.AvatarURL)

	return err
}
