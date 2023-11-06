package model

import "time"

type VideoBulletComment struct {
	// Id, UserId, Content, CommentTime
	Id          uint      `default:"0" json:"id"`
	UserId      uint      `default:"0" json:"-"`
	Content     string    `default:"" json:"content"`
	CommentAt   float64   `default:"0" json:"comment_at"`
	CommentTime time.Time `default:"0000-00-00 00:00:00" json:"-"`

	User           interface{} `json:"user"`
	CommentTimeStr string      `default:"" json:"comment_time"`
}
