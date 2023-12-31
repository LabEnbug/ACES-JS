package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type User struct {
	Id               uint            `default:"0" json:"user_id"`
	Username         string          `default:"" json:"username"`
	Password         string          `default:"" json:"-"`
	Nickname         string          `default:"" json:"nickname"`
	Avatar           sql.NullString  `default:"" json:"-"`
	DepositAmount    float64         `default:"" json:"-"`                 // add on 20231105
	Balance          float64         `default:"" json:"balance,omitempty"` // add on 20231105
	FollowCount      uint            `default:"0" json:"follow_count"`
	BeFollowedCount  uint            `default:"0" json:"be_followed_count"`
	BeLikedCount     uint            `default:"0" json:"be_liked_count"`
	BeFavoriteCount  uint            `default:"0" json:"be_favorite_count"`
	BeCommentedCount uint            `default:"0" json:"be_commented_count"`
	BeForwardedCount uint            `default:"0" json:"be_forwarded_count"`
	BeWatchedCount   uint            `default:"0" json:"be_watched_count"`
	RecommendMatrix  json.RawMessage `json:"-"`
	RegTime          time.Time       `default:"0000-00-00 00:00:00" json:"-"`
	LastLoginTime    sql.NullTime    `default:"0000-00-00 00:00:00" json:"-"`

	BeFollowed bool   `default:"false" json:"be_followed"`
	RegTimeStr string `default:"" json:"reg_time"`
	IsSelf     bool   `default:"false" json:"is_self"`
	AvatarUrl  string `default:"" json:"avatar_url,omitempty"`
}
