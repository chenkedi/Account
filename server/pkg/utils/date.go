package utils

import "time"

func NowUTC() time.Time {
	return time.Now().UTC()
}

func ToUTC(t time.Time) time.Time {
	return t.UTC()
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
