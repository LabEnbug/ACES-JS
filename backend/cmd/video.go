package cmd

import (
	"backend/config"
	"backend/database/mysql"
	"backend/tool"
	"backend/algorithm"
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
	videoList := mysql.GetVideoList(queryType, queryUserId, queryRelation, queryLimit, queryStart, userId)
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

func GetVideoInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/{videoUid} Get video info
	 * @apiName GetVideoInfo
	 *
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	//// check method, only accept POST
	//if r.Method != "POST" {
	//	status = 0
	//	errorMsg = "Invalid request method."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

	//// parse form
	//err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	//if err != nil {
	//	status = 0
	//	errorMsg = "Failed to parse form."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}
	//queryVideoUid := r.PostFormValue("video_uid")
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

func DoVideoAction(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/action Do video action
	 * @apiName DoVideoAction
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {String} action Action.
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

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
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
	queryVideoUid := r.PostFormValue("video_uid")
	queryAction := r.PostFormValue("action")

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// do action: (like, unlike), (favorite, unfavorite), (top, untop)
	ok, errNo := mysql.DoVideoAction(videoId, userId, queryAction)
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

	data = map[string]interface{}{
		"action": queryAction,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func RecordWatchedVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/watch Record watched video
	 * @apiName RecordWatchedVideo
	 *
	 * @apiParam {String} video_uid Video uid.
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

	// check token
	_, userId, _, _ := FindAndCheckToken(r)
	// not logged-in user should also be able to watch video, just record as userId=0

	//if !tokenValid {
	//	status = 0
	//	errorMsg = "Not logged in."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// record video history
	ok := mysql.GuestWatchedVideo(videoId, userId)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GuestForwardVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/forward Forward video
	 * @apiName GuestForwardVideo
	 *
	 * @apiParam {String} video_uid Video uid.
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

	// check token
	_, userId, _, _ := FindAndCheckToken(r)
	// not logged-in user should also be able to watch video, just record as userId=0

	//if !tokenValid {
	//	status = 0
	//	errorMsg = "Not logged in."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// record video history
	ok := mysql.GuestForwardVideo(videoId, userId)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
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
	 * @api {post} /v1/video/upload/file Upload video file
	 * @apiName UploadVideoFile
	 *
	 * @apiParam {File} file Video file.
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
	 * @api {post} /v1/video/upload/confirm Confirm video upload
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

	// check method, only accept POST
	if r.Method != "POST" {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

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
	queryVideoUid := r.PostFormValue("video_uid")
	queryVideoType, _ := strconv.Atoi(r.PostFormValue("video_type"))
	queryVideoContent := r.PostFormValue("video_content")
	queryVideoKeyword := r.PostFormValue("video_keyword")

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

	//// check method, only accept POST
	//if r.Method != "POST" {
	//	status = 0
	//	errorMsg = "Invalid request method."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

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

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")

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

func TopVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/top Top video
	 * @apiName TopVideo
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {String} type Top type (top || untop).
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

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
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
	queryVideoUid := r.PostFormValue("video_uid")
	queryType := r.PostFormValue("type")

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
	ok := mysql.TopVideo(queryVideoUid, queryType)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func PrivateVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/private Make video private
	 * @apiName PrivateVideo
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {String} type Private type (private || unprivate).
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

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
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
	queryVideoUid := r.PostFormValue("video_uid")
	queryType := r.PostFormValue("type")

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
	ok := mysql.PrivateVideo(queryVideoUid, queryType)
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
	 * @api {get|post} /v1/video/types Get video types
	 * @apiName GetVideoTypes
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method
	if r.Method != "GET" && r.Method != "POST" {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

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

func GetRecommendVideoList(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/userRecomList Get recommend video list by userId
	 *  
	 * @apiName GetRecommendVideoList
	 *
	 * @apiParam {Number} user_id User id. uint
	 * @apiParam {Number} limit the number of video needed. int
	 * @apiParam {Number} start Start at. int
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

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryUserIdTmp, _ := strconv.ParseUint(r.PostFormValue("user_id"), 10, 32)
	queryUserId := uint(queryUserIdTmp)
	queryLimit, _ := strconv.Atoi(r.PostFormValue("limit"))
	queryStart, _ := strconv.Atoi(r.PostFormValue("start"))

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
	videoList := algorithm.GetRecommendVideoList(queryUserId, queryLimit, queryStart)
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

func PromoteVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/promote Promote video
	 * @apiName PromoteVideo
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {Number} count Bandwidth count.
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
	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
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
	queryVideoUid := r.PostFormValue("video_uid")
	queryCountTmp := r.PostFormValue("count")
	queryCount, _ := strconv.Atoi(queryCountTmp)
	if queryCount < 1 || queryCount > 1000 {
		status = 0
		errorMsg = "Invalid count."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check video
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	// todo: promotion may be able to be created by not owner, like "help to promote"
	if !mysql.CheckUserVideoRelation(userId, videoId, "uploaded") {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if user have enough balance
	userBalance := mysql.GetUserBalance(userId)
	queryCountFloat := float64(queryCount)
	if userBalance < queryCountFloat*config.PromotePrice {
		status = 0
		errorMsg = "Insufficient balance."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// promote
	ok := mysql.PromoteVideo(videoId, userId, queryCount)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func AdvertiseVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/advertise Advertise video
	 * @apiName AdvertiseVideo
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {Number} count Bandwidth count.
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

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
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
	queryVideoUid := r.PostFormValue("video_uid")
	queryCountTmp := r.PostFormValue("count")
	queryCount, _ := strconv.Atoi(queryCountTmp)
	if queryCount < 1 || queryCount > 1000 {
		status = 0
		errorMsg = "Invalid count."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

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

	// check if user have enough balance
	userBalance := mysql.GetUserBalance(userId)
	queryCountFloat := float64(queryCount)
	if userBalance < queryCountFloat*config.AdvertisePrice {
		status = 0
		errorMsg = "Insufficient balance."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// advertise
	ok := mysql.AdvertiseVideo(videoId, userId, queryCount)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}
