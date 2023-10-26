package main

import (
	"backend/auth"
	"backend/config"
	"backend/database"
	"backend/model"
	"backend/tool"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type Response struct {
	Status   int         `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	ErrorMsg string      `json:"err_msg,omitempty"`
}

func SendJSONResponse(w http.ResponseWriter, status int, data interface{}, errorMsg string) {
	response := Response{
		Status:   status,
		Data:     data,
		ErrorMsg: errorMsg,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		return
	}
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		return
	}
	//log.Println(r.Form)
	//log.Println("path", r.URL.Path)
	//log.Println("scheme", r.URL.Scheme)
	//log.Println(r.Form["url_long"])
	for k, v := range r.Form {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), "key:", k, ", val:", strings.Join(v, ""))
	}
	_, err = fmt.Fprintf(w, "Hello")
	if err != nil {
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryUsername := r.FormValue("username")
	queryPassword := r.FormValue("password")
	user, ok, errNo := database.CheckUserPassword(queryUsername, queryPassword)
	if !ok {
		if errNo == 1 { // user not found
			errorMsg = "User not found."
			status = 0
		} else if errNo == 2 {
			errorMsg = "Wrong password."
			status = 0
		} else {
			errorMsg = "Unknown error."
			status = 0
		}
	} else {
		// user match
		token, exp, err := auth.CreateToken(user.Id)
		if err != nil {
			status = 0
			errorMsg = "Failed to create token."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		err = database.StoreToken(token)
		if err != nil {
			status = 0
			errorMsg = "Failed to store token."
		} else {
			database.SetLastLogin(user.Id)
			status = 200
			data = map[string]interface{}{
				"token": token,
				"exp":   tool.UnixTimeToRFC3339(exp),
				"user": map[string]interface{}{
					"user_id":  user.Id,
					"username": user.Username,
					"nickname": user.Nickname,
					"reg_time": user.RegTime,
				},
			}
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func logout(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// find token
	token, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	//check token exist
	isExist, err := database.CheckTokenIsExist(token)
	if err != nil {
		status = 0
		errorMsg = "Failed to check token."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	if !isExist {
		status = 200
	} else {
		// test token will never be revoke
		if token != "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg0Njk0MzUsInVzZXJJZCI6MX0.oWpg0PRAgTWApqZ8AObPPlPrmbgBSTWF2zcMWJxoxcxASsyLLOM-ehfpg_A6D7ksRCSG_UQ0mtTp4ongY-gtwA" {
			err = database.RevokeToken(token)
			if err != nil {
				status = 0
				errorMsg = "Failed to revoke token."
			}
		}

	}
	SendJSONResponse(w, status, data, errorMsg)
}

func signup(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	//implement sign up function
	queryUsername := r.FormValue("username")
	queryPassword := r.FormValue("password")
	queryNickname := r.FormValue("nickname")
	user, ok, errNo := database.CreateUser(queryUsername, queryPassword, queryNickname)
	if !ok {
		if errNo == 1 { // user already exists
			errorMsg = "User already exists."
			status = 0
		} else {
			errorMsg = "Unknown error."
			status = 0
		}
	} else {
		user, _, _ = database.GetUserInfoById(user.Id)
		status = 200
		data = map[string]interface{}{
			"user": map[string]interface{}{
				"user_id":  user.Id,
				"username": user.Username,
				"nickname": user.Nickname,
				"reg_time": user.RegTime,
			},
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func findAndCheckToken(r *http.Request) (bool, uint, int64, string) {
	// find token
	// Authorization: Bearer xxx
	token, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return false, 0, 0, ""
	}

	// check token
	isExist, err := database.CheckTokenIsExist(token)
	if err != nil {
		return false, 0, 0, ""
	} else {
		if isExist {
			userId, exp, err := auth.GetInfoFromToken(token)
			if err != nil {
				return false, 0, 0, ""
			} else {
				return true, userId, exp, token
			}
		} else {
			return false, 0, 0, ""
		}
	}
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	tokenValid, userId, exp, token := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get user info
	user, ok, errNo := database.GetUserInfoById(userId)
	if !ok {
		if errNo == 1 { // user not found
			errorMsg = "User not found."
			status = 0
		} else {
			errorMsg = "Unknown error."
			status = 0
		}
	} else {
		status = 200
		data = map[string]interface{}{
			"token": token,
			"exp":   tool.UnixTimeToRFC3339(exp),
			"user": map[string]interface{}{
				"user_id":  user.Id,
				"username": user.Username,
				"nickname": user.Nickname,
				"reg_time": tool.DatabaseTimeToRFC3339(user.RegTime),
			},
		}
	}
	SendJSONResponse(w, status, data, errorMsg)
}

func setUserInfo(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryNickname := r.FormValue("nickname")
	// set Nickname by userId
	ok := database.SetUserInfo(userId, queryNickname)
	if !ok {
		errorMsg = "Unknown error."
		status = 0
	} else {
		user, _, _ := database.GetUserInfoById(userId)
		status = 200
		data = map[string]interface{}{
			"user": map[string]interface{}{
				"user_id":  user.Id,
				"username": user.Username,
				"nickname": user.Nickname,
				"reg_time": tool.DatabaseTimeToRFC3339(user.RegTime),
			},
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoList(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryType, _ := strconv.Atoi(r.FormValue("type"))
	parseUint, _ := strconv.ParseUint(r.FormValue("user_id"), 10, 32)
	queryUserId := uint(parseUint)
	queryActionHistory := r.FormValue("action_history")
	queryLimit, _ := strconv.Atoi(r.FormValue("limit"))
	queryPage, _ := strconv.Atoi(r.FormValue("page"))

	// for some bad parameter
	if queryLimit > 20 {
		queryLimit = 20
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryPage < 0 {
		queryPage = 1
	}

	videoList := database.GetVideoList(queryType, queryUserId, queryActionHistory, queryLimit, queryPage)
	if len(videoList) == 0 {
		status = 0
		errorMsg = "No more video found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	status = 200
	data = map[string]interface{}{
		"video_list": videoList,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoDetail(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.FormValue("video_uid")
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	status = 200
	data = map[string]interface{}{
		"video": video,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func doVideoAction(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	queryVideoUid := r.FormValue("video_uid")
	queryAction := r.FormValue("action")
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// do action
	// todo: like, unlike, favorite, unfavorite, comment
	ok, errNo := database.GuestDoVideoAction(video, userId, queryAction)
	if !ok {
		if errNo == 1 { // already done
			errorMsg = "Already done."
			status = 0
		} else {
			errorMsg = "Unknown error."
			status = 0
		}
	} else {
		status = 200
		data = map[string]interface{}{
			"action": queryAction,
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func recordVideoHistory(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	queryVideoUid := r.FormValue("video_uid")
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// record this history
	ok := database.GuestRecordVideoHistory(video, userId)
	if !ok {
		errorMsg = "Unknown error."
		status = 0
	} else {
		status = 200
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// only accept POST
	if r.Method != "POST" {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// receive file and save to {{BaseLocalFileDir}}/tmp/
	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check file size by MaxBytesReader, max 50MB
	MaxUploadSize := int64(50 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		status = 0
		errorMsg = fmt.Sprintf("File size limit to %dMB.", int(math.Floor(float64(MaxUploadSize)/1024/1024)))
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		status = 0
		errorMsg = "Failed to get file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	defer file.Close()
	// check file type (accept all video types)
	fileType := fileHeader.Header.Get("Content-Type")
	if !strings.Contains(fileType, "video") {
		status = 0
		errorMsg = "File type not supported."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// get file title, type, and description
	queryVideoType, _ := strconv.Atoi(r.FormValue("video_type"))
	queryVideoTitle := r.FormValue("video_title")
	queryVideoContent := r.FormValue("video_content")
	if queryVideoTitle == "" || queryVideoContent == "" {
		status = 0
		errorMsg = "Video title or content cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// generate video uid
	videoUid := uuid.New().String()
	// save file to {{BaseLocalFileDir}}/tmp/{{video_uid}}
	fileSavePath := path.Join(config.BaseLocalFileDir, "tmp", videoUid)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		status = 0
		errorMsg = "Failed to receive file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	err = os.WriteFile(fileSavePath, fileBytes, 0644)
	if err != nil {
		status = 0
		errorMsg = "Failed to save file."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// save video info to database
	video, ok := database.UserCreateVideo(userId, videoUid, queryVideoType, queryVideoTitle, queryVideoContent)
	if !ok {
		status = 0
		errorMsg = "Failed to save video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// 2 choices: 1. upload file to qiniu and wait for auto hls callback, 2. upload file to qiniu and transcode to hls by ffmpeg
	// upload file to qiniu, different dir by config.VideoProcessType
	subDir := "origin/local_process"  // if transcode by local
	if config.VideoProcessType == 1 { // if transcode by qiniu
		subDir = "origin/remote_process"
	}
	ok = tool.UploadFileToQiniu(fileSavePath, path.Join(config.BaseRemoteFileDir, subDir, videoUid))
	if !ok {
		status = 0
		errorMsg = "Failed to upload file to cloud."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// make choice
	if config.VideoProcessType == 1 { // wait for auto hls callback, return now
		data = map[string]interface{}{
			"video": video,
		}
	} else if config.VideoProcessType == 2 { // transcode to hls by ffmpeg
		// not implement now
	}
	// delete file
	_ = os.Remove(fileSavePath)

	SendJSONResponse(w, status, data, errorMsg)
}

func setVideoInfo(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	queryVideoUid := r.FormValue("video_uid")
	queryVideoTitle := r.FormValue("video_title")
	queryVideoContent := r.FormValue("video_content")
	queryVideoType, _ := strconv.Atoi(r.FormValue("video_uid"))
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	if video.UserId != userId {
		status = 0
		errorMsg = "This video isn't uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// edit video info
	ok := database.SetVideoInfo(queryVideoUid, queryVideoTitle, queryVideoContent, queryVideoType)
	if !ok {
		errorMsg = "Unknown error."
		status = 0
	}

	// get video info
	video = database.GetVideoDetail(queryVideoUid)
	status = 200
	data = map[string]interface{}{
		"video": video,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func deleteVideo(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	queryVideoUid := r.FormValue("video_uid")
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	if video.UserId != userId {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// delete video
	ok := database.DeleteVideo(queryVideoUid)
	if !ok {
		errorMsg = "Unknown error."
		status = 0
	} else {
		status = 200
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func qiniuHlsCallback(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check method
	if r.Method == "POST" {
		postData, err := io.ReadAll(r.Body)
		if err != nil {
			status = 0
			errorMsg = "Failed to read post data."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
		// {"version":"v3","id":"cn-east-2.01c201cwib1vvu4ido00mvdadi000avw","reqid":"CAAAAJQnVO8copEX","pipeline":"1382943942.default.sys","input":{"kodo_file":{"bucket":"aces-js","key":"origin/remote_process/95e255ae-72f7-11ee-9ed6-0242ac1a0003-123.mp4"}},"code":0,"desc":"successfully completed","ops":[{"id":"node2_avthumb","fop":{"cmd":"avthumb/m3u8/vcodec/libx264/acodec/libfdk_aac/ar/48000/ab/138k/crf/20/pattern/{{_base64_join \"video/\" .meta.fname \"/\" .meta.count \".ts\"}}","input_from":"__origin__","result":{"code":0,"desc":"successfully completed","has_output":false}}},{"id":"node4_saveas","fop":{"cmd":"saveas/YWNlcy1qczp2aWRlby85NWUyNTVhZS03MmY3LTExZWUtOWVkNi0wMjQyYWMxYTAwMDMtMTIzL2luZGV4Lm0zdTg=","result":{"code":0,"desc":"successfully completed","has_output":true,"kodo_file":{"bucket":"aces-js","key":"video/95e255ae-72f7-11ee-9ed6-0242ac1a0003-123/index.m3u8","hash":"FmAwuxAOQaKTE8Xv3f0YYoZZvXFk"}}},"depends":["node2_avthumb"]}],"created_at":1698316775387}

		// parse json
		var qiniuCallbackData model.QiniuCallbackData
		err = json.Unmarshal(postData, &qiniuCallbackData)
		if err != nil {
			status = 0
			errorMsg = "Failed to parse json."
			SendJSONResponse(w, status, data, errorMsg)
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "failed to parse json:", string(postData))
			return
		}
		file := qiniuCallbackData.Input.KodoFile.Key
		// get video uid
		videoUid := strings.Split(file, "/")[2][:36]
		// get video info
		video := database.GetVideoDetail(videoUid)
		if video.Id == 0 {
			status = 0
			errorMsg = "Video not found."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
		// update video info
		ok := database.CallbackUpdateVideoHLS(videoUid)
		if !ok {
			status = 0
			errorMsg = "Failed to update video info."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
	} else {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func qiniuScreenshotCallback(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check method
	if r.Method == "POST" {
		postData, err := io.ReadAll(r.Body)
		if err != nil {
			status = 0
			errorMsg = "Failed to read post data."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		// parse json
		var qiniuCallbackData model.QiniuCallbackData
		err = json.Unmarshal(postData, &qiniuCallbackData)
		if err != nil {
			status = 0
			errorMsg = "Failed to parse json."
			SendJSONResponse(w, status, data, errorMsg)
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "failed to parse json:", string(postData))
			return
		}
		file := qiniuCallbackData.Input.KodoFile.Key
		// get video uid
		videoUid := strings.Split(file, "/")[2]
		// get video info
		video := database.GetVideoDetail(videoUid)
		if video.Id == 0 {
			status = 0
			errorMsg = "Video not found."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
		// update video info
		ok := database.CallbackUpdateVideoScreenshot(videoUid)
		if !ok {
			status = 0
			errorMsg = "Failed to update video info."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
	} else {
		status = 0
		errorMsg = "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoTypes(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	videoTypes := database.GetVideoTypes()
	if len(videoTypes) == 0 {
		status = 0
		errorMsg = "No video type found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	status = 200
	data = map[string]interface{}{
		"video_types": videoTypes,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func setToken(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	userIdTmp, _ := strconv.ParseUint(r.FormValue("user_id"), 10, 32)
	userId := uint(userIdTmp)
	token, exp, err := auth.CreateToken(userId)
	if err != nil {
		status = 0
		errorMsg = "Failed to create token."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	err = database.StoreToken(token)
	if err != nil {
		status = 0
		errorMsg = "Failed to store token."
	} else {
		status = 200
		data = map[string]interface{}{
			"token":   token,
			"user_id": userId,
			"exp":     tool.UnixTimeToRFC3339(exp),
		}
	}
	SendJSONResponse(w, status, data, errorMsg)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	// parse form
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	token := r.FormValue("token")
	isExist, err := database.CheckTokenIsExist(token)
	if err != nil {
		status = 0
		errorMsg = "Failed to check token."
	} else {
		if isExist {
			userId, exp, err := auth.GetInfoFromToken(token)
			if err != nil {
				status = 0
				errorMsg = "Failed to parse token."
			} else {
				status = 200
				data = map[string]interface{}{
					"user_id": userId,
					"exp":     tool.UnixTimeToRFC3339(exp),
				}
			}
		} else {
			status = 0
			errorMsg = "Token not exist."
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func main() {
	config.Test = true // set to false to disable test api
	config.ShowLog = true
	config.VideoProcessType = 1 // 1: wait for qiniu transcode callback; 2: transcode by ffmpeg locally

	config.InitConfig()

	database.InitDBPool()
	defer database.CloseDBPool()
	database.InitRedisPool()
	defer database.CloseRedisPool()
	auth.InitEcdsa()

	http.HandleFunc("/", defaultPage)
	http.HandleFunc("/v1/login", login)
	http.HandleFunc("/v1/logout", logout)
	http.HandleFunc("/v1/signup", signup)
	http.HandleFunc("/v1/user/info", getUserInfo)
	http.HandleFunc("/v1/user/edit", setUserInfo)
	http.HandleFunc("/v1/video/list", getVideoList)
	http.HandleFunc("/v1/video/detail", getVideoDetail)
	http.HandleFunc("/v1/video/do", doVideoAction)
	http.HandleFunc("/v1/video/record", recordVideoHistory)
	http.HandleFunc("/v1/video/upload", uploadVideo)
	http.HandleFunc("/v1/video/edit", setVideoInfo)
	http.HandleFunc("/v1/video/delete", deleteVideo)

	http.HandleFunc("/v1/video/types", getVideoTypes)

	http.HandleFunc("/callback/qiniu/hls", qiniuHlsCallback)
	http.HandleFunc("/callback/qiniu/screenshot", qiniuScreenshotCallback)

	if config.Test {
		http.HandleFunc("/v1/test/setToken", setToken)
		http.HandleFunc("/v1/test/getToken", getToken)
	}
	err := http.ListenAndServe(":8051", nil)
	if err != nil {
		log.Fatal("Startup err: ", err)
	}

}
