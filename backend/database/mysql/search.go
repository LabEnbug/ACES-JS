package mysql

import (
	"backend/common"
	"backend/config"
	"backend/model"
	"log"
	"math/rand"
	"runtime"
	"time"
)

func SearchVideo(keyword string, limit int, start int, currentUserId uint) []model.Video {
	var videoList []model.Video
	// setup ngram before using search
	// `CREATE FULLTEXT INDEX ngram_content ON video (content) WITH PARSER ngram;`
	// `CREATE FULLTEXT INDEX ngram_keyword ON video (keyword) WITH PARSER ngram;`
	rows, err := DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND content IS NOT NULL AND hls=1 AND (MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) OR MATCH(keyword) AGAINST(? IN NATURAL LANGUAGE MODE)) ORDER BY id DESC LIMIT ?, ?", keyword, keyword, start, limit)
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

	// log search history
	if _, err := DB.Exec("INSERT INTO search_history_video (keyword, search_user_id, search_time) VALUES (?, ?, NOW())", keyword, currentUserId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "log search history ERR: ", err)
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
	rows, err := DB.Query("SELECT keyword FROM search_history_video GROUP BY keyword ORDER BY COUNT(1) DESC LIMIT ?", limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
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
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
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
	rows, err := DB.Query("SELECT id, username, password, nickname, follow_count, be_followed_count, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, reg_time, last_login_time FROM user WHERE username LIKE ? OR nickname LIKE ? ORDER BY last_login_time DESC LIMIT ?, ?", "%"+keyword+"%", "%"+keyword+"%", start, limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
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
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
		}
		if currentUserId != 0 {
			user.BeFollowed = CheckUserFollowed(user.Id, currentUserId)
			user.IsSelf = user.Id == currentUserId
		}
		user = common.MakeUserSupInfo(user)
		userList = append(userList, user)
	}

	// log search history
	if _, err := DB.Exec("INSERT INTO search_history_user (keyword, search_user_id, search_time) VALUES (?, ?, NOW())", keyword, currentUserId); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "log search history ERR: ", err)
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
	rows, err := DB.Query("SELECT keyword FROM search_history_user GROUP BY keyword ORDER BY COUNT(1) DESC LIMIT ?", limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
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
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
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
