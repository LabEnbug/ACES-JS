package model

import (
	"time"
)

type VideoForward struct {
	// Id, VideoId, UserId, ForwardTime
	Id          uint      `default:"0"`
	VideoId     uint      `default:"0"`
	UserId      uint      `default:"0"`
	ForwardTime time.Time `default:"0000-00-00 00:00:00"`
}
