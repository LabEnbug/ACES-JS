package database

import (
	"backend/common"
	"backend/config"
	"backend/model"
	"backend/tool"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var (
	db  *sql.DB
	err error
)

func InitDBPool() {
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DefaultDatabase.User, config.DefaultDatabase.Pass, config.DefaultDatabase.Host, config.DefaultDatabase.Port,
		config.DefaultDatabase.Name))
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err:", err)
		}
		panic(err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err := db.Ping()
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err:", err)
		}
		panic(err)
	}
}

func CloseDBPool() {
	err := db.Close()
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return
	}
}

func GetUserIsFollowed(userId uint, followUserId uint) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM user_follow WHERE user_id=? AND follow_user_id=? AND unfollow_time IS NULL LIMIT 1", userId, followUserId).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func GetUserInfoById(userId uint, currentUserId uint) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	err := db.QueryRow("SELECT id, username, password, nickname, follow_count, be_followed_count, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, reg_time, last_login_time FROM user WHERE id=? LIMIT 1", userId).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.FollowCount, &user.BeFollowedCount, &user.BeLikedCount, &user.BeFavoriteCount, &user.BeCommentedCount, &user.BeForwardedCount, &user.BeWatchedCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		if errors.Is(err, sql.ErrNoRows) { // no such user
			ok = false
			errNo = 1
		} else {
			ok = false
			errNo = 9999
		}
	}
	if currentUserId != 0 {
		user.BeFollowed = GetUserIsFollowed(currentUserId, user.Id)
	}
	user = common.MakeUserSupInfo(user)
	return user, ok, errNo
}

func GetUserInfoByUsername(userName string, currentUserId uint) (model.User, bool) {
	var user model.User
	err := db.QueryRow("SELECT id, username, password, nickname, follow_count, be_followed_count, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, reg_time, last_login_time FROM user WHERE username=? LIMIT 1", userName).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.FollowCount, &user.BeFollowedCount, &user.BeLikedCount, &user.BeFavoriteCount, &user.BeCommentedCount, &user.BeForwardedCount, &user.BeWatchedCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return user, false
	}

	if currentUserId != 0 {
		user.BeFollowed = GetUserIsFollowed(currentUserId, user.Id)
	}
	user = common.MakeUserSupInfo(user)
	return user, true
}

func CheckUserPassword(username string, password string) (model.User, bool, int) {
	ok := true
	errNo := 0

	user, userExist := GetUserInfoByUsername(username, 0)

	if !userExist {
		ok = false
		errNo = 1
	} else if password != user.Password {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err:", "wrong password")
		}
		ok = false
		errNo = 2
	}
	// no need to check be followed
	user = common.MakeUserSupInfo(user)
	return user, ok, errNo
}

func CreateUser(username string, password string, nickname string) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	_, userExist := GetUserInfoByUsername(username, 0)

	if !userExist { // no such user
		// create user
		if nickname == "" { // empty nickname
			nickname = "用户" + strconv.Itoa(rand.Intn(100000)+100000) // random nickname
		}
		res, err := db.Exec("INSERT INTO user (username, password, nickname, reg_time) VALUES (?, ?, ?, NOW())",
			username, password, nickname)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			ok = false
			errNo = 9999
		} else {
			lastId, err := res.LastInsertId()
			if err != nil {
				if config.ShowLog {
					funcName, _, _, _ := runtime.Caller(0)
					log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
				}
				ok = false
				errNo = 9999
			} else {
				if config.ShowLog {
					funcName, _, _, _ := runtime.Caller(0)
					log.Println(runtime.FuncForPC(funcName).Name(), "signup lastId: ", lastId)
				}
				user.Id = uint(lastId)
				user.Username = username
				user.Password = password
				user.Nickname = nickname
			}
		}
	} else { // username exist
		ok = false
		errNo = 1
	}
	// no need to check be followed
	user = common.MakeUserSupInfo(user)
	return user, ok, errNo
}

