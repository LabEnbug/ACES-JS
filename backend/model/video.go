package model

import (
	"time"
)

type Video struct {
	// Id, UserId, VideoUid, Type, Content, Keyword, UploadTime, BeLikedCount, BeFavoriteCount, BeCommentedCount,
	// BeForwardedCount, CoverUrl, PlayUrl, IsUserLiked, IsUserFavorite, IsUserUploaded, IsUserWatched, IsUserLastPlay
	Id               uint        `default:"0" json:"-"`
	UserId           uint        `default:"0" json:"-"`
	User             interface{} `json:"user"`
	VideoUid         string      `default:"" json:"video_uid"`
	Type             int8        `default:"0" json:"type"`
	Content          string      `default:"" json:"content"`
	Keyword          string      `default:"" json:"keyword"`
	UploadTime       time.Time   `default:"0000-00-00 00:00:00" json:"-"`
	UploadTimeStr    string      `default:"" json:"upload_time"`
	BeLikedCount     uint        `default:"0" json:"be_liked_count"`
	BeFavoriteCount  uint        `default:"0" json:"be_favorite_count"`
	BeCommentedCount uint        `default:"0" json:"be_commented_count"`
	BeForwardedCount uint        `default:"0" json:"be_forwarded_count"`
	BeWatchedCount   uint        `default:"0" json:"be_watched_count"`
	CoverUrl         string      `default:"" json:"cover_url"`
	PlayUrl          string      `default:"" json:"play_url"`
	IsUserLiked      bool        `default:"false" json:"is_user_liked"`
	IsUserFavorite   bool        `default:"false" json:"is_user_favorite"`
	IsUserUploaded   bool        `default:"false" json:"is_user_uploaded"`
	IsUserWatched    bool        `default:"false" json:"is_user_watched"`
	IsUserLastPlay   bool        `default:"false" json:"is_user_last_play"`
}
