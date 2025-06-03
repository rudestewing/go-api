package timezone

import (
	"go-api/config"
	"time"
)

// Now returns the current time in the configured timezone
func Now() time.Time {
	return config.Now()
}

// NowUTC returns the current time in UTC
func NowUTC() time.Time {
	return config.NowUTC()
}

// ToLocal converts a UTC time to the configured timezone
func ToLocal(t time.Time) time.Time {
	return config.ToLocalTime(t)
}

// ToUTC converts a time from the configured timezone to UTC
func ToUTC(t time.Time) time.Time {
	return config.ToUTC(t)
}

// GetTimezone returns the configured timezone string
func GetTimezone() string {
	return config.GetTimezone()
}

// GetLocation returns the configured timezone location
func GetLocation() *time.Location {
	return config.GetTimezoneLocation()
}
