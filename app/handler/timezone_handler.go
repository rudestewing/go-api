package handler

import (
	"go-api/app/shared/response"
	"go-api/app/shared/timezone"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TimezoneHandler struct{}

func NewTimezoneHandler() *TimezoneHandler {
	return &TimezoneHandler{}
}

// GetTimezoneInfo returns current timezone information
func (h *TimezoneHandler) GetTimezoneInfo(c *fiber.Ctx) error {
	info := timezone.TimezoneInfo()

	// Add some additional helpful information
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
}

// GetCurrentTime returns current time in various formats
func (h *TimezoneHandler) GetCurrentTime(c *fiber.Ctx) error {
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
}

// ConvertTime converts time between timezones
func (h *TimezoneHandler) ConvertTime(c *fiber.Ctx) error {
	type ConvertRequest struct {
		DateTime     string `json:"datetime" validate:"required"`
		FromTimezone string `json:"from_timezone" validate:"required"`
		ToTimezone   string `json:"to_timezone" validate:"required"`
		Format       string `json:"format"` // optional, defaults to RFC3339
	}

	var req ConvertRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	// Set default format if not provided
	if req.Format == "" {
		req.Format = time.RFC3339
	}

	// Validate timezones
	if err := timezone.ValidateTimezone(req.FromTimezone); err != nil {
		return response.BadRequest(c, "Invalid from_timezone: "+err.Error())
	}

	if err := timezone.ValidateTimezone(req.ToTimezone); err != nil {
		return response.BadRequest(c, "Invalid to_timezone: "+err.Error())
	}

	// Load source timezone
	fromLoc, err := time.LoadLocation(req.FromTimezone)
	if err != nil {
		return response.BadRequest(c, "Failed to load from_timezone: "+err.Error())
	}

	// Load target timezone
	toLoc, err := time.LoadLocation(req.ToTimezone)
	if err != nil {
		return response.BadRequest(c, "Failed to load to_timezone: "+err.Error())
	}

	// Parse the input time in the source timezone
	inputTime, err := time.ParseInLocation(req.Format, req.DateTime, fromLoc)
	if err != nil {
		return response.BadRequest(c, "Failed to parse datetime with the specified format: "+err.Error())
	}

	// Convert to target timezone
	convertedTime := inputTime.In(toLoc)

	data := map[string]interface{}{
		"input": map[string]interface{}{
			"datetime":  req.DateTime,
			"timezone":  req.FromTimezone,
			"timestamp": inputTime.Unix(),
		},
		"converted": map[string]interface{}{
			"datetime":  convertedTime.Format(req.Format),
			"timezone":  req.ToTimezone,
			"timestamp": convertedTime.Unix(),
		},
		"utc": map[string]interface{}{
			"datetime":  inputTime.UTC().Format(req.Format),
			"timezone":  "UTC",
			"timestamp": inputTime.Unix(),
		},
	}

	return response.Success(c, data, "Time converted successfully")
}

// GetAvailableTimezones returns list of available timezones
func (h *TimezoneHandler) GetAvailableTimezones(c *fiber.Ctx) error {
	timezones := timezone.GetAvailableTimezones()

	// Add current time for each timezone for demonstration
	timezoneData := make([]map[string]interface{}, len(timezones))
	baseTime := time.Now().UTC()

	for i, tz := range timezones {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			continue
		}

		localTime := baseTime.In(loc)
		_, offset := localTime.Zone()

		timezoneData[i] = map[string]interface{}{
			"timezone":       tz,
			"current_time":   localTime.Format(time.RFC3339),
			"offset_seconds": offset,
			"offset_hours":   float64(offset) / 3600,
			"is_dst":         localTime.IsDST(),
		}
	}

	data := map[string]interface{}{
		"total_timezones": len(timezones),
		"timezones":       timezoneData,
	}

	return response.Success(c, data, "Available timezones retrieved successfully")
}
