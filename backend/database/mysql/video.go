package mysql

import (
	"backend/common"
	"backend/config"
	"backend/model"
	"database/sql"
	"errors"
	"log"
	"runtime"
)

func GetVideoList(videoType int, userId uint, relation string, limit int, start int, currentUserId uint) []model.Video {
	// get video list
	var videoList []model.Video
	var rows *sql.Rows
	var err error
	if userId == 0 {
		if videoType == 0 {
			if relation == "followed" {
				// get user's followed user's video list
				rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND user_id IN (SELECT follow_user_id FROM user_follow WHERE user_id=? AND unfollow_time IS NULL ORDER BY id DESC) ORDER BY id DESC LIMIT ?, ?", currentUserId, start, limit)
			} else if relation == "watched" {
				// get user's watched video list, fix wrong order on 20231104
				rows, err = DB.Query("SELECT v.id, v.user_id, v.video_uid, v.type, v.content, v.keyword, v.upload_time, v.be_liked_count, v.be_favorite_count, v.be_commented_count, v.be_forwarded_count, v.be_watched_count, v.top, v.private, v.screenshot, v.hls FROM (SELECT video_id FROM (SELECT video_id, MAX(id) AS max_id FROM video_watch WHERE user_id=? GROUP BY video_id ) vw ORDER BY max_id DESC LIMIT ?, ?) AS subquery JOIN video v ON v.id = subquery.video_id WHERE v.deleted=0 AND v.private=0 AND v.hls=1", currentUserId, start, limit)
			} else {
				// comprehensive video list, [deprecated on 202311107]
				rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 ORDER BY RAND() DESC LIMIT ?, ?", start, limit)
			}
		} else {
			// specific type video list
			rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND type=? ORDER BY RAND() DESC LIMIT ?, ?", videoType, start, limit)
		}
	} else {
		// specific user
		if videoType == 0 {
			// all relations
			if relation == "liked" {
				// get user's liked video list, fix wrong order on 20231104
				rows, err = DB.Query("SELECT v.id, v.user_id, v.video_uid, v.type, v.content, v.keyword, v.upload_time, v.be_liked_count, v.be_favorite_count, v.be_commented_count, v.be_forwarded_count, v.be_watched_count, v.top, v.private, v.screenshot, v.hls FROM (SELECT video_id FROM (SELECT video_id, MAX(id) AS max_id FROM video_like WHERE user_id=? AND unlike_time IS NULL GROUP BY video_id ) vl ORDER BY max_id DESC LIMIT ?, ?) AS subquery JOIN video v ON v.id = subquery.video_id WHERE v.deleted=0 AND v.private=0 AND v.hls=1", userId, start, limit)
			} else if relation == "favorite" {
				// get user's favorite video list, fix wrong order on 20231104
				rows, err = DB.Query("SELECT v.id, v.user_id, v.video_uid, v.type, v.content, v.keyword, v.upload_time, v.be_liked_count, v.be_favorite_count, v.be_commented_count, v.be_forwarded_count, v.be_watched_count, v.top, v.private, v.screenshot, v.hls FROM (SELECT video_id FROM (SELECT video_id, MAX(id) AS max_id FROM video_favorite WHERE user_id=? AND unfavorite_time IS NULL GROUP BY video_id ) vf ORDER BY max_id DESC LIMIT ?, ?) AS subquery JOIN video v ON v.id = subquery.video_id WHERE v.deleted=0 AND v.private=0 AND v.hls=1", userId, start, limit)
			} else {
				// get user's all uploaded video list, add top order, private check on 20231104
				if currentUserId != 0 && userId == currentUserId {
					rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND user_id=? ORDER BY (top IS NOT NULL) DESC, top DESC, id DESC LIMIT ?, ?", userId, start, limit)
				} else {
					// remove private video
					rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND user_id=? ORDER BY (top IS NOT NULL) DESC, top DESC, id DESC LIMIT ?, ?", userId, start, limit)
				}
			}
		} else {
			// qActionHistory valid only when videoType == 0, else ignored.
			// get user's uploaded video list by type
			// [seems useless now], uploaded video will not be list type
			if currentUserId != 0 && userId == currentUserId {
				rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND hls=1 AND user_id=? AND type=? ORDER BY id DESC LIMIT ?, ?", userId, videoType, start, limit)
			} else {
				// remove private video
				rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND user_id=? AND type=? ORDER BY id DESC LIMIT ?, ?", userId, videoType, start, limit)
			}
		}
	}

	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return videoList
	}
	defer rows.Close()

	for rows.Next() {
		var video model.Video
		var screenshotOk uint8
		var hlsOk uint8

		err := rows.Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Content, &video.Keyword, &video.UploadTime, &video.BeLikedCount, &video.BeFavoriteCount, &video.BeCommentedCount, &video.BeForwardedCount, &video.BeWatchedCount, &video.Top, &video.Private, &screenshotOk, &hlsOk)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return videoList
		}
		user, _, _ := GetUserInfoById(video.UserId, currentUserId)
		video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
		videoList = append(videoList, video)
	}
	return videoList
}

