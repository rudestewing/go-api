package router

import (
	"go-api/app"
	auth "go-api/domain/auth/handler"
	healthcheck "go-api/domain/healthcheck/handler"
)

type Handler struct {
	health *healthcheck.HealthHandler
	auth *auth.AuthHandler
}

func NewHandler(app *app.Provider) *Handler {
	return &Handler{
		health: healthcheck.NewHealthHandler(app),
		auth: auth.NewAuthHandler(app),
	}
}