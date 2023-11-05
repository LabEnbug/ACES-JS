package mysql

import (
	"backend/config"
	"backend/model"
	"backend/tool"
	"database/sql"
	"log"
	"runtime"
)

func GetVideoBulletComment(bulletCommentId uint, currentUserId uint) model.VideoBulletComment {
	var bulletComment model.VideoBulletComment
	err := DB.QueryRow("SELECT id, user_id, content, comment_at, comment_time FROM video_bullet_comment WHERE id=? AND delete_time IS NULL LIMIT 1", bulletCommentId).
		Scan(&bulletComment.Id, &bulletComment.UserId, &bulletComment.Content, &bulletComment.CommentAt, &bulletComment.CommentTime)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return bulletComment
	}
	bulletComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(bulletComment.CommentTime)
	user, _, _ := GetUserInfoById(bulletComment.UserId, currentUserId)
	bulletComment.User = user
	return bulletComment
}

func GetVideoBulletCommentList(videoId uint, limit int, start int, currentUserId uint) ([]model.VideoBulletComment, int) {
	var bulletCommentList []model.VideoBulletComment
	var rows *sql.Rows
	var err error
	var bulletCommentCount int

	rows, err = DB.Query("SELECT id, user_id, content, comment_at, comment_time FROM video_bullet_comment WHERE video_id=? AND delete_time IS NULL ORDER BY id DESC LIMIT ?, ?", videoId, start, limit)

	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err, "videoId:", videoId, "start:", start, "limit:", limit)
		}
		return bulletCommentList, bulletCommentCount
	}

	defer rows.Close()

	for rows.Next() {
		var bulletComment model.VideoBulletComment
		err := rows.Scan(&bulletComment.Id, &bulletComment.UserId, &bulletComment.Content, &bulletComment.CommentAt, &bulletComment.CommentTime)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err, "videoId:", videoId, "start:", start, "limit:", limit)
			}
			return bulletCommentList, bulletCommentCount
		}
		bulletComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(bulletComment.CommentTime)
		user, _, _ := GetUserInfoById(bulletComment.UserId, currentUserId)
		bulletComment.User = user
		bulletCommentList = append(bulletCommentList, bulletComment)
	}

	// get bullet comment count
	err = DB.QueryRow("SELECT COUNT(1) FROM video_bullet_comment WHERE video_id=? AND delete_time IS NULL", videoId).
		Scan(&bulletCommentCount)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err, "videoId:", videoId, "start:", start, "limit:", limit)
		}
		return bulletCommentList, bulletCommentCount
	}

	return bulletCommentList, bulletCommentCount
}

func MakeVideoBulletComment(videoId uint, userId uint, content string, commentAt float64) (int64, bool) {
	var res sql.Result
	var err error
	var commentId int64
	res, err = DB.Exec("INSERT INTO video_bullet_comment (video_id, user_id, content, comment_at, comment_time) VALUES (?, ?, ?, ?, NOW())", videoId, userId, content, commentAt)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err, "videoId:", videoId, "userId:", userId, "content:", content, "commentAt:", commentAt)
		}
		return commentId, false
	}

	// get bullet comment id
	commentId, err = res.LastInsertId()
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err, "videoId:", videoId, "userId:", userId, "content:", content, "commentAt:", commentAt)
		}
		return commentId, false
	}

	return commentId, true
}

func DeleteVideoBulletComment(bulletCommentId uint) bool {
	if _, err := DB.Exec("UPDATE video_bullet_comment SET delete_time=NOW() WHERE id=?", bulletCommentId); err != nil {
		return false
	}
	return true
}
