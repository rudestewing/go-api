package router

import (
	"go-api/app/handler"
	"go-api/app/middleware"
	"go-api/app/service"
	"go-api/app/shared/response"
	"go-api/app/shared/timezone"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *handler.Handler, authService *service.AuthService) {
	v1 := app.Group("/api/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, nil, "hello from v1", "API is running")
	})

	// Apply timeout middleware
	v1.Use(middleware.TimeoutMiddleware(30*time.Second, "User operation timed out"))

	// Auth routes (public) - no auth middleware needed
	auth := v1.Group("/auth")
	auth.Post("/register", handler.AuthHandler.Register)
	auth.Post("/login", handler.AuthHandler.Login)
	auth.Post("/logout", middleware.AuthMiddleware(authService), handler.AuthHandler.Logout)
	auth.Post("/logout-all", middleware.AuthMiddleware(authService), handler.AuthHandler.LogoutAll)

	// Protected routes with auth middleware
	protected := v1.Group("/user")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.Get("/profile", handler.UserHandler.GetProfile)

	// Timezone routes (public) - for demonstration and debugging
	timezoneGroup := v1.Group("/timezone")
	timezoneGroup.Get("/info", func(c *fiber.Ctx) error {
		info := timezone.TimezoneInfo()
		now := timezone.Now()
		utcNow := timezone.NowUTC()

		data := map[string]interface{}{
			"timezone_info": info,
			"formatted_times": map[string]interface{}{
				"local_datetime":    timezone.FormatDateTime(now),
				"local_date":        timezone.FormatDate(now),
				"local_time":        timezone.FormatTime(now),
				"local_rfc3339":     timezone.FormatRFC3339(now),
				"utc_datetime":      timezone.FormatDateTime(utcNow),
				"utc_rfc3339":       timezone.FormatRFC3339(utcNow),
			},
			"day_boundaries": map[string]interface{}{
				"today_start": timezone.FormatDateTime(timezone.TodayStart()),
				"today_end":   timezone.FormatDateTime(timezone.TodayEnd()),
			},
			"available_timezones": timezone.GetAvailableTimezones(),
		}

		return response.Success(c, data, "Timezone information retrieved successfully")
	})

	timezoneGroup.Get("/current", func(c *fiber.Ctx) error {
		now := timezone.Now()
		utcNow := timezone.NowUTC()

		data := map[string]interface{}{
			"current_timezone": timezone.GetTimezone(),
			"local_time": map[string]interface{}{
				"timestamp":  now.Unix(),
				"iso":        timezone.FormatRFC3339(now),
				"datetime":   timezone.FormatDateTime(now),
				"date":       timezone.FormatDate(now),
				"time":       timezone.FormatTime(now),
			},
			"utc_time": map[string]interface{}{
				"timestamp":  utcNow.Unix(),
				"iso":        timezone.FormatRFC3339(utcNow),
				"datetime":   timezone.FormatDateTime(utcNow),
				"date":       timezone.FormatDate(utcNow),
				"time":       timezone.FormatTime(utcNow),
			},
		}

		return response.Success(c, data, "Current time retrieved successfully")
	})

	v1.Use(middleware.TimeoutMiddleware(30 * time.Second))
	v1.Get("asynchronous", handler.AsynchronousHandler.RunAsyncTask)
}