func SetLastLogin(userId uint) {
	_, err := db.Exec("UPDATE user SET last_login_time=NOW() WHERE id=?", userId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
	}
}

func FollowUser(followUserId uint, userId uint, action string) bool {
	if action == "follow" {
		_, err := db.Exec("INSERT INTO user_follow (follow_user_id, user_id, follow_time) VALUES (?, ?, NOW())", followUserId, userId)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return false
		}
	} else if action == "unfollow" {
		_, err := db.Exec("UPDATE user_follow SET unfollow_time=NOW() WHERE follow_user_id=? AND user_id=?", followUserId, userId)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return false
		}
	} else {
		return false
	}
	return true
}

func SetUserInfo(userId uint, nickname string) bool {
	_, err := db.Exec("UPDATE user SET nickname=? WHERE id=?", nickname, userId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false
	}
	return true
}

func GetVideoList(qType int, qUserId uint, qRelation string, limit int, start int, currentUserId uint) []model.Video {
	// get video list
	var videoList []model.Video
	var rows *sql.Rows
	var err error
	if qUserId == 0 {
		if qType == 0 {
			if qRelation == "followed" {
				// get user's followed user's video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND user_id IN (SELECT follow_user_id FROM user_follow WHERE user_id=? AND unfollow_time IS NULL ORDER BY id DESC) ORDER BY id DESC LIMIT ?, ?", currentUserId, start, limit)
			} else if qRelation == "watched" {
				// get user's watched video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND id IN (SELECT DISTINCT(video_id) FROM video_watch WHERE user_id=? ORDER BY id DESC) LIMIT ?, ?", currentUserId, start, limit)
			} else {
				rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 ORDER BY id DESC LIMIT ?, ?", start, limit)
			}
		} else {
			rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND type=? ORDER BY id DESC LIMIT ?, ?", qType, start, limit)

		}
	} else {
		// specific user
		if qType == 0 {
			// all relations
			if qRelation == "liked" {
				// get user's liked video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND id IN (SELECT video_id FROM video_like WHERE user_id=? AND unlike_time IS NULL ORDER BY id DESC) LIMIT ?, ?", qUserId, start, limit)
			} else if qRelation == "favorite" {
				// get user's favorite video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND id IN (SELECT video_id FROM video_favorite WHERE user_id=? AND unfavorite_time IS NULL ORDER BY id DESC) LIMIT ?, ?", qUserId, start, limit)
			} else {
				// get user's all uploaded video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND user_id=? ORDER BY id DESC LIMIT ?, ?", qUserId, start, limit)
			}
		} else {
			// qActionHistory valid only when qType == 0, else ignored.
			// get user's uploaded video list by type
			rows, err = db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND user_id=? AND type=? ORDER BY id DESC LIMIT ?, ?", qUserId, qType, start, limit)
		}
	}

	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return videoList
	}
	defer rows.Close()

	for rows.Next() {
		var video model.Video
		var screenshotOk uint8
		var hlsOk uint8

		err := rows.Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Content, &video.Keyword, &video.UploadTime, &video.BeLikedCount, &video.BeFavoriteCount, &video.BeCommentedCount, &video.BeForwardedCount, &video.BeWatchedCount, &screenshotOk, &hlsOk)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return videoList
		}
		user, _, _ := GetUserInfoById(video.UserId, currentUserId)
		video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
		videoList = append(videoList, video)
	}
	return videoList
}

func GetVideoNum(qUserId uint, qRelation string) int {
	// get video num
	var err error
	var count int
	if qRelation == "like" {
		// get user's liked video num
		err = db.QueryRow("SELECT COUNT(1) FROM video_like WHERE user_id=? AND unlike_time IS NULL", qUserId).Scan(&count)
	} else if qRelation == "favorite" {
		// get user's favorite video num
		err = db.QueryRow("SELECT COUNT(1) FROM video_favorite WHERE user_id=? AND unfavorite_time IS NULL", qUserId).Scan(&count)
	} else if qRelation == "watched" {
		// get user's watched video num
		err = db.QueryRow("SELECT COUNT(DISTINCT(video_id)) FROM video_watch WHERE user_id=?", qUserId).Scan(&count)
	} else {
		// get user's all uploaded video num
		err = db.QueryRow("SELECT COUNT(1) FROM video WHERE deleted=0 AND user_id=?", qUserId).Scan(&count)
	}
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err:", err)
		}
		return 0
	}
	return count
}

