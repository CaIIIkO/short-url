package url

type RepositoryInterface interface {
}

type Service struct {
	repo RepositoryInterface
}

func NewURLService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}
