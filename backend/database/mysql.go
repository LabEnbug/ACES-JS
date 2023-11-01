package database

import (
	"backend/common"
	"backend/config"
	"backend/model"
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
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err:", err)
		panic(err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err := db.Ping()
	if err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err:", err)
		panic(err)
	}
}

func CloseDBPool() {
	err := db.Close()
	if err != nil {
		return
	}
}

func GetUserInfoById(userId uint) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	err := db.QueryRow("SELECT id, username, password, nickname, like_count, favorite_count, comment_count, reg_time, last_login_time FROM user WHERE id=? LIMIT 1", userId).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.LikeCount, &user.FavoriteCount, &user.CommentCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // no such user
			ok = false
			errNo = 1
		} else {
			ok = false
			errNo = 9999
		}
	}
	return user, ok, errNo
}

func GetUserInfoByUsername(userName string) (model.User, bool) {
	var user model.User
	err := db.QueryRow("SELECT id, username, password, nickname, like_count, favorite_count, comment_count, reg_time, last_login_time FROM user WHERE username=? LIMIT 1", userName).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.LikeCount, &user.FavoriteCount, &user.CommentCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		return user, false
	}
	return user, true
}

func CheckUserPassword(username string, password string) (model.User, bool, int) {
	ok := true
	errNo := 0

	user, userExist := GetUserInfoByUsername(username)

	if !userExist {
		ok = false
		errNo = 1
	} else if password != user.Password {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err:", "wrong password")
		ok = false
		errNo = 2
	}
	return user, ok, errNo
}

func CreateUser(username string, password string, nickname string) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	_, userExist := GetUserInfoByUsername(username)

	if !userExist { // no such user
		// create user
		if nickname == "" { // empty nickname
			nickname = "用户" + strconv.Itoa(rand.Intn(100000)+100000) // random nickname
		}
		res, err := db.Exec("INSERT INTO user (username, password, nickname, reg_time) VALUES (?, ?, ?, NOW())",
			username, password, nickname)
		if err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			ok = false
			errNo = 9999
		} else {
			lastId, err := res.LastInsertId()
			if err != nil {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
				ok = false
				errNo = 9999
			} else {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "signup lastId: ", lastId)
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
	return user, ok, errNo
}

func SetLastLogin(userId uint) {
	_, err := db.Exec("UPDATE user SET last_login_time=NOW() WHERE id=?", userId)
	if err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
	}
}

func SetUserInfo(userId uint, nickname string) bool {
	_, err := db.Exec("UPDATE user SET nickname=? WHERE id=?", nickname, userId)
	if err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		return false
	}
	return true
}

func GetVideoList(qType int, qUserId uint, qActionHistory string, limit int, page int) []model.Video {
	// get video list
	var videoList []model.Video
	var rows *sql.Rows
	var err error
	if qUserId == 0 {
		if qType == 0 {
			rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 ORDER BY id DESC LIMIT ?, ?", (page-1)*limit, limit)
		} else {
			rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND type=? ORDER BY id DESC LIMIT ?, ?", qType, (page-1)*limit, limit)
		}
	} else {
		// specific user
		if qType == 0 {
			// all types
			if qActionHistory == "like" {
				// get user's liked video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND id IN (SELECT video_id FROM video_like WHERE user_id=? AND unlike_time IS NULL ORDER BY id DESC) LIMIT ?, ?", qUserId, (page-1)*limit, limit)
			} else if qActionHistory == "favorite" {
				// get user's favorite video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND id IN (SELECT video_id FROM video_favorite WHERE user_id=? AND unfavorite_time IS NULL ORDER BY id DESC) LIMIT ?, ?", qUserId, (page-1)*limit, limit)
			} else if qActionHistory == "history" {
				// get user's history video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND id IN (SELECT video_id FROM video_history WHERE user_id=? ORDER BY id DESC) LIMIT ?, ?", qUserId, (page-1)*limit, limit)
			} else {
				// get user's all uploaded video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND user_id=? LIMIT ?, ?", qUserId, (page-1)*limit, limit)
			}
		} else {
			// qActionHistory valid only when qType == 0, else ignored.
			// get user's uploaded video list by type
			rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE deleted=0 AND user_id=? AND type=? ORDER BY id DESC LIMIT ?, ?", qUserId, qType, (page-1)*limit, limit)
		}
	}

	if err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		return videoList
	}
	defer rows.Close()

	for rows.Next() {
		var video model.Video
		var screenshotOk uint8
		var hlsOk uint8

		err := rows.Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Title, &video.Content, &video.UploadTime, &video.LikeCount, &video.FavoriteCount, &video.CommentCount, &screenshotOk, &hlsOk)
		if err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return videoList
		}
		if screenshotOk == 1 {
			video.CoverUrl = fmt.Sprintf("http://s348vstvo.bkt.clouddn.com/video/%s/cover.png", video.VideoUid)
		}
		if hlsOk == 1 {
			video.PlayUrl = fmt.Sprintf("http://s348vstvo.bkt.clouddn.com/video/%s/index.m3u8", video.VideoUid)
		}
		user, _, _ := GetUserInfoById(video.UserId)
		video.User = common.GetVisibleUserInfo(user)
		videoList = append(videoList, video)
	}
	return videoList
}

