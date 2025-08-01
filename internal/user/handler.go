package user

type ServiceInterface interface {
}

type Handler struct {
	service ServiceInterface
}

func NewUserHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}
