package model

import (
	"time"
)

type VideoWatch struct {
	// Id, VideoId, UserId, WatchedTime
	Id          uint      `default:"0"`
	VideoId     uint      `default:"0"`
	UserId      uint      `default:"0"`
	WatchedTime time.Time `default:"0000-00-00 00:00:00"`
}
