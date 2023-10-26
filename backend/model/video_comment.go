package model

import (
	"database/sql"
	"time"
)

type VideoComment struct {
	// Id, VideoId, UserId, Content, QuoteId, CommentTime, DeleteTime
	Id          uint         `default:"0"`
	VideoId     uint         `default:"0"`
	UserId      uint         `default:"0"`
	Content     string       `default:""`
	QuoteId     uint         `default:"0"`
	CommentTime time.Time    `default:"0000-00-00 00:00:00"`
	DeleteTime  sql.NullTime `default:"0000-00-00 00:00:00"`
}