func GetVideoNum(userId uint, relation string, currentUserId uint) int {
	// get video num
	var err error
	var count int
	if relation == "liked" {
		// get user's liked video num
		err = DB.QueryRow("SELECT COUNT(1) FROM video_like WHERE user_id=? AND unlike_time IS NULL", userId).Scan(&count)
	} else if relation == "favorite" {
		// get user's favorite video num
		err = DB.QueryRow("SELECT COUNT(1) FROM video_favorite WHERE user_id=? AND unfavorite_time IS NULL", userId).Scan(&count)
	} else if relation == "watched" {
		// get user's watched video num
		err = DB.QueryRow("SELECT COUNT(DISTINCT(video_id)) FROM video_watch WHERE user_id=?", userId).Scan(&count)
	} else {
		// get user's all uploaded video num
		if currentUserId != 0 && userId == currentUserId {
			err = DB.QueryRow("SELECT COUNT(1) FROM video WHERE deleted=0 AND hls=1 AND user_id=?", userId).Scan(&count)
		} else {
			// remove private video
			err = DB.QueryRow("SELECT COUNT(1) FROM video WHERE deleted=0 AND private=0 AND hls=1 AND user_id=?", userId).Scan(&count)
		}
	}
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR:", err)
		}
		return 0
	}
	return count
}

func GetVideoIdByVideoUid(videoUid string) uint {
	var videoId uint
	err := DB.QueryRow("SELECT id FROM video WHERE video_uid=? LIMIT 1", videoUid).Scan(&videoId)
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

	err := DB.QueryRow("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE video_uid=? LIMIT 1", videoUid).
		Scan(&video.Id, &video.UserId, &video.VideoUid, &video.Type, &video.Content, &video.Keyword, &video.UploadTime, &video.BeLikedCount, &video.BeFavoriteCount, &video.BeCommentedCount, &video.BeForwardedCount, &video.BeWatchedCount, &video.Top, &video.Private, &screenshotOk, &hlsOk)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		if errors.Is(err, sql.ErrNoRows) { // no such video
			return video
		} else {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return video
		}
	}
	user, _, _ := GetUserInfoById(video.UserId, currentUserId)
	video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
	return video
}