func GetVideoIdByVideoUid(videoUid string) uint {
	var videoId uint
	err := db.QueryRow("SELECT id FROM video WHERE video_uid=? LIMIT 1", videoUid).Scan(&videoId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return 0
	}
	return videoId
}

func GetVideoDetailByVideoUid(videoUid string, currentUserId uint) model.Video {
	var video model.Video
	var screenshotOk uint8
	var hlsOk uint8

	err := db.QueryRow("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE video_uid=? LIMIT 1", videoUid).
		Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Content, &video.Keyword, &video.UploadTime, &video.BeLikedCount, &video.BeFavoriteCount, &video.BeCommentedCount, &video.BeForwardedCount, &video.BeWatchedCount, &screenshotOk, &hlsOk)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		if errors.Is(err, sql.ErrNoRows) { // no such video
			return video
		} else {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return video
		}
	}
	user, _, _ := GetUserInfoById(video.UserId, currentUserId)
	video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
	return video
}

func CheckUserVideoRelation(userId uint, videoId uint, relationType string) bool {
	switch relationType {
	case "liked":
		var count int
		err := db.QueryRow("SELECT COUNT(1) FROM video_like WHERE video_id=? AND user_id=? AND unlike_time IS NULL LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "favorite":
		var count int
		err := db.QueryRow("SELECT COUNT(1) FROM video_favorite WHERE video_id=? AND user_id=? AND unfavorite_time IS NULL LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "uploaded":
		var count int
		err := db.QueryRow("SELECT COUNT(1) FROM video WHERE id=? AND user_id=? LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "watched":
		var count int
		err := db.QueryRow("SELECT COUNT(1) FROM video_watch WHERE video_id=? AND user_id=? LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "last_play":
		var videoIdReturn uint
		err := db.QueryRow("SELECT video_id FROM video_watch WHERE user_id=? ORDER BY id DESC LIMIT 1", userId).Scan(&videoIdReturn)
		if err != nil || videoIdReturn != videoId {
			return false
		}
		return true
	default:
		return false
	}
}

func CheckUserVideoAllRelation(userId uint, videoId uint) (bool, bool, bool, bool, bool) {
	/*
	 * return: isUserLiked, isUserFavorite, isUserUploaded, isUserWatched, isUserLastPlay
	 */
	return CheckUserVideoRelation(userId, videoId, "liked"),
		CheckUserVideoRelation(userId, videoId, "favorite"),
		CheckUserVideoRelation(userId, videoId, "uploaded"),
		CheckUserVideoRelation(userId, videoId, "watched"),
		CheckUserVideoRelation(userId, videoId, "last_play")
}

func checkVideoActionHistory(videoId uint, userId uint, actionType string) (model.VideoAction, error) {
	var videoAction model.VideoAction
	var err error
	if actionType == "like" || actionType == "unlike" {
		err = db.QueryRow("SELECT id, video_id, user_id, like_time, unlike_time FROM video_like WHERE video_id=? AND user_id=? ORDER BY id DESC LIMIT 1", videoId, userId).
			Scan(&videoAction.Id, &videoAction.VideoId, &videoAction.UserId, &videoAction.ActionTime, &videoAction.UnActionTime)
	} else if actionType == "favorite" || actionType == "unfavorite" {
		err = db.QueryRow("SELECT id, video_id, user_id, favorite_time, unfavorite_time FROM video_favorite WHERE video_id=? AND user_id=? ORDER BY id DESC LIMIT 1", videoId, userId).
			Scan(&videoAction.Id, &videoAction.VideoId, &videoAction.UserId, &videoAction.ActionTime, &videoAction.UnActionTime)
	}
	return videoAction, err
}

func GuestDoVideoAction(videoId uint, userId uint, actionType string) (bool, int) {
	videoAction, err := checkVideoActionHistory(videoId, userId, actionType)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false, 1
	}
	switch actionType {
	case "like":
		if err == nil && !videoAction.UnActionTime.Valid {
			return false, 1 // already liked
		}

		if _, err := db.Exec("INSERT INTO video_like (video_id, user_id, like_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return false, 9999
		}
	case "unlike":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return false, 1 // already unliked
		}

		if _, err := db.Exec("UPDATE video_like SET unlike_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return false, 9999
		}
	case "favorite":
		if err == nil && !videoAction.UnActionTime.Valid {
			return false, 1 // already favorite
		}

		if _, err := db.Exec("INSERT INTO video_favorite (video_id, user_id, favorite_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return false, 9999
		}
	case "unfavorite":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return false, 1 // already unfavorite
		}

		if _, err := db.Exec("UPDATE video_favorite SET unfavorite_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return false, 9999
		}
	case "forward":

	default:
		return false, 9999
	}
	return true, 0
}

func GuestWatchedVideo(videoId uint, userId uint) bool {
	if _, err := db.Exec("INSERT INTO video_watch (video_id, user_id, watch_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false
	}
	return true
}

func GuestForwardVideo(videoId uint, userId uint) bool {
	if _, err := db.Exec("INSERT INTO video_forward (video_id, user_id, watch_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false
	}
	return true
}

func GetVideoCommentList(videoId uint, limit int, start int, currentUserId uint) []model.VideoComment {
	var videoComments []model.VideoComment
	// inner join video for video_uid
	rows, err := db.Query("SELECT video_comment.id, video.video_uid, video_comment.user_id, "+
		"video_comment.content, video_comment.quote_comment_id, video_comment.comment_time FROM video_comment "+
		"INNER JOIN video ON video_comment.video_id=video.id WHERE video_comment.video_id=? AND "+
		"video_comment.delete_time IS NULL ORDER BY video_comment.id DESC LIMIT ?, ?", videoId, start, limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return videoComments
	}
	defer rows.Close()

	for rows.Next() {
		var videoComment model.VideoComment
		err := rows.Scan(&videoComment.Id, &videoComment.VideoUid, &videoComment.UserId, &videoComment.Content, &videoComment.QuoteCommentId, &videoComment.CommentTime)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return videoComments
		}
		videoComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(videoComment.CommentTime)
		user, _, _ := GetUserInfoById(videoComment.UserId, currentUserId)
		videoComment.User = user
		videoComments = append(videoComments, videoComment)
	}
	return videoComments
}

