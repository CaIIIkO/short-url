package user

import (
	"context"
	"errors"
	"net/mail"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type RepositoryInterface interface {
	Create(ctx context.Context, u *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewUserService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// Register - регистрация пользователя
func (s *Service) Register(ctx context.Context, input *RegisterRequest) (*User, error) {
	if err := s.validateRegisterInput(ctx, input); err != nil {
		return nil, err
	}

	//Хэширование пароля
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        input.Email,
		PasswordHash: string(hashed),
	}

	user, err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var (
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]{6,30}$`)
)

// validateRegisterInput проверяет корректность входных данных при регистрации
func (s *Service) validateRegisterInput(ctx context.Context, input *RegisterRequest) error {
	// Валидация email
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		return errors.New("invalid email")
	}

	//Проверка уникальности email
	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("user email already exists")
	}

	// Валидация пароля
	if len(input.Password) < 6 || len(input.Password) > 30 {
		return errors.New("invalid password: must be at least 6 - 30 characters")
	}

	if !passwordRegex.MatchString(input.Password) {
		return errors.New(`invalid password: ust contain only a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?`)
	}

	return nil
}
