package main

import (
	"backend/auth"
	"backend/config"
	"backend/database"
	"backend/tool"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5/request"
	"log"
	"net/http"
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
	err := r.ParseForm()
	if err != nil {
		return
	}
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
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
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	var queryUsername string
	var queryPassword string
	for k, v := range r.Form {
		if k == "username" {
			queryUsername = strings.Join(v, "")
		}
		if k == "password" {
			queryPassword = strings.Join(v, "")
		}
	}
	user, errNo := database.CheckUserPassword(queryUsername, queryPassword)
	if errNo != 0 {
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
					"userid":   user.Id,
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
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	//implement sign up function
	var queryUsername string
	var queryPassword string
	for k, v := range r.Form {
		if k == "username" {
			queryUsername = strings.Join(v, "")
		}
		if k == "password" {
			queryPassword = strings.Join(v, "")
		}
	}
	user, errNo := database.CreateUser(queryUsername, queryPassword)
	if errNo != 0 {
		if errNo == 1 { // user already exists
			errorMsg = "User already exists."
			status = 0
		} else {
			errorMsg = "Unknown error."
			status = 0
		}
	} else {
		user, _ = database.GetUserInfoById(user.Id)
		status = 200
		data = map[string]interface{}{
			"user": map[string]interface{}{
				"userid":   user.Id,
				"username": user.Username,
				"nickname": user.Nickname,
				"reg_time": user.RegTime,
			},
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func findAndCheckToken(r *http.Request) (bool, int, int64, string) {
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
	user, errNo := database.GetUserInfoById(userId)
	if errNo != 0 {
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
				"userid":   user.Id,
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
	// set Nickname by userId
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	var queryNickname string
	for k, v := range r.Form {
		if k == "nickname" {
			queryNickname = strings.Join(v, "")
		}
	}
	errNo := database.SetUserInfo(userId, queryNickname)
	if errNo != 0 {
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
			"nickname": queryNickname,
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoList(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	var queryType int
	var queryUser int
	var queryActionHistory string
	var queryLimit int
	var queryPage int
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	for k, v := range r.Form {
		if k == "type" { // video type
			queryType, _ = strconv.Atoi(strings.Join(v, ""))
		}
		if k == "user" { // if list single user or all
			queryUser, _ = strconv.Atoi(strings.Join(v, ""))
		}
		if k == "action_history" { // if list liked/favorite/history of single user
			queryActionHistory = strings.Join(v, "")
		}
		if k == "limit" {
			queryLimit, _ = strconv.Atoi(strings.Join(v, ""))
		}
		if k == "page" {
			queryPage, _ = strconv.Atoi(strings.Join(v, ""))
		}
	}
	// for some bad parameter
	if queryLimit > 20 {
		queryLimit = 20
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryPage < 0 {
		queryPage = 1
	}

	videoList := database.GetVideoList(queryType, queryUser, queryActionHistory, queryLimit, queryPage)
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
	var queryVideoUid string
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	for k, v := range r.Form {
		if k == "video_uid" {
			queryVideoUid = strings.Join(v, "")
		}
	}
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
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	var queryVideoUid string
	var queryAction string
	for k, v := range r.Form {
		if k == "video_uid" {
			queryVideoUid = strings.Join(v, "")
		}
		if k == "action" {
			queryAction = strings.Join(v, "")
		}
	}
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// do action
	// todo: like, unlike, favorite, unfavorite, comment
	errNo := database.GuestDoVideoAction(video, userId, queryAction)
	if errNo != 0 {
		if errNo == 1 { //already do before
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
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	var queryVideoUid string
	for k, v := range r.Form {
		if k == "video_uid" {
			queryVideoUid = strings.Join(v, "")
		}
	}
	video := database.GetVideoDetail(queryVideoUid)
	if video.Id == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// record this history
	errNo := database.GuestRecordVideoHistory(video, userId)
	if errNo != 0 {
		errorMsg = "Unknown error."
		status = 0
	} else {
		status = 200
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {

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

	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	// check video
	var queryVideoUid string
	var queryVideoTitle string
	var queryVideoContent string
	var queryVideoType int
	for k, v := range r.Form {
		if k == "video_uid" {
			queryVideoUid = strings.Join(v, "")
		}
		if k == "video_title" {
			queryVideoTitle = strings.Join(v, "")
		}
		if k == "video_content" {
			queryVideoContent = strings.Join(v, "")
		}
		if k == "video_type" {
			queryVideoType, _ = strconv.Atoi(strings.Join(v, ""))
		}
	}
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
	errNo := database.SetVideoInfo(queryVideoUid, queryVideoTitle, queryVideoContent, queryVideoType)
	if errNo != 0 {
		errorMsg = "Unknown error."
		status = 0
	} else {
		status = 200
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func deleteVideo(w http.ResponseWriter, r *http.Request) {

}

func setToken(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	var userId int
	for k, v := range r.Form {
		if k == "userid" {
			userId, _ = strconv.Atoi(strings.Join(v, ""))
		}
	}
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
			"token":  token,
			"userid": userId,
			"exp":    tool.UnixTimeToRFC3339(exp),
		}
	}
	SendJSONResponse(w, status, data, errorMsg)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""
	err := r.ParseForm()
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	var token string
	for k, v := range r.Form {
		if k == "token" {
			token = strings.Join(v, "")
		}
	}
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
					"userid": userId,
					"exp":    tool.UnixTimeToRFC3339(exp),
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
	test := true // set to false to disable test api

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

	if test {
		http.HandleFunc("/v1/test/setToken", setToken)
		http.HandleFunc("/v1/test/getToken", getToken)
	}
	err := http.ListenAndServe(":8051", nil)
	if err != nil {
		log.Fatal("Startup err: ", err)
	}

}
