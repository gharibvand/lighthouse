package service

import (
	"errors"
	"lighthouse/internal/domain"
	"lighthouse/internal/repository/postgres"
	"lighthouse/pkg/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *postgres.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *postgres.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) RegisterWithEmail(email, password string) (*domain.User, string, error) {
	// Check if user exists
	existing, _ := s.userRepo.GetByEmail(email)
	if existing != nil {
		return nil, "", errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &domain.User{
		Email:    &email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) RegisterWithPhone(phone, password string) (*domain.User, string, error) {
	// Check if user exists
	existing, _ := s.userRepo.GetByPhone(phone)
	if existing != nil {
		return nil, "", errors.New("user with this phone already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &domain.User{
		Phone:    &phone,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login accepts either email or phone as identifier
func (s *AuthService) Login(identifier, password string) (*domain.User, string, error) {
	var user *domain.User
	var err error

	// Try email first, then phone
	user, err = s.userRepo.GetByEmail(identifier)
	if err != nil || user == nil {
		// Try phone
		user, err = s.userRepo.GetByPhone(identifier)
		if err != nil || user == nil {
			return nil, "", errors.New("invalid credentials")
		}
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}