func GetVideoDetail(videoUid string) model.Video {
	var video model.Video
	var screenshotOk uint8
	var hlsOk uint8

	err := db.QueryRow("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count, screenshot, hls FROM video WHERE video_uid=? LIMIT 1", videoUid).
		Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Title, &video.Content, &video.UploadTime, &video.LikeCount, &video.FavoriteCount, &video.CommentCount, &screenshotOk, &hlsOk)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // no such video
			return video
		} else {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return video
		}
	}
	if screenshotOk == 1 {
		video.CoverUrl = fmt.Sprintf("http://s348vstvo.bkt.clouddn.com/video/%s/cover.png", video.VideoUid)
	}
	if hlsOk == 1 {
		video.PlayUrl = fmt.Sprintf("http://s348vstvo.bkt.clouddn.com/video/%s/index.m3u8", video.VideoUid)
	}
	user, _, _ := GetUserInfoById(video.UserId)
	video.User = common.GetVisibleUserInfo(user)
	return video
}

func CheckUserVideoRelation(userId uint, videoId uint, relationType string) bool {
	switch relationType {
	case "like":
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
	case "history":
		var count int
		err := db.QueryRow("SELECT COUNT(1) FROM video_history WHERE video_id=? AND user_id=? LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "last_play":
		var videoIdReturn uint
		err := db.QueryRow("SELECT video_id FROM video_history WHERE user_id=? ORDER BY id DESC LIMIT 1", videoId, userId).Scan(&videoIdReturn)
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
	 * return: isUserLiked, isUserFavorite, isUserUploaded, isUserHistory, isUserLastPlay
	 */
	return CheckUserVideoRelation(userId, videoId, "like"),
		CheckUserVideoRelation(userId, videoId, "favorite"),
		CheckUserVideoRelation(userId, videoId, "uploaded"),
		CheckUserVideoRelation(userId, videoId, "history"),
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

func GuestDoVideoAction(video model.Video, userId uint, actionType string) (bool, int) {
	videoAction, err := checkVideoActionHistory(video.Id, userId, actionType)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		return false, 1
	}

	switch actionType {
	case "like":
		if err == nil && !videoAction.UnActionTime.Valid {
			return false, 1 // already liked
		}

		if _, err := db.Exec("INSERT INTO video_like (video_id, user_id, like_time) VALUES (?, ?, NOW())", video.Id, userId); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

		if _, err := db.Exec("UPDATE video SET like_count=like_count+1 WHERE id=?", video.Id); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

	case "unlike":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return false, 1 // already unliked
		}

		if _, err := db.Exec("UPDATE video_like SET unlike_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

		if _, err := db.Exec("UPDATE video SET like_count=like_count-1 WHERE id=?", video.Id); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

	case "favorite":
		if err == nil && !videoAction.UnActionTime.Valid {
			return false, 1 // already favorite
		}

		if _, err := db.Exec("INSERT INTO video_favorite (video_id, user_id, favorite_time) VALUES (?, ?, NOW())", video.Id, userId); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

		if _, err := db.Exec("UPDATE video SET favorite_count=favorite_count+1 WHERE id=?", video.Id); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

	case "unfavorite":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return false, 1 // already unfavorite
		}

		if _, err := db.Exec("UPDATE video_favorite SET unfavorite_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}

		if _, err := db.Exec("UPDATE video SET favorite_count=favorite_count-1 WHERE id=?", video.Id); err != nil {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			return false, 9999
		}
	default:
		return false, 9999
	}
	return true, 0
}

func GuestRecordVideoHistory(video model.Video, userId uint) bool {
	if _, err := db.Exec("INSERT INTO video_history (video_id, user_id, record_time) VALUES (?, ?, NOW())", video.Id, userId); err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
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

func UserCreateVideo(userId uint, videoUid string, videoType int, videoTitle string, videoContent string) (model.Video, bool) {
	var video model.Video
	if _, err := db.Exec("INSERT INTO video (user_id, video_uid, type, title, content, upload_time) VALUES (?, ?, ?, ?, ?, NOW())", userId, videoUid, videoType, videoTitle, videoContent); err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		return video, false
	}
	// get video info after insert
	video = GetVideoDetail(videoUid)

	return video, true
}

func SetVideoInfo(videoUid string, videoTitle string, videoContent string, videoType int8) bool {
	if _, err := db.Exec("UPDATE video SET title=?, content=?, type=? WHERE video_uid=?", videoTitle, videoContent, videoType, videoUid); err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
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
		return videoTypes
	}
	defer rows.Close()

	for rows.Next() {
		var videoType model.VideoType
		err := rows.Scan(&videoType.Id, &videoType.TypeName)
		if err != nil {
			return videoTypes
		}
		videoTypes = append(videoTypes, videoType)
	}
	return videoTypes
}

func CallbackUpdateVideoHLS(videoUid string) bool {
	if _, err := db.Exec("UPDATE video SET hls=1 WHERE video_uid=?", videoUid); err != nil {
		return false
	}
	return true
}

func CallbackUpdateVideoScreenshot(videoUid string) bool {
	if _, err := db.Exec("UPDATE video SET screenshot=1 WHERE video_uid=?", videoUid); err != nil {
		return false
	}
	return true
}
