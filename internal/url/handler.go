package url

type ServiceInterface interface {
}

type Handler struct {
	service ServiceInterface
}

func NewURLHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}
