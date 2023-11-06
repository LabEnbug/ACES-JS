package model

import "time"

type VideoCommentChild struct {
	// Id, UserId, Content, QuoteRootCommentId, QuoteUserId, QuoteChildCommentId, CommentTime
	Id                  uint      `default:"0" json:"id"`
	UserId              uint      `default:"0" json:"-"`
	Content             string    `default:"" json:"content"`
	QuoteRootCommentId  uint      `default:"0" json:"quote_comment_id"`
	QuoteUserId         uint      `default:"0" json:"-"`
	QuoteChildCommentId uint      `default:"0" json:"quote_child_comment_id"`
	CommentTime         time.Time `default:"0000-00-00 00:00:00" json:"-"`

	User           interface{} `json:"user"`
	QuoteUser      interface{} `json:"quote_user"`
	CommentTimeStr string      `default:"" json:"comment_time"`
}