func CheckUserVideoRelation(userId uint, videoId uint, relation string) bool {
	switch relation {
	case "liked":
		var count int
		err := DB.QueryRow("SELECT COUNT(1) FROM video_like WHERE video_id=? AND user_id=? AND unlike_time IS NULL LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "favorite":
		var count int
		err := DB.QueryRow("SELECT COUNT(1) FROM video_favorite WHERE video_id=? AND user_id=? AND unfavorite_time IS NULL LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "uploaded":
		var count int
		err := DB.QueryRow("SELECT COUNT(1) FROM video WHERE id=? AND user_id=? AND deleted=0 LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "watched":
		var count int
		err := DB.QueryRow("SELECT COUNT(1) FROM video_watch WHERE video_id=? AND user_id=? LIMIT 1", videoId, userId).Scan(&count)
		if err != nil || count == 0 {
			return false
		}
		return true
	case "last_play":
		var videoIdReturn uint
		err := DB.QueryRow("SELECT video_id FROM video_watch WHERE user_id=? ORDER BY id DESC LIMIT 1", userId).Scan(&videoIdReturn)
		if err != nil || videoIdReturn != videoId {
			return false
		}
		return true
	case "top":
		var top sql.NullTime
		err := DB.QueryRow("SELECT top FROM video WHERE id=? LIMIT 1", videoId).Scan(&top)
		if err != nil || !top.Valid {
			return false
		}
		return true
	default:
		return false
	}
}

func CheckUserVideoAllRelation(userId uint, videoId uint) (bool, bool, bool, bool, bool) {
	/*
	 * return: isUserLiked, isUserFavorite, isUserUploaded, isUserWatched, isUserLastPlay, isTop
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
		err = DB.QueryRow("SELECT id, video_id, user_id, like_time, unlike_time FROM video_like WHERE video_id=? AND user_id=? ORDER BY id DESC LIMIT 1", videoId, userId).
			Scan(&videoAction.Id, &videoAction.VideoId, &videoAction.UserId, &videoAction.ActionTime, &videoAction.UnActionTime)
	} else if actionType == "favorite" || actionType == "unfavorite" {
		err = DB.QueryRow("SELECT id, video_id, user_id, favorite_time, unfavorite_time FROM video_favorite WHERE video_id=? AND user_id=? ORDER BY id DESC LIMIT 1", videoId, userId).
			Scan(&videoAction.Id, &videoAction.VideoId, &videoAction.UserId, &videoAction.ActionTime, &videoAction.UnActionTime)
	}
	return videoAction, err
}

func GuestLikeOrFavoriteVideo(videoId uint, userId uint, actionType string) (bool, int) {
	videoAction, err := checkVideoActionHistory(videoId, userId, actionType)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false, 1
	}
	switch actionType {
	case "like":
		if err == nil && !videoAction.UnActionTime.Valid {
			return false, 1 // already liked
		}

		if _, err := DB.Exec("INSERT INTO video_like (video_id, user_id, like_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return false, 9999
		}
	case "unlike":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return false, 1 // already unliked
		}

		if _, err := DB.Exec("UPDATE video_like SET unlike_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return false, 9999
		}
	case "favorite":
		if err == nil && !videoAction.UnActionTime.Valid {
			return false, 1 // already favorite
		}

		if _, err := DB.Exec("INSERT INTO video_favorite (video_id, user_id, favorite_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return false, 9999
		}
	case "unfavorite":
		if errors.Is(err, sql.ErrNoRows) || videoAction.UnActionTime.Valid {
			return false, 1 // already unfavorite
		}

		if _, err := DB.Exec("UPDATE video_favorite SET unfavorite_time=NOW() WHERE id=?", videoAction.Id); err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return false, 9999
		}
	default:
		return false, 9999
	}
	return true, 0
}

func GuestWatchVideo(videoId uint, userId uint) bool {
	if _, err := DB.Exec("INSERT INTO video_watch (video_id, user_id, watch_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func GuestForwardVideo(videoId uint, userId uint) bool {
	if _, err := DB.Exec("INSERT INTO video_forward (video_id, user_id, forward_time) VALUES (?, ?, NOW())", videoId, userId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func CheckVideoDuplicate(videoContent string, videoKeyword string, videoType int, userId uint) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(1) FROM video WHERE content=? AND keyword=? AND type=? AND user_id=? LIMIT 1", videoContent, videoKeyword, videoType, userId).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func CheckVideoType(videoType int) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(1) FROM video_type WHERE id=? LIMIT 1", videoType).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func CreateVideoWithoutInfo(userId uint, videoUid string) bool {
	_, err := DB.Exec("INSERT INTO video (user_id, video_uid, upload_time) VALUES (?, ?, NOW())", userId, videoUid)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func ConfirmCreateVideoWithInfo(videoUid string, videoType int, videoContent string, videoKeyword string) bool {
	if _, err := DB.Exec("UPDATE video SET content=?, keyword=?, type=? WHERE video_uid=?", videoContent, videoKeyword, videoType, videoUid); err != nil {
		return false
	}
	return true
}

func CreateVideo(userId uint, videoUid string, videoType int, videoContent string, videoKeyword string) (model.Video, bool) {
	var video model.Video
	if _, err := DB.Exec("INSERT INTO video (user_id, video_uid, type, content, keyword, upload_time) VALUES (?, ?, ?, ?, ?, NOW())", userId, videoUid, videoType, videoContent, videoKeyword); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return video, false
	}

	// get video info after insert
	video = GetVideoDetailByVideoUid(videoUid, 0)

	return video, true
}

func SetVideoInfo(videoUid string, videoContent string, videoKeyword string, videoType int8) bool {
	if _, err := DB.Exec("UPDATE video SET content=?, keyword=?, type=? WHERE video_uid=?", videoContent, videoKeyword, videoType, videoUid); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func DeleteVideo(videoUid string) bool {
	if _, err := DB.Exec("UPDATE video SET deleted=1 WHERE video_uid=?", videoUid); err != nil {
		return false
	}
	return true
}

func UserTopOrPrivateVideo(videoUid string, actionType string) bool {
	if actionType == "top" {
		if _, err := DB.Exec("UPDATE video SET top=now() WHERE video_uid=?", videoUid); err != nil {
			return false
		}
	} else if actionType == "untop" {
		if _, err := DB.Exec("UPDATE video SET top=NULL WHERE video_uid=?", videoUid); err != nil {
			return false
		}
	} else if actionType == "private" {
		if _, err := DB.Exec("UPDATE video SET private=1 WHERE video_uid=?", videoUid); err != nil {
			return false
		}
	} else if actionType == "unprivate" {
		if _, err := DB.Exec("UPDATE video SET private=0 WHERE video_uid=?", videoUid); err != nil {
			return false
		}
	}
	return true
}

func GetVideoTypes() []model.VideoType {
	var videoTypes []model.VideoType
	rows, err := DB.Query("SELECT id, type_name FROM video_type")
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
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
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			return videoTypes
		}
		videoTypes = append(videoTypes, videoType)
	}
	return videoTypes
}

func GetRemainPromoteAndAdvertiseCount(videoId uint) (int, int) {
	var remainPromoteCount int
	var remainAdvertiseCount int
	err := DB.QueryRow("SELECT SUM(remain_count) FROM video_promote WHERE video_id=? AND remain_count>0", videoId).Scan(&remainPromoteCount)
	if err != nil {
		remainPromoteCount = 0
	}
	err = DB.QueryRow("SELECT SUM(remain_count) FROM video_advertise WHERE video_id=? AND remain_count>0", videoId).Scan(&remainAdvertiseCount)
	if err != nil {
		remainAdvertiseCount = 0
	}
	return remainPromoteCount, remainAdvertiseCount
}

func PromoteVideo(videoId uint, userId uint, count int) bool {
	price := float64(count) * config.PromotePrice
	if _, err := DB.Exec("INSERT INTO video_promote (video_id, user_id, count, remain_count, price, promote_time) VALUES (?, ?, ?, ?, ?, NOW())", videoId, userId, count, count, price); err != nil {
		return false
	}
	return true
}

func AdvertiseVideo(videoId uint, userId uint, count int) bool {
	price := float64(count) * config.AdvertisePrice
	if _, err := DB.Exec("INSERT INTO video_advertise (video_id, user_id, count, remain_count, price, advertise_time) VALUES (?, ?, ?, ?, ?, NOW())", videoId, userId, count, count, price); err != nil {
		return false
	}
	return true
}
