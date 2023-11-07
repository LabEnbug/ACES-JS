package cmd

import (
	"backend/algorithm"
	"backend/config"
	"backend/database/mysql"
	"backend/model"
	"backend/tool"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func GetVideoList(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get} /v1/videos Get video list including normal, user-(uploaded, liked, favorite, watched, followed)
	 * @apiName GetVideoList
	 *
	 * @apiParam {Number} type Video type.
	 * @apiParam {Number} user_id User id.
	 * @apiParam {String} relation Relation.
	 * @apiParam {Number} limit Max number of videos.
	 * @apiParam {Number} start Start at.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	queryParams := r.URL.Query()
	queryType, _ := strconv.Atoi(queryParams.Get("type"))
	queryUserIdTmp, _ := strconv.ParseUint(queryParams.Get("user_id"), 10, 32)
	queryUserId := uint(queryUserIdTmp)
	queryRelation := queryParams.Get("relation")
	queryLimit, _ := strconv.Atoi(queryParams.Get("limit"))
	queryStart, _ := strconv.Atoi(queryParams.Get("start"))

	// for some bad parameter, strict limit
	if queryLimit > 24 {
		queryLimit = 24
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check user
	tokenValid, userId, _, _ := FindAndCheckToken(r)

	// get video list
	var videoList []model.Video
	if queryType == 0 && queryUserId == 0 && queryRelation == "" {
		log.Println(1111111)
		videoList = algorithm.GetRecommendVideoList(queryLimit, queryStart, userId)
	} else {
		videoList = mysql.GetVideoList(queryType, queryUserId, queryRelation, queryLimit, queryStart, userId)
	}
	if len(videoList) == 0 {
		status = 0
		errorMsg = "No more video found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check relation between user and video
	// only need to check when user logged in
	if tokenValid { // user logged in
		for i := 0; i < len(videoList); i++ {
			videoList[i].IsUserLiked, videoList[i].IsUserFavorite, videoList[i].IsUserUploaded,
				videoList[i].IsUserWatched, videoList[i].IsUserLastPlay = mysql.CheckUserVideoAllRelation(userId, videoList[i].Id)
		}
	}

	// if queryUserId is not 0, get video numbers
	if queryUserId != 0 || queryRelation == "watched" {
		var videoNumbers int
		if queryRelation == "watched" {
			// only able to get self watched
			videoNumbers = mysql.GetVideoNum(userId, queryRelation, userId)
		} else {
			videoNumbers = mysql.GetVideoNum(queryUserId, queryRelation, userId)
		}
		data = map[string]interface{}{
			"video_list": videoList,
			"video_num":  videoNumbers,
		}
	} else {
		data = map[string]interface{}{
			"video_list": videoList,
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GetRelatedVideoList(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	queryParams := r.URL.Query()
	queryLimit, _ := strconv.Atoi(queryParams.Get("limit"))
	queryStart, _ := strconv.Atoi(queryParams.Get("start"))

	// for some bad parameter, strict limit
	if queryLimit > 24 {
		queryLimit = 24
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check video
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check user
	tokenValid, userId, _, _ := FindAndCheckToken(r)

	// get related video list
	videoList := mysql.GetRecommendVideoBySimilarity(videoId, queryLimit, queryStart, userId)
	if len(videoList) == 0 {
		status = 0
		errorMsg = "No more video found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check relation between user and video
	// only need to check when user logged in
	if tokenValid { // user logged in
		for i := 0; i < len(videoList); i++ {
			videoList[i].IsUserLiked, videoList[i].IsUserFavorite, videoList[i].IsUserUploaded,
				videoList[i].IsUserWatched, videoList[i].IsUserLastPlay = mysql.CheckUserVideoAllRelation(userId, videoList[i].Id)
		}
	}

	data = map[string]interface{}{
		"video_list": videoList,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GetVideoInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/{videoUid} Get video info
	 * @apiName GetVideoInfo
	 *
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	// check user
	tokenValid, userId, _, _ := FindAndCheckToken(r)

	// get video info
	video := mysql.GetVideoDetailByVideoUid(queryVideoUid, userId)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// private not own by the user will not be visible
	if video.IsPrivate && video.UserId != userId {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check relation between user and video
	// only need to check when user logged in
	if tokenValid { // user logged in
		// check relationship
		video.IsUserLiked, video.IsUserFavorite, video.IsUserUploaded, video.IsUserWatched,
			video.IsUserLastPlay = mysql.CheckUserVideoAllRelation(userId, video.Id)
	}

	if video.IsUserUploaded {
		remainPromoteCount, remainAdvertiseCount := mysql.GetRemainPromoteAndAdvertiseCount(video.Id)
		data = map[string]interface{}{
			"video":                  video,
			"remain_promote_count":   remainPromoteCount,
			"remain_advertise_count": remainAdvertiseCount,
		}
	} else {
		data = map[string]interface{}{
			"video": video,
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func HandleVideoAction(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post|delete} /v1/videos/{videoUid}/{action} Handle video action
	 * @apiName HandleVideoAction
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]
	queryAction := vars["action"]

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)

	needOwnerAction := map[string]bool{
		"top":       true,
		"private":   true,
		"promote":   true,
		"advertise": true,
	}

	needLoginAction := map[string]bool{
		"like":      true,
		"favorite":  true,
		"top":       true,
		"private":   true,
		"promote":   true,
		"advertise": true,
	}

	if needLoginAction[queryAction] && !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	if needOwnerAction[queryAction] && !mysql.CheckUserVideoRelation(userId, videoId, "uploaded") {
		status = 0
		errorMsg = "Not owner."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	switch queryAction {
	case "like", "favorite":
		if r.Method == "POST" {
			// do action
			ok, errNo := mysql.GuestLikeOrFavoriteVideo(videoId, userId, queryAction)
			if !ok {
				if errNo == 1 { // already done
					errorMsg = "Already done."
				} else {
					errorMsg = "Unknown error."
					status = 0
				}
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
		} else if r.Method == "DELETE" {
			// do un-action
			ok, errNo := mysql.GuestLikeOrFavoriteVideo(videoId, userId, "un"+queryAction)
			if !ok {
				if errNo == 1 { // already done
					errorMsg = "Already done."
				} else {
					errorMsg = "Unknown error."
					status = 0
				}
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
		} else {
			status = 0
			errorMsg = "Invalid action."
		}
	case "watch":
		if r.Method == "POST" {
			ok := mysql.GuestWatchVideo(videoId, userId)
			if !ok {
				status = 0
				errorMsg = "Unknown error."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
		} else {
			status = 0
			errorMsg = "Invalid action."
		}
	case "forward":
		if r.Method == "POST" {
			ok := mysql.GuestForwardVideo(videoId, userId)
			if !ok {
				status = 0
				errorMsg = "Unknown error."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
		} else {
			status = 0
			errorMsg = "Invalid action."
		}
	case "top", "private":
		if r.Method == "POST" {
			// top video
			ok := mysql.UserTopOrPrivateVideo(queryVideoUid, queryAction)
			if !ok {
				status = 0
				errorMsg = "Unknown error."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
		} else if r.Method == "DELETE" {
			// do un-action
			// top video
			ok := mysql.UserTopOrPrivateVideo(queryVideoUid, "un"+queryAction)
			if !ok {
				status = 0
				errorMsg = "Unknown error."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
		} else {
			status = 0
			errorMsg = "Invalid action."
		}
	case "promote", "advertise":
		if r.Method == "POST" {
			// parse form
			err := r.ParseMultipartForm(config.MaxNormalPostSize64)
			if err != nil {
				status = 0
				errorMsg = "Failed to parse form."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}
			queryCountTmp := r.FormValue("count")
			queryCount, _ := strconv.Atoi(queryCountTmp)
			if queryCount < 1 || queryCount > 1000 {
				status = 0
				errorMsg = "Invalid count."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}

			// check if user have enough balance
			userBalance := mysql.GetUserBalance(userId)
			queryCountFloat := float64(queryCount)
			if (queryAction == "promote" && userBalance < queryCountFloat*config.PromotePrice) ||
				(queryAction == "advertise" && userBalance < queryCountFloat*config.AdvertisePrice) {
				status = 0
				errorMsg = "Insufficient balance."
				SendJSONResponse(w, status, data, errorMsg)
				return
			}

			if queryAction == "promote" {
				ok := mysql.PromoteVideo(videoId, userId, queryCount)
				if !ok {
					status = 0
					errorMsg = "Unknown error."
					SendJSONResponse(w, status, data, errorMsg)
					return
				}
			} else if queryAction == "advertise" {
				ok := mysql.AdvertiseVideo(videoId, userId, queryCount)
				if !ok {
					status = 0
					errorMsg = "Unknown error."
					SendJSONResponse(w, status, data, errorMsg)
					return
				}
			}
		} else {
			status = 0
			errorMsg = "Invalid action."
		}
	default:
		status = 0
		errorMsg = "Invalid action."
	}

	data = map[string]interface{}{
		"action": queryAction,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/upload Upload video
	 * @apiName UploadVideo
	 *
	 * @apiParam {File} file Video file.
	 * @apiParam {Number} video_type Video type.
	 * @apiParam {String} video_content Video content.
	 * @apiParam {String} video_keyword Video keywords.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// receive file and save to {{BaseLocalFileDir}}/tmp/ for temp
	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// parse form
	// check file size by MaxBytesReader, limit to {{config.MaxUploadVideoSize64}}MB
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadVideoSize64)
	err := r.ParseMultipartForm(config.MaxUploadVideoSize64)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = fmt.Sprintf("File size limit to %dMB.", config.MaxUploadVideoSize)
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoType, _ := strconv.Atoi(r.PostFormValue("video_type"))
	queryVideoContent := r.PostFormValue("video_content")
	queryVideoKeyword := r.PostFormValue("video_keyword")

	// check file title, type, and description
	if queryVideoContent == "" {
		status = 0
		errorMsg = "Video content cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video duplicated
	if mysql.CheckVideoDuplicate(queryVideoContent, queryVideoKeyword, queryVideoType, userId) {
		status = 0
		errorMsg = "Video already exists."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check video type
	if queryVideoType < 1 || !mysql.CheckVideoType(queryVideoType) {
		status = 0
		errorMsg = "Invalid video type."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = "Failed to get file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	defer file.Close()

	// check file type (accept all video types)
	// check from filename, lower filename must include (.mp4 or .mov or .avi or .wmv or .flv or .mkv or .webm
	// or .ts or .rm or .rmvb or .3gp or .mpeg or .mpg or .m4v or .f4v)
	fileName := strings.ToLower(fileHeader.Filename)
	extensionList := []string{".mp4", ".mov", ".avi", ".wmv", ".flv", ".mkv", ".webm", ".ts", ".rm",
		".rmvb", ".3gp", ".mpeg", ".mpg", ".m4v", ".f4v"}
	// check if filename contains any extension name
	extensionNameExist := false
	for _, extensionName := range extensionList {
		if strings.HasSuffix(fileName, extensionName) {
			extensionNameExist = true
			break
		}
	}
	if !extensionNameExist {
		status = 0
		errorMsg = "File type not supported."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// generate video uid
	videoUid := uuid.New().String()

	// save file to {{BaseLocalFileDir}}/tmp/{{video_uid}}
	fileSavePath := path.Join(config.BaseLocalFileDir, "tmp", videoUid)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = "Failed to receive file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	err = os.WriteFile(fileSavePath, fileBytes, 0644)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = "Failed to save file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// save video info to database
	video, ok := mysql.CreateVideo(userId, videoUid, queryVideoType, queryVideoContent, queryVideoKeyword)
	if !ok {
		status = 0
		errorMsg = "Failed to save video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// process video, 2 choices:
	// 1. upload file directly to qiniu and wait for auto hls transcode callback,
	// 2. transcode to hls by ffmpeg and upload to qiniu

	// make choice
	if config.VideoProcessType == 1 { // choice 1
		subDir := "origin/remote_process"

		// upload file to qiniu
		// try to upload 3 times
		uploadOk := false
		for i := 0; i < 3; i++ {
			ok = tool.UploadFileToQiniu(fileSavePath, path.Join(config.BaseRemoteFileDir, subDir, videoUid))
			if ok {
				uploadOk = true
				break
			}
		}

		if !uploadOk {
			status = 0
			errorMsg = "Failed to upload file to cloud."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		// just return, wait for callback asynchronously
		data = map[string]interface{}{
			"video": video,
		}
	} else { // choice 2
		// not implement now
		//subDir := "origin/local_process" // if transcode locally
	}

	// delete tmp file
	_ = os.Remove(fileSavePath)

	SendJSONResponse(w, status, data, errorMsg)
}

func UploadVideoRemote(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/uploadRemote Upload video remote (create first, upload later) (admin only)
	 * @apiName UploadVideoRemote
	 *
	 * @apiParam {Number} video_type Video type.
	 * @apiParam {String} video_content Video content.
	 * @apiParam {String} video_keyword Video keywords.
	 * @apiParam {Number} user_id User Id.
	 * @apiParam {String} key Secret Key.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// receive file and save to {{BaseLocalFileDir}}/tmp/ for temp
	// check token
	//tokenValid, _, _, _ := FindAndCheckToken(r)
	//if !tokenValid {
	//	status := 0
	//	data := map[string]interface{}{}
	//	errorMsg := "Not logged in."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

	// parse form
	err := r.ParseMultipartForm(config.MaxUploadVideoSize64)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = fmt.Sprintf("File size limit to %dMB.", config.MaxUploadVideoSize)
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoType, _ := strconv.Atoi(r.PostFormValue("video_type"))
	queryVideoContent := r.PostFormValue("video_content")
	queryVideoKeyword := r.PostFormValue("video_keyword")
	queryUserIdTmp, _ := strconv.Atoi(r.PostFormValue("user_id"))
	queryUserId := uint(queryUserIdTmp)
	queryKey := r.PostFormValue("key")
	if queryKey != "acesaces" {
		status = 0
		errorMsg = "Invalid key."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check file title, type, and description
	if queryVideoContent == "" {
		status = 0
		errorMsg = "Video content cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video duplicated
	if mysql.CheckVideoDuplicate(queryVideoContent, queryVideoKeyword, queryVideoType, queryUserId) {
		status = 0
		errorMsg = "Video already exists."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check video type
	if queryVideoType < 1 || !mysql.CheckVideoType(queryVideoType) {
		status = 0
		errorMsg = "Invalid video type."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// generate video uid
	videoUid := uuid.New().String()

	// save video info to database
	video, ok := mysql.CreateVideo(queryUserId, videoUid, queryVideoType, queryVideoContent, queryVideoKeyword)
	if !ok {
		status = 0
		errorMsg = "Failed to save video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"video": video,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func UploadVideoFile(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/upload Upload video file
	 * @apiName UploadVideoFile
	 *
	 * @apiParam {File} file Video file.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// receive file and save to {{BaseLocalFileDir}}/tmp/ for temp
	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// parse form
	// check file size by MaxBytesReader, limit to {{config.MaxUploadVideoSize64}}MB
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadVideoSize64)
	err := r.ParseMultipartForm(config.MaxUploadVideoSize64)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = fmt.Sprintf("File size limit to %dMB.", config.MaxUploadVideoSize)
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = "Failed to get file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	defer file.Close()

	// check file type (accept all video types)
	// check from filename, lower filename must include (.mp4 or .mov or .avi or .wmv or .flv or .mkv or .webm
	// or .ts or .rm or .rmvb or .3gp or .mpeg or .mpg or .m4v or .f4v)
	fileName := strings.ToLower(fileHeader.Filename)
	extensionList := []string{".mp4", ".mov", ".avi", ".wmv", ".flv", ".mkv", ".webm", ".ts", ".rm",
		".rmvb", ".3gp", ".mpeg", ".mpg", ".m4v", ".f4v"}
	// check if filename contains any extension name
	extensionNameExist := false
	for _, extensionName := range extensionList {
		if strings.HasSuffix(fileName, extensionName) {
			extensionNameExist = true
			break
		}
	}
	if !extensionNameExist {
		status = 0
		errorMsg = "File type not supported."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// generate video uid
	videoUid := uuid.New().String()

	// save file to {{BaseLocalFileDir}}/tmp/{{video_uid}}
	fileSavePath := path.Join(config.BaseLocalFileDir, "tmp", videoUid)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = "Failed to receive file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	err = os.WriteFile(fileSavePath, fileBytes, 0644)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = "Failed to save file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// save video info to database
	ok := mysql.CreateVideoWithoutInfo(userId, videoUid)
	if !ok {
		status = 0
		errorMsg = "Failed to save video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"video_uid": videoUid,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func ConfirmVideoUpload(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {put} /v1/video/upload/{videoUid} Confirm video upload
	 * @apiName ConfirmVideoUpload
	 *
	 * @apiParam {string} video_uid Video uid.
	 * @apiParam {Number} video_type Video type.
	 * @apiParam {String} video_content Video content.
	 * @apiParam {String} video_keyword Video keywords.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoType, _ := strconv.Atoi(r.FormValue("video_type"))
	queryVideoContent := r.FormValue("video_content")
	queryVideoKeyword := r.FormValue("video_keyword")

	if len(queryVideoContent) > 120 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Content too long."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	if len(queryVideoKeyword) > 100 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Keyword too long."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if user has created this video
	if !mysql.CheckUserVideoRelation(userId, mysql.GetVideoIdByVideoUid(queryVideoUid), "uploaded") {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// update video info
	ok := mysql.ConfirmCreateVideoWithInfo(queryVideoUid, queryVideoType, queryVideoContent, queryVideoKeyword)
	if !ok {
		status = 0
		errorMsg = "Failed to save video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get video info
	video := mysql.GetVideoDetailByVideoUid(queryVideoUid, userId)

	// upload file to qiniu
	fileSavePath := path.Join(config.BaseLocalFileDir, "tmp", queryVideoUid)

	// process video, 2 choices:
	// 1. upload file directly to qiniu and wait for auto hls transcode callback,
	// 2. transcode to hls by ffmpeg and upload to qiniu

	// make choice
	if config.VideoProcessType == 1 { // choice 1
		subDir := "origin/remote_process"

		// upload file to qiniu
		// try to upload 3 times
		uploadOk := false
		for i := 0; i < 3; i++ {
			ok := tool.UploadFileToQiniu(fileSavePath, path.Join(config.BaseRemoteFileDir, subDir, queryVideoUid))
			if ok {
				uploadOk = true
				break
			}
		}

		if !uploadOk {
			status = 0
			errorMsg = "Failed to upload file to cloud."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		// just return, wait for callback asynchronously
		data = map[string]interface{}{
			"video": video,
		}
	} else { // choice 2
		// not implement now
		//subDir := "origin/local_process" // if transcode locally
	}

	// delete tmp file
	_ = os.Remove(fileSavePath)

	SendJSONResponse(w, status, data, errorMsg)
}

func SetVideoInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {put} /v1/video/{videoUid} Set video info
	 * @apiName SetVideoInfo
	 *
	 * @apiParam {String} video_content Video content.
	 * @apiParam {String} video_keyword Video keywords.
	 * @apiParam {Number} video_type Video type.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoContent := r.FormValue("video_content")
	queryVideoKeyword := r.FormValue("video_keyword")
	queryVideoTypeTmp, _ := strconv.Atoi(r.FormValue("video_type"))
	queryVideoType := int8(queryVideoTypeTmp)

	if len(queryVideoContent) > 120 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Content too long."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	if len(queryVideoKeyword) > 100 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Keyword too long."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	if !mysql.CheckUserVideoRelation(userId, videoId, "uploaded") {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// edit video info
	ok := mysql.SetVideoInfo(queryVideoUid, queryVideoContent, queryVideoKeyword, queryVideoType)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get video info
	video := mysql.GetVideoDetailByVideoUid(queryVideoUid, userId)

	data = map[string]interface{}{
		"video": video,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {delete} /v1/video/{videoUid} Delete video
	 * @apiName DeleteVideo
	 *
	 * @apiParam {String} video_uid Video uid.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	// check video
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	if !mysql.CheckUserVideoRelation(userId, videoId, "uploaded") {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// delete video
	ok := mysql.DeleteVideo(queryVideoUid)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GetVideoTypes(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get} /v1/video/types Get video types
	 * @apiName GetVideoTypes
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// get video types
	videoTypes := mysql.GetVideoTypes()
	if len(videoTypes) == 0 {
		status = 0
		errorMsg = "No video type found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"video_types": videoTypes,
	}

	SendJSONResponse(w, status, data, errorMsg)
}
