package handler

import "go-api/container"

type Handler struct {
	AuthHandler *AuthHandler
}

func NewHandler(container *container.Container) *Handler {
	authHandler := NewAuthHandler(
		container.AuthService,
		container.EmailService,
	)

	return &Handler{
		AuthHandler: authHandler,
	}
}