func CheckVideoCommentDuplicate(videoId uint, userId uint, content string, quoteCommentId uint) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM video_comment WHERE user_id=? AND video_id=? AND content=? AND quote_comment_id=? LIMIT 1", userId, videoId, content, quoteCommentId).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func CheckVideoCommentExist(videoCommentId uint) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM video_comment WHERE id=? LIMIT 1", videoCommentId).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func UserMakeVideoComment(videoId uint, userId uint, content string, quoteCommentId uint) bool {
	var err error
	if quoteCommentId == 0 { // not quote
		_, err = db.Exec("INSERT INTO video_comment (video_id, user_id, content, comment_time) VALUES (?, ?, ?, NOW())", videoId, userId, content)
	} else {
		_, err = db.Exec("INSERT INTO video_comment (video_id, user_id, content, quote_comment_id, comment_time) VALUES (?, ?, ?, ?, NOW())", videoId, userId, content, quoteCommentId)
	}
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err, "videoId:", videoId, "userId:", userId, "content:", content, "quoteCommentId:", quoteCommentId)
		}
		return false
	}

	// update video comment count
	if _, err := db.Exec("UPDATE video SET comment_count=comment_count+1 WHERE id=?", videoId); err != nil {
		return false
	}

	return true
}

func GetVideoComment(videoCommentId uint, currentUserId uint) model.VideoComment {
	var videoComment model.VideoComment
	err := db.QueryRow("SELECT video_comment.id, video.video_uid, video_comment.user_id, "+
		"video_comment.content, video_comment.quote_comment_id, video_comment.comment_time FROM video_comment "+
		"INNER JOIN video ON video_comment.video_id=video.id WHERE video_comment.id=? AND "+
		"video_comment.delete_time IS NULL LIMIT 1", videoCommentId).Scan(&videoComment.Id, &videoComment.VideoUid, &videoComment.UserId, &videoComment.Content, &videoComment.QuoteCommentId, &videoComment.CommentTime)
	if err != nil {
		return videoComment
	}
	videoComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(videoComment.CommentTime)
	user, _, _ := GetUserInfoById(videoComment.UserId, currentUserId)
	videoComment.User = user
	return videoComment
}

