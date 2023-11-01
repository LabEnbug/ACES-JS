package model

import (
	"database/sql"
	"time"
)

type User struct {
	// Id, Username, Password, Nickname, FollowCount, BeFollowed, BeFollowedCount, BeLikedCount, BeFavoriteCount,
	// BeCommentedCount, BeForwardedCount, BeWatchedCount, RegTime, LastLoginTime
	Id               uint         `default:"0" json:"user_id"`
	Username         string       `default:"" json:"username"`
	Password         string       `default:"" json:"-"`
	Nickname         string       `default:"" json:"nickname"`
	FollowCount      uint         `default:"0" json:"follow_count"`
	BeFollowed       bool         `default:"false" json:"be_followed"`
	BeFollowedCount  uint         `default:"0" json:"be_followed_count"`
	BeLikedCount     uint         `default:"0" json:"be_liked_count"`
	BeFavoriteCount  uint         `default:"0" json:"be_favorite_count"`
	BeCommentedCount uint         `default:"0" json:"be_commented_count"`
	BeForwardedCount uint         `default:"0" json:"be_forwarded_count"`
	BeWatchedCount   uint         `default:"0" json:"be_watched_count"`
	RegTime          time.Time    `default:"0000-00-00 00:00:00" json:"-"`
	RegTimeStr       string       `default:"" json:"reg_time"`
	LastLoginTime    sql.NullTime `default:"0000-00-00 00:00:00" json:"-"`
}
