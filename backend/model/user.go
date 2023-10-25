package model

import (
	"database/sql"
	"time"
)

type User struct {
	// Id, Username, Password, Nickname, LikeCount, FavoriteCount, CommentCount, RegTime, LastLoginTime
	Id            int          `default:"0" json:"id"`
	Username      string       `default:"" json:"username"`
	Password      string       `default:"" json:"password"`
	Nickname      string       `default:"" json:"nickname"`
	LikeCount     int          `default:"0" json:"like_count"`
	FavoriteCount int          `default:"0" json:"favorite_count"`
	CommentCount  int          `default:"0" json:"comment_count"`
	RegTime       time.Time    `default:"0000-00-00 00:00:00" json:"reg_time"`
	LastLoginTime sql.NullTime `default:"0000-00-00 00:00:00" json:"last_login_time"`
}
