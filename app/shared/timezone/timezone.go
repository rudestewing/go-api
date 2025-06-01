package timezone

import (
	"fmt"
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

// Format formats a time in the configured timezone using the given layout
func Format(t time.Time, layout string) string {
	return ToLocal(t).Format(layout)
}

// FormatNow formats the current time in the configured timezone using the given layout
func FormatNow(layout string) string {
	return Now().Format(layout)
}

// ParseInTimezone parses a time string in the configured timezone
func ParseInTimezone(layout, value string) (time.Time, error) {
	loc := GetLocation()
	return time.ParseInLocation(layout, value, loc)
}

// AddExpiry adds duration to current time and returns it in the configured timezone
func AddExpiry(duration time.Duration) time.Time {
	return Now().Add(duration)
}

// IsExpired checks if a time is expired (before current time)
func IsExpired(t time.Time) bool {
	return t.Before(Now())
}

// TimezoneInfo returns information about the current timezone
func TimezoneInfo() map[string]interface{} {
	now := Now()
	_, offset := now.Zone()

	return map[string]interface{}{
		"timezone":    GetTimezone(),
		"current_time": now.Format(time.RFC3339),
		"utc_time":    NowUTC().Format(time.RFC3339),
		"offset_seconds": offset,
		"offset_hours":   float64(offset) / 3600,
		"is_dst":        now.IsDST(),
	}
}

// Common time formats with timezone
const (
	DateTimeFormat     = "2006-01-02 15:04:05"
	DateFormat         = "2006-01-02"
	TimeFormat         = "15:04:05"
	RFC3339Format      = time.RFC3339
	CustomFormat       = "02/01/2006 15:04:05"
	LogFormat          = "2006-01-02T15:04:05.000Z07:00"
)

// Helper functions for common formatting
func FormatDateTime(t time.Time) string {
	return Format(t, DateTimeFormat)
}

func FormatDate(t time.Time) string {
	return Format(t, DateFormat)
}

func FormatTime(t time.Time) string {
	return Format(t, TimeFormat)
}

func FormatRFC3339(t time.Time) string {
	return Format(t, RFC3339Format)
}

func FormatLog(t time.Time) string {
	return Format(t, LogFormat)
}

// StartOfDay returns the start of day (00:00:00) in the configured timezone
func StartOfDay(t time.Time) time.Time {
	local := ToLocal(t)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, GetLocation())
}

// EndOfDay returns the end of day (23:59:59.999999999) in the configured timezone
func EndOfDay(t time.Time) time.Time {
	local := ToLocal(t)
	return time.Date(local.Year(), local.Month(), local.Day(), 23, 59, 59, 999999999, GetLocation())
}

// TodayStart returns the start of today in the configured timezone
func TodayStart() time.Time {
	return StartOfDay(Now())
}

// TodayEnd returns the end of today in the configured timezone
func TodayEnd() time.Time {
	return EndOfDay(Now())
}

// GetAvailableTimezones returns a list of common timezones
func GetAvailableTimezones() []string {
	return []string{
		"UTC",
		"Asia/Jakarta",      // WIB (UTC+7)
		"Asia/Makassar",     // WITA (UTC+8)
		"Asia/Jayapura",     // WIT (UTC+9)
		"Asia/Singapore",    // SGT (UTC+8)
		"Asia/Kuala_Lumpur", // MYT (UTC+8)
		"Asia/Bangkok",      // ICT (UTC+7)
		"Asia/Manila",       // PHT (UTC+8)
		"Asia/Tokyo",        // JST (UTC+9)
		"Asia/Seoul",        // KST (UTC+9)
		"Asia/Shanghai",     // CST (UTC+8)
		"Asia/Hong_Kong",    // HKT (UTC+8)
		"Asia/Kolkata",      // IST (UTC+5:30)
		"Asia/Dubai",        // GST (UTC+4)
		"Europe/London",     // GMT/BST (UTC+0/+1)
		"Europe/Paris",      // CET/CEST (UTC+1/+2)
		"Europe/Berlin",     // CET/CEST (UTC+1/+2)
		"Europe/Rome",       // CET/CEST (UTC+1/+2)
		"Europe/Moscow",     // MSK (UTC+3)
		"America/New_York",  // EST/EDT (UTC-5/-4)
		"America/Chicago",   // CST/CDT (UTC-6/-5)
		"America/Denver",    // MST/MDT (UTC-7/-6)
		"America/Los_Angeles", // PST/PDT (UTC-8/-7)
		"America/Sao_Paulo", // BRT (UTC-3)
		"Australia/Sydney",  // AEST/AEDT (UTC+10/+11)
		"Australia/Melbourne", // AEST/AEDT (UTC+10/+11)
		"Pacific/Auckland",  // NZST/NZDT (UTC+12/+13)
	}
}

// ValidateTimezone checks if a timezone string is valid
func ValidateTimezone(tz string) error {
	_, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone '%s': %w", tz, err)
	}
	return nil
}
