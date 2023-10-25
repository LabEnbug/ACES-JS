package model

import (
	"database/sql"
	"time"
)

type VideoComment struct {
	// Id, VideoId, UserId, Content, QuoteId, CommentTime, DeleteTime
	Id          int64        `default:"0"`
	VideoId     int          `default:"0"`
	UserId      int          `default:"0"`
	Content     string       `default:""`
	QuoteId     int64        `default:"0"`
	CommentTime time.Time    `default:"0000-00-00 00:00:00"`
	DeleteTime  sql.NullTime `default:"0000-00-00 00:00:00"`
}
