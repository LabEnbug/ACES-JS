package model

import (
	"database/sql"
	"time"
)

type VideoAction struct {
	// Id, VideoId, UserId, ActionTime, UnActionTime
	Id           uint         `default:"0"`
	VideoId      uint         `default:"0"`
	UserId       uint         `default:"0"`
	ActionTime   time.Time    `default:"0000-00-00 00:00:00"`
	UnActionTime sql.NullTime `default:"0000-00-00 00:00:00"`
}