func GuestDeleteVideoComment(videoCommentId uint) bool {
	if _, err := db.Exec("UPDATE video_comment SET delete_time=NOW() WHERE id=?", videoCommentId); err != nil {
		return false
	}
	return true
}

func CheckVideoDuplicate(videoContent string, videoKeyword string, videoType int, userId uint) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM video WHERE content=? AND keyword=? AND type=? AND user_id=? LIMIT 1", videoContent, videoKeyword, videoType, userId).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func CheckVideoType(videoType int) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM video_type WHERE id=? LIMIT 1", videoType).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func UserCreateVideo(userId uint, videoUid string, videoType int, videoContent string, videoKeyword string) (model.Video, bool) {
	var video model.Video
	if _, err := db.Exec("INSERT INTO video (user_id, video_uid, type, content, keyword, upload_time) VALUES (?, ?, ?, ?, ?, NOW())", userId, videoUid, videoType, videoContent, videoKeyword); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return video, false
	}

	// get video info after insert
	video = GetVideoDetailByVideoUid(videoUid, 0)

	return video, true
}

func SetVideoInfo(videoUid string, videoContent string, videoKeyword string, videoType int8) bool {
	if _, err := db.Exec("UPDATE video SET content=?, keyword=?, type=? WHERE video_uid=?", videoContent, videoKeyword, videoType, videoUid); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false
	}
	return true
}

func DeleteVideo(videoUid string) bool {
	if _, err := db.Exec("UPDATE video SET deleted=1 WHERE video_uid=?", videoUid); err != nil {
		return false
	}
	return true
}

func GetVideoTypes() []model.VideoType {
	var videoTypes []model.VideoType
	rows, err := db.Query("SELECT id, type_name FROM video_type")
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return videoTypes
	}
	defer rows.Close()

	for rows.Next() {
		var videoType model.VideoType
		err := rows.Scan(&videoType.Id, &videoType.TypeName)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return videoTypes
		}
		videoTypes = append(videoTypes, videoType)
	}
	return videoTypes
}

func SearchVideo(keyword string, limit int, start int, currentUserId uint) []model.Video {
	var videoList []model.Video
	// setup ngram before using search
	// `CREATE FULLTEXT INDEX ngram_content ON video (content) WITH PARSER ngram;`
	// `CREATE FULLTEXT INDEX ngram_keyword ON video (keyword) WITH PARSER ngram;`
	rows, err := db.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND (MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) OR MATCH(keyword) AGAINST(? IN NATURAL LANGUAGE MODE)) ORDER BY id DESC LIMIT ?, ?", keyword, keyword, start, limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return videoList
	}
	defer rows.Close()

	for rows.Next() {
		var video model.Video
		var screenshotOk uint8
		var hlsOk uint8

		err := rows.Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Content, &video.Keyword, &video.UploadTime, &video.BeLikedCount, &video.BeFavoriteCount, &video.BeCommentedCount, &video.BeForwardedCount, &video.BeWatchedCount, &screenshotOk, &hlsOk)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return videoList
		}
		user, _, _ := GetUserInfoById(video.UserId, currentUserId)
		video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
		videoList = append(videoList, video)
	}

	// log search history
	if _, err := db.Exec("INSERT INTO search_history_video (keyword, search_user_id, search_time) VALUES (?, ?, NOW())", keyword, currentUserId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "log search history err: ", err)
		}
	}

	return videoList
}

