package database

import (
	"backend/config"
	"backend/model"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
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
		log.Println(err)
		panic(err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(time.Minute * 60) // mysql default conn timeout=8h, should < mysql_timeout
	err := db.Ping()
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func CloseDBPool() {
	err := db.Close()
	if err != nil {
		return
	}
}

func GetUserInfoById(userId int) (model.User, int) {
	var user model.User
	errNo := 0

	err := db.QueryRow("SELECT id, username, password, nickname, like_count, favorite_count, comment_count, reg_time, last_login_time FROM user WHERE id=? LIMIT 1", userId).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.LikeCount, &user.FavoriteCount, &user.CommentCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // no such user
			errNo = 1
		} else {
			errNo = 9999
		}
	}
	return user, errNo
}

func GetUserInfoByUsername(userName string) (model.User, int) {
	var user model.User
	errNo := 0

	err := db.QueryRow("SELECT id, username, password, nickname, like_count, favorite_count, comment_count, reg_time, last_login_time FROM user WHERE username=? LIMIT 1", userName).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.LikeCount, &user.FavoriteCount, &user.CommentCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // no such user
			errNo = 1
		} else {
			errNo = 9999
		}
	}
	return user, errNo
}

func CheckUserIfExist(username string) (bool, model.User) {
	user, errNo := GetUserInfoByUsername(username)
	if errNo != 0 {
		return false, user
	}
	return true, user
}

func CheckUserPassword(username string, password string) (model.User, int) {
	var user model.User
	errNo := 0
	userExist := false

	userExist, user = CheckUserIfExist(username)

	if !userExist {
		errNo = 1
	} else if password != user.Password {
		log.Println("err: wrong password")
		errNo = 2
	}
	return user, errNo
}

func CreateUser(username string, password string) (model.User, int) {
	var user model.User
	errNo := 0

	userExist, _ := CheckUserIfExist(username)

	if !userExist { // no such user
		// create user
		nickname := "用户" + strconv.Itoa(rand.Intn(100000)+100000) //random nickname
		res, err := db.Exec("INSERT INTO user (username, password, nickname, reg_time) VALUES (?, ?, ?, NOW())",
			username, password, nickname)
		if err != nil {
			log.Println("err: ", err)
			errNo = 9999
		} else {
			lastId, err := res.LastInsertId()
			if err != nil {
				log.Println("err: ", err)
				errNo = 9999
			} else {
				log.Println("signup lastId: ", lastId)
				user.Id = int(lastId)
				user.Username = username
				user.Password = password
				user.Nickname = nickname
			}
		}
	} else { // username exist
		errNo = 1
	}
	return user, errNo
}

func SetLastLogin(userId int) {
	_, err := db.Exec("UPDATE user SET last_login_time=NOW() WHERE id=?", userId)
	if err != nil {
		log.Println("err: ", err)
	}
}

func SetUserInfo(userId int, nickname string) int {
	_, err := db.Exec("UPDATE user SET nickname=? WHERE id=?", nickname, userId)
	if err != nil {
		log.Println("err: ", err)
		return 9999
	}
	return 0
}

func GetVideoList(qType int, qUser int, qActionHistory string, limit int, page int) []model.Video {
	// get video list
	var videoList []model.Video
	var rows *sql.Rows
	var err error
	if qUser == 0 {
		if qType == 0 {
			rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video ORDER BY id DESC LIMIT ?, ?", (page-1)*limit, limit)
		} else {
			rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE type=? ORDER BY id DESC LIMIT ?, ?", qType, (page-1)*limit, limit)
		}
	} else {
		// specific user
		if qType == 0 {
			// all types
			if qActionHistory == "like" {
				// get user's liked video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE id IN (SELECT video_id FROM video_like WHERE user_id=? AND unlike_time IS NULL ORDER BY id DESC) LIMIT ?, ?", qUser, (page-1)*limit, limit)
			} else if qActionHistory == "favorite" {
				// get user's favorite video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE id IN (SELECT video_id FROM video_favorite WHERE user_id=? AND unfavorite_time IS NULL ORDER BY id DESC) LIMIT ?, ?", qUser, (page-1)*limit, limit)
			} else if qActionHistory == "history" {
				// get user's history video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE id IN (SELECT video_id FROM video_history WHERE user_id=? ORDER BY id DESC) LIMIT ?, ?", qUser, (page-1)*limit, limit)
			} else {
				// get user's all uploaded video list
				rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE user_id=? LIMIT ?, ?", qUser, (page-1)*limit, limit)
			}
		} else {
			// qActionHistory valid only when qType == 0, else ignored.
			// get user's uploaded video list by type
			rows, err = db.Query("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE user_id=? AND type=? ORDER BY id DESC LIMIT ?, ?", qUser, qType, (page-1)*limit, limit)
		}
	}

	if err != nil {
		log.Println("err: ", err)
		return videoList
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var video model.Video
		err := rows.Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Title, &video.Content, &video.UploadTime, &video.LikeCount, &video.FavoriteCount, &video.CommentCount)
		if err != nil {
			log.Println("err: ", err)
			return videoList
		}
		videoList = append(videoList, video)
	}
	return videoList
}

