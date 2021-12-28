package time

import (
	"time"
)

const (
	// ISO8601 ISO8601
	ISO8601 = "2006-01-02T15:04:05.999Z"
)

// Time time now utc
func Time() time.Time {
	return time.Now().UTC()
}

// NowUnix current UTC timestamp.
func NowUnix() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

// Now current time in ISO8601 format.
func Now() string {
	return time.Now().UTC().Format(ISO8601)
}

// Format timestamp to time string.
func Format(t int64) string {
	return time.Unix(0, t*1e6).UTC().Format(ISO8601)
}

// Unix time string to timestamp.
func Unix(ts string) (int64, error) {
	t, err := time.Parse(ISO8601, ts)
	if err != nil {
		return 0, err
	}

	return t.UTC().UnixNano() / 1e6, nil
}
