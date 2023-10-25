package model

import "time"

type Video struct {
	// Id, UserId, VideoUid, Type, Title, Content, UploadTime, LikeCount, FavoriteCount, CommentCount
	Id            int       `default:"0" json:"id"`
	UserId        int       `default:"0" json:"user_id"`
	VideoUid      string    `default:"" json:"video_uid"`
	Type          int       `default:"0" json:"type"`
	Title         string    `default:"" json:"title"`
	Content       string    `default:"" json:"content"`
	UploadTime    time.Time `default:"0000-00-00 00:00:00" json:"upload_time"`
	LikeCount     int       `default:"0" json:"like_count"`
	FavoriteCount int       `default:"0" json:"favorite_count"`
	CommentCount  int       `default:"0" json:"comment_count"`
}
