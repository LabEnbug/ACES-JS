package model

import (
	"time"
)

type VideoHistory struct {
	// Id, VideoId, UserId, RecordTime
	Id         int64     `default:"0"`
	VideoId    int       `default:"0"`
	UserId     int       `default:"0"`
	RecordTime time.Time `default:"0000-00-00 00:00:00"`
}