func GetVideoDetail(videoUid string) model.Video {
	var video model.Video
	err := db.QueryRow("SELECT id, user_id, video_uid, type, title, content, upload_time, like_count, favorite_count, comment_count FROM video WHERE video_uid=? LIMIT 1", videoUid).
		Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Title, &video.Content, &video.UploadTime, &video.LikeCount, &video.FavoriteCount, &video.CommentCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // no such video
			return video
		} else {
			log.Println("err: ", err)
			return video
		}
	}
	return video
}

func checkVideoActionHistory(videoId int, userId int, actionType string) (model.VideoAction, error) {
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

func GuestDoVideoAction(video model.Video, userId int, actionType string) int {
	videoAction, err := checkVideoActionHistory(video.Id, userId, actionType)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("err: ", err)
		return 9999
	}

	switch actionType {
	case "like":
		if err == nil && !videoAction.UnActionTime.Valid {
			return 1 // already liked
		}

		if _, err := db.Exec("INSERT INTO video_like (video_id, user_id, like_time) VALUES (?, ?, NOW())", video.Id, userId); err != nil {
			log.Println("err: ", err)
			return 9999
		}

		if _, err := db.Exec("UPDATE video SET like_count=like_count+1 WHERE id=?", video.Id); err != nil {
			log.Println("err: ", err)
			return 9999
		}

	case "unlike":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return 1 // already unliked
		}

		if _, err := db.Exec("UPDATE video_like SET unlike_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			log.Println("err: ", err)
			return 9999
		}

		if _, err := db.Exec("UPDATE video SET like_count=like_count-1 WHERE id=?", video.Id); err != nil {
			log.Println("err: ", err)
			return 9999
		}
	case "favorite":
		if err == nil && !videoAction.UnActionTime.Valid {
			return 1 // already favorite
		}

		if _, err := db.Exec("INSERT INTO video_favorite (video_id, user_id, favorite_time) VALUES (?, ?, NOW())", video.Id, userId); err != nil {
			log.Println("err: ", err)
			return 9999
		}

		if _, err := db.Exec("UPDATE video SET favorite_count=favorite_count+1 WHERE id=?", video.Id); err != nil {
			log.Println("err: ", err)
			return 9999
		}

	case "unfavorite":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return 1 // already unfavorite
		}

		if _, err := db.Exec("UPDATE video_favorite SET unfavorite_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			log.Println("err: ", err)
			return 9999
		}

		if _, err := db.Exec("UPDATE video SET favorite_count=favorite_count-1 WHERE id=?", video.Id); err != nil {
			log.Println("err: ", err)
			return 9999
		}
	default:
		return 9999
	}
	return 0
}

func GuestRecordVideoHistory(video model.Video, userId int) int {
	if _, err := db.Exec("INSERT INTO video_history (video_id, user_id, record_time) VALUES (?, ?, NOW())", video.Id, userId); err != nil {
		log.Println("err: ", err)
		return 9999
	}
	return 0
}

func SetVideoInfo(videoUid string, videoTitle string, videoContent string, videoType int) int {
	if _, err := db.Exec("UPDATE video SET title=?, content=?, type=? WHERE video_uid=?", videoTitle, videoContent, videoType, videoUid); err != nil {
		log.Println("err: ", err)
		return 9999
	}
	return 0
}
