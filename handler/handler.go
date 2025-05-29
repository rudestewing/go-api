package handler

import "go-api/container"

type Handler struct {
	AuthHandler		*AuthHandler
	UserHandler		*UserHandler
}

func NewHandler(container *container.Container) *Handler {
	return &Handler{
		AuthHandler: NewAuthHandler(container.AuthService),
		UserHandler: NewUserHandler(container.UserService),
	}
}
