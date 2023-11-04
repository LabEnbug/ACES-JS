package model

import "time"

type VideoCommentChild struct {
	// Id, UserId, Content, QuoteRootCommentId, QuoteUserId, QuoteChildCommentId, CommentTime
	Id                  uint        `default:"0" json:"id"`
	VideoUid            string      `default:"" json:"-"`
	Content             string      `default:"" json:"content"`
	QuoteRootCommentId  uint        `default:"0" json:"quote_comment_id"`
	QuoteUserId         uint        `default:"0" json:"-"`
	QuoteUser           interface{} `json:"quote_user"`
	QuoteChildCommentId uint        `default:"0" json:"quote_child_comment_id"`
	CommentTime         time.Time   `default:"0000-00-00 00:00:00" json:"-"`
	UserId              uint        `default:"0" json:"-"`

	CommentTimeStr string      `default:"" json:"comment_time"`
	User           interface{} `json:"user"`
}