func GetSearchVideoHotkeys(maxCount int) []string {
	var hotkeys []string
	limit := 20
	if maxCount > limit {
		limit = maxCount
	}
	rows, err := db.Query("SELECT keyword FROM search_history_video GROUP BY keyword ORDER BY COUNT(1) DESC LIMIT ?", limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return hotkeys
	}
	defer rows.Close()

	for rows.Next() {
		var hotkey string
		err := rows.Scan(&hotkey)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return hotkeys
		}
		hotkeys = append(hotkeys, hotkey)
	}

	// randomly select {{maxCount}} of hotkeys
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(hotkeys), func(i, j int) { hotkeys[i], hotkeys[j] = hotkeys[j], hotkeys[i] })
	hotkeysLen := len(hotkeys)
	if hotkeysLen < maxCount {
		maxCount = hotkeysLen
	}
	hotkeys = hotkeys[:maxCount]
	return hotkeys
}

func SearchUser(keyword string, limit int, start int, currentUserId uint) []model.User {
	var userList []model.User
	rows, err := db.Query("SELECT id, username, password, nickname, follow_count, be_followed_count, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, reg_time, last_login_time FROM user WHERE username LIKE ? OR nickname LIKE ? ORDER BY last_login_time DESC LIMIT ?, ?", "%"+keyword+"%", "%"+keyword+"%", start, limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return userList
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.FollowCount, &user.BeFollowedCount, &user.BeLikedCount, &user.BeFavoriteCount, &user.BeCommentedCount, &user.BeForwardedCount, &user.BeWatchedCount, &user.RegTime, &user.LastLoginTime)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
		}
		if currentUserId != 0 {
			user.BeFollowed = GetUserIsFollowed(currentUserId, user.Id)
		}
		user = common.MakeUserSupInfo(user)
		userList = append(userList, user)
	}

	// log search history
	if _, err := db.Exec("INSERT INTO search_history_user (keyword, search_user_id, search_time) VALUES (?, ?, NOW())", keyword, currentUserId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "log search history err: ", err)
		}
	}

	return userList
}

func GetSearchUserHotkeys(maxCount int) []string {
	var hotkeys []string
	limit := 20
	if maxCount > limit {
		limit = maxCount
	}
	rows, err := db.Query("SELECT keyword FROM search_history_user GROUP BY keyword ORDER BY COUNT(1) DESC LIMIT ?", limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return hotkeys
	}
	defer rows.Close()

	for rows.Next() {
		var hotkey string
		err := rows.Scan(&hotkey)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			return hotkeys
		}
		hotkeys = append(hotkeys, hotkey)
	}

	// randomly select {{maxCount}} of hotkeys
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(hotkeys), func(i, j int) { hotkeys[i], hotkeys[j] = hotkeys[j], hotkeys[i] })
	hotkeysLen := len(hotkeys)
	if hotkeysLen < maxCount {
		maxCount = hotkeysLen
	}
	hotkeys = hotkeys[:maxCount]
	return hotkeys
}

func CallbackUpdateVideoHLS(videoUid string) bool {
	if _, err := db.Exec("UPDATE video SET hls=1 WHERE video_uid=?", videoUid); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false
	}
	return true
}

func CallbackUpdateVideoScreenshot(videoUid string) bool {
	if _, err := db.Exec("UPDATE video SET screenshot=1 WHERE video_uid=?", videoUid); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return false
	}
	return true
}
