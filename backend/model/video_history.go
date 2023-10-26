package model

import (
	"time"
)

type VideoHistory struct {
	// Id, VideoId, UserId, RecordTime
	Id         uint      `default:"0"`
	VideoId    uint      `default:"0"`
	UserId     uint      `default:"0"`
	RecordTime time.Time `default:"0000-00-00 00:00:00"`
}
