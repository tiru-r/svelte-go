package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo      *Repository
	jwtSecret []byte
	eventBus  EventPublisher
}

type EventPublisher interface {
	Publish(event string, data any) error
}

func NewService(repo *Repository, jwtSecret string, eventBus EventPublisher) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		eventBus:  eventBus,
	}
}

func (s *Service) Register(req RegisterRequest) (*LoginResponse, error) {
	log.Printf("🔐 Registering new user: %s", req.Email)

	user := &User{
		Email:    req.Email,
		Username: req.Username,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.eventBus.Publish("user.registered", map[string]any{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Username,
	})

	user.Password = ""
	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	log.Printf("🔐 User login attempt: %s", req.Email)

	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	if !user.VerifyPassword(req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	s.eventBus.Publish("user.logged_in", map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
	})

	user.Password = ""
	return &LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) VerifyToken(tokenString string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		userID, _ := (*claims)["user_id"].(string)
		username, _ := (*claims)["username"].(string)
		email, _ := (*claims)["email"].(string)
		exp, _ := (*claims)["exp"].(float64)
		iat, _ := (*claims)["iat"].(float64)

		return &AuthClaims{
			UserID:    userID,
			Username:  username,
			Email:     email,
			ExpiresAt: int64(exp),
			IssuedAt:  int64(iat),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *Service) GenerateJWT(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) Logout(userID string) error {
	log.Printf("🔐 User logout: %s", userID)

	err := s.repo.InvalidateUserSessions(userID)
	if err != nil {
		log.Printf("Warning: Failed to invalidate sessions for user %s: %v", userID, err)
	}

	s.eventBus.Publish("user.logged_out", map[string]any{
		"user_id": userID,
	})

	return nil
}

func (s *Service) GetUserByID(id string) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *Service) RefreshToken(tokenString string) (string, error) {
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	if time.Now().Unix() > claims.ExpiresAt-3600 {
		user, err := s.repo.GetUserByID(claims.UserID)
		if err != nil {
			return "", err
		}

		return s.GenerateJWT(user)
	}

	return tokenString, nil
}

func GenerateSecretKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate secret key")
	}
	return base64.StdEncoding.EncodeToString(key)
}