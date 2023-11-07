package mysql

import (
	"backend/common"
	"backend/config"
	"backend/model"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"runtime"
	"fmt"
	"strings"
)

func SetUserRecommendMatrix(userId uint, recommendMatrix []byte) bool {
	_, err := DB.Exec("UPDATE user SET recommend_matrix=? WHERE id=?", recommendMatrix, userId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func GetRecommendMatrixByUserId(userId uint) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	err := DB.QueryRow("SELECT id, recommend_matrix FROM user WHERE id=? LIMIT 1", userId).
		Scan(&user.Id, &user.RecommendMatrix)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
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

func GetOneRecommendVideoByProbabilityMatrix(recomType int, choice int, limit int, start int) []model.Video {
	// get one recommend video, but as list return 
	var videoList []model.Video
	var rows *sql.Rows
	var err error

	switch recomType {
	case 1:
		videoType := choice
		rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND type=? ORDER BY RAND() DESC LIMIT ?, ?", videoType, start, limit)
	case 2:
		userId := uint(choice)
		rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND user_id=? ORDER BY RAND() DESC LIMIT ?, ?", userId, start, limit)
	case 3:
		switch choice {
		// watched + liked + favorite + comment + forwarded
		case 1:
			watchedThreshold := uint(20)  // assume 20 watched is a popular video, cause the num of user is small. same as follows
			rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND be_watched_count>? ORDER BY RAND() DESC LIMIT ?, ?", watchedThreshold, start, limit)
		case 2:
			likedThreshold := uint(10)
			rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND be_liked_count>? ORDER BY RAND() DESC LIMIT ?, ?", likedThreshold, start, limit)
		case 3:
			favoriteThreshold := uint(5)
			rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND be_favorite_count>? ORDER BY RAND() DESC LIMIT ?, ?", favoriteThreshold, start, limit)
		case 4:
			commentThreshold := uint(2)
			rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND be_watched_count>? ORDER BY RAND() DESC LIMIT ?, ?", commentThreshold, start, limit)
		case 5:
			forwardThreshold := uint(2)
			rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND be_watched_count>? ORDER BY RAND() DESC LIMIT ?, ?", forwardThreshold, start, limit)
		}
	case 4:
		// traffic pool
		row := DB.QueryRow("SELECT video_id FROM video_promote WHERE remain_count>0 GROUP BY video_id limit 1")
		var videoId uint
		err = row.Scan(&videoId)
		if err != nil {
			log.Fatal(err)
		}
		rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND id=? ORDER BY RAND() DESC LIMIT ?, 1", videoId, start)
	case 5:
		// advertise pool
		row := DB.QueryRow("SELECT video_id FROM video_advertise WHERE remain_count>0 GROUP BY video_id limit 1")
		var videoId uint
		err = row.Scan(&videoId)
		if err != nil {
			log.Fatal(err)
		}
		rows, err = DB.Query("SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND id=? ORDER BY RAND() DESC LIMIT ?, 1", videoId, start)
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
		currentUserId := uint(0)
		user, _, _ := GetUserInfoById(video.UserId, currentUserId)
		video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
		videoList = append(videoList, video)
	}
	return videoList
}

// 生成指定数量的占位符字符串，用于构建查询语句
func generatePlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func GetRecommendVideoBySimilarity(videoId uint, requiredVideo int, start int) []model.Video {
	var videoList []model.Video
	var rows *sql.Rows
	var err error

	// Top 200 videos with the highest similarity
	row := DB.QueryRow("SELECT related_videos FROM video WHERE id=? limit 1", videoId)
	var jsonData string
	err = row.Scan(&jsonData)
	if err != nil {
		log.Fatal(err)
	}
	// 解析JSON数据
	var data map[string][]uint
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("解析JSON失败：", err)
		return videoList
	}
	var relatedVis []uint
	relatedVis = data["relatedVis"][start: start + requiredVideo]  // get required related videos num
	// fmt.Println(relatedVis)

	query := "SELECT id, user_id, video_uid, type, content, keyword, upload_time, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, top, private, screenshot, hls FROM video WHERE deleted=0 AND private=0 AND hls=1 AND id IN (" + generatePlaceholders(len(relatedVis)) + ")"

	// 构建查询参数
	args := make([]interface{}, len(relatedVis))
	for i, id := range relatedVis {
		args[i] = id
	}

	// 执行查询
	rows, err = DB.Query(query, args...)
	if err != nil {
		fmt.Println("执行查询失败：", err)
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
		currentUserId := uint(0)
		user, _, _ := GetUserInfoById(video.UserId, currentUserId)
		video = common.MakeVideoSupInfo(video, screenshotOk, hlsOk, user)
		videoList = append(videoList, video)
	}
	// fmt.Println(videoList)
	return videoList
}
