package model

import (
	"time"
)

type Video struct {
	// Id, UserId, VideoUid, Type, Title, Content, UploadTime, LikeCount, FavoriteCount, CommentCount, CoverUrl, PlayUrl
	Id             uint        `default:"0" json:"id"`
	UserId         uint        `default:"0" json:"user_id"`
	User           interface{} `json:"user"`
	VideoUid       string      `default:"" json:"video_uid"`
	Type           int8        `default:"0" json:"type"`
	Title          string      `default:"" json:"title"`
	Content        string      `default:"" json:"content"`
	UploadTime     time.Time   `default:"0000-00-00 00:00:00" json:"upload_time"`
	LikeCount      uint        `default:"0" json:"like_count"`
	FavoriteCount  uint        `default:"0" json:"favorite_count"`
	CommentCount   uint        `default:"0" json:"comment_count"`
	CoverUrl       string      `default:"" json:"cover_url"`
	PlayUrl        string      `default:"" json:"play_url"`
	IsUserLiked    bool        `default:"false" json:"is_user_liked"`
	IsUserFavorite bool        `default:"false" json:"is_user_favorite"`
	IsUserUploaded bool        `default:"false" json:"is_user_uploaded"`
	IsUserHistory  bool        `default:"false" json:"is_user_history"`
	IsUserLastPlay bool        `default:"false" json:"is_user_last_play"`
}
