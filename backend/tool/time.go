package tool

import "time"

func UnixTimeToRFC3339(t int64) string {
	return time.Unix(t, 0).UTC().Format("2006-01-02T15:04:05.000Z")
}

func DatabaseTimeToRFC3339(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
