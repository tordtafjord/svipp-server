package util

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
	"time"
)

const MicrosecondsPerHour = 3600_000_000
const MicrosecondsPerMinute = 60_000_000

// Form input type time to pg type time. Returns null type if not valid
func TimeInputToPgTime(timeStr string) pgtype.Time {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return pgtype.Time{Valid: false}
	}

	microsecondsSinceMidnight := (int64(t.Hour()) * MicrosecondsPerHour) +
		(int64(t.Minute()) * MicrosecondsPerMinute)

	return pgtype.Time{
		Microseconds: microsecondsSinceMidnight,
		Valid:        true,
	}
}

// Helper function to convert string to *string, returning nil if empty/whitespace
func StringToPtr(s string) *string {
	// Trim whitespace
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}
