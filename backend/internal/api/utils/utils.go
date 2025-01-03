package utils

import (
	"database/sql"
	"fmt"
	"go-fitsync/backend/internal/database/sqlc"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: upgrade these to accept a variadic input for valid checking in instances like "01-seed-users.go"

// Used to parse client's timezone from the custom HTTP header "X-User-Timezone" and convert the client's timezone to UTC, maintaining the same calendar date in client's time. Client must still send its current time in the request body.
func FromClientTimezoneToUTC(clientTime time.Time, r *http.Request) (time.Time, error) {
	clientTZ := r.Header.Get("X-User-Timezone")
	if clientTZ == "" {
		clientTZ = "UTC" // fallback
	}

	location, err := time.LoadLocation(clientTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %w", err)
	}

	utcTime := time.Date(
		clientTime.Year(),
		clientTime.Month(),
		clientTime.Day(),
		0, 0, 0, 0,
		location,
	)

	return utcTime, nil
}

func GetIDFromPath(path string) (int32, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path")
	}

	id, err := strconv.ParseInt(parts[len(parts)-1], 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

// Used in APIs' SQLc params section for compact conversions.
func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// Used in APIs' SQLc params section for compact conversions.
func ToNullInt32(num interface{}) sql.NullInt32 {
	switch n := num.(type) {
	case int:
		return sql.NullInt32{
			Int32: int32(n),
			Valid: true,
		}
	case int64:
		return sql.NullInt32{
			Int32: int32(n),
			Valid: true,
		}
	case int32:
		return sql.NullInt32{
			Int32: n,
			Valid: true,
		}
	default:
		return sql.NullInt32{Valid: false} // invalid NullInt32 is safer than panicking
	}
}

// Used in APIs' SQLc params section for compact conversions.
func ToNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

// Used in APIs with optional fields that utilize pointers
func IntPtr(i int32) *int32         { return &i }
func StrPtr(s string) *string       { return &s }
func Float32Ptr(f float32) *float32 { return &f }

// For converting pointers to SQL null types
func NullIntFromIntPtr(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func NullStringFromStringPtr(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NullStringFromFloat32Ptr(f *float32) sql.NullString {
	if f == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: fmt.Sprintf("%.1f", *f), Valid: true}
}

func NullResistanceTypeEnumFromStringPtr(s *string) sqlc.NullResistanceTypeEnum {
	if s == nil {
		return sqlc.NullResistanceTypeEnum{
			Valid: false,
		}
	}
	return sqlc.NullResistanceTypeEnum{
		ResistanceTypeEnum: sqlc.ResistanceTypeEnum(*s),
		Valid:              true,
	}
}
