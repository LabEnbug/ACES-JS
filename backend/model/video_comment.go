package model

import "time"

type VideoComment struct {
	// Id, VideoUid, UserId, Content, QuoteCommentId, CommentTime
	Id             uint        `default:"0" json:"id"`
	VideoUid       string      `default:"" json:"video_uid"`
	UserId         uint        `default:"0" json:"-"`
	User           interface{} `json:"user"`
	Content        string      `default:"" json:"content"`
	QuoteCommentId uint        `default:"0" json:"quote_comment_id"`
	CommentTime    time.Time   `default:"0000-00-00 00:00:00" json:"-"`
	CommentTimeStr string      `default:"" json:"comment_time"`
}
