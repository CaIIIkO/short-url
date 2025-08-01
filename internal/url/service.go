package url

import (
	"context"
	"errors"
	"math/rand"
	"short-url/internal/auth"
	"time"

	"github.com/google/uuid"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type RepositoryInterface interface {
	CreateLink(ctx context.Context, userID uuid.UUID, originalURL, code string) (*Link, error)
	GetLink(ctx context.Context, code string) (*Link, error)
	GetUserLinks(ctx context.Context, userID uuid.UUID) (*[]Link, error)
	LogClick(ctx context.Context, linkID uuid.UUID, ip, userAgent, referrer string) error
	GetClicks(ctx context.Context, linkID uuid.UUID) (*[]Click, error)
}

type Service struct {
	repo RepositoryInterface
}

func NewURLService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateLink(ctx context.Context, originalURL string) (*Link, error) {
	userID, _ := auth.UserIDFromContext(ctx)
	//Добавить валидацию оригинальной ссылки
	code := generateCode(6)

	return s.repo.CreateLink(ctx, userID, originalURL, code)
}

// Пересмотреть реализацию алгоритма!
func generateCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *Service) RedirectAndLog(ctx context.Context, code string, ip string, userAgent string, referrer string) (*Link, error) {
	link, err := s.repo.GetLink(ctx, code)
	if err != nil {
		return nil, errors.New("link not found")
	}

	//Логирование перехода по ссылке
	_ = s.repo.LogClick(ctx, link.ID, ip, userAgent, referrer)

	return link, nil
}

func (s *Service) Stats(ctx context.Context, code string) (*[]Click, *Link, error) {
	userID, _ := auth.UserIDFromContext(ctx)

	link, err := s.repo.GetLink(ctx, code)
	if err != nil || link == nil || link.UserID != userID {
		return nil, nil, errors.New("access denied or not found")
	}

	clicks, err := s.repo.GetClicks(ctx, link.ID)
	if err != nil {
		return nil, nil, err
	}
	return clicks, link, nil
}

func (s *Service) LinksList(ctx context.Context) (*[]Link, error) {
	userID, _ := auth.UserIDFromContext(ctx)
	links, err := s.repo.GetUserLinks(ctx, userID)
	if err != nil {
		return nil, err
	}
	return links, nil
}
