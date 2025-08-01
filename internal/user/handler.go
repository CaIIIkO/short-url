package user

type ServiceInterface interface {
}

type Handler struct {
	service ServiceInterface
}

func NewAdHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}
