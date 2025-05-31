package handler

import "go-api/container"

type Handler struct {
	AuthHandler		*AuthHandler
	UserHandler		*UserHandler
	AsynchronousHandler *AsynchronousHandler
}

func NewHandler(container *container.Container) *Handler {
	return &Handler{
		AuthHandler: NewAuthHandler(container.AuthService),
		UserHandler: NewUserHandler(container.UserService),
		AsynchronousHandler: NewAsynchronousHandler(),
	}
}
