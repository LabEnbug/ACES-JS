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
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return
	}
}

func defaultPage(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return
	}
	if config.ShowLog {
		for k, v := range r.Form {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "key:", k, ", val:", strings.Join(v, ""))
		}
		for k, v := range r.PostForm {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "key:", k, ", val:", strings.Join(v, ""))
		}
	}
	_, err = fmt.Fprintf(w, "Hello")
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		return
	}
}

func findAndCheckToken(r *http.Request) (bool, uint, int64, string) {
	// find token
	// Authorization: Bearer xxx
	token, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return false, 0, 0, ""
	}

	// check token
	isExist, _ := database.CheckTokenIsExist(token)
	if isExist {
		userId, exp, err := auth.GetInfoFromToken(token)
		if err == nil {
			return true, userId, exp, token
		}
	}
	return false, 0, 0, ""
}

func login(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/user/login User login
	 * @apiName UserLogin
	 *
	 * @apiParam {String} username Username.
	 * @apiParam {String} password Password (origin).
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
	queryUsername := r.PostFormValue("username")
	queryPassword := r.PostFormValue("password")

	// check user
	user, ok, errNo := database.CheckUserPassword(queryUsername, queryPassword)
	if !ok {
		if errNo == 1 { // user not found
			status = 0
			errorMsg = "User not found."
		} else if errNo == 2 { // wrong password
			status = 0
			errorMsg = "Wrong password."
		} else {
			status = 0
			errorMsg = "Unknown error."
		}
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// user match, create token and store it
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
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// update last login time
	database.SetLastLogin(user.Id)
	status = 200
	data = map[string]interface{}{
		"token": token,
		"exp":   tool.UnixTimeToRFC3339(exp),
		"user":  user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func logout(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get|post} /v1/user/logout User logout
	 * @apiName UserLogout
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

	// check user
	tokenValid, _, _, token := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// test token of user1 will never be revoked
	if !config.Test && token != "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg0Njk0MzUsInVzZXJJZCI6MX0.oWpg0PRAgTWApqZ8AObPPlPrmbgBSTWF2zcMWJxoxcxASsyLLOM-ehfpg_A6D7ksRCSG_UQ0mtTp4ongY-gtwA" {
		err := database.RevokeToken(token)
		if err != nil {
			status = 0
			errorMsg = "Failed to revoke token."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func signup(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/user/signup User sign up
	 * @apiName UserSignUp
	 *
	 * @apiParam {String} username Username.
	 * @apiParam {String} password Password (origin).
	 * @apiParam {String} nickname Nickname.
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
	queryUsername := r.PostFormValue("username")
	queryPassword := r.PostFormValue("password")
	queryNickname := r.PostFormValue("nickname")

	// create user
	user, ok, errNo := database.CreateUser(queryUsername, queryPassword, queryNickname)
	if !ok {
		if errNo == 1 { // user already exists
			status = 0
			errorMsg = "User already exists."
		} else {
			status = 0
			errorMsg = "Unknown error."
		}
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get user info
	user, _, _ = database.GetUserInfoById(user.Id, 0)
	data = map[string]interface{}{
		"user": user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get|post} /v1/user/info Get user info
	 * @apiName GetUserInfo
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

	// check token
	tokenValid, userId, exp, token := findAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get user info
	user, ok, errNo := database.GetUserInfoById(userId, 0)
	if !ok {
		if errNo == 1 { // user not found
			status = 0
			errorMsg = "User not found."
		} else {
			status = 0
			errorMsg = "Unknown error."
		}
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"token": token,
		"exp":   tool.UnixTimeToRFC3339(exp),
		"user":  user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getOtherUserInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/user/query Get other user info
	 * @apiName GetOtherUserInfo
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
	queryUsername := r.PostFormValue("username")
	if queryUsername == "" {
		status = 0
		errorMsg = "Username cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check user
	_, userId, _, _ := findAndCheckToken(r)

	// get user info
	user, userExist := database.GetUserInfoByUsername(queryUsername, userId)
	if !userExist {
		status = 0
		errorMsg = "User not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"user": user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func followUser(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/user/follow Follow user
	 * @apiName FollowUser
	 *
	 * @apiParam {String} username Username.
	 * @apiParam {String} action Follow or Unfollow.
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
	queryUserIdTmp, _ := strconv.ParseUint(r.PostFormValue("user_id"), 10, 32)
	queryUserId := uint(queryUserIdTmp)
	queryAction := r.PostFormValue("action")

	// check user
	user, userExist, _ := database.GetUserInfoById(queryUserId, userId)
	if !userExist {
		status = 0
		errorMsg = "User not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	if user.Id == userId {
		status = 0
		errorMsg = "Cannot follow yourself."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// follow or unfollow user
	ok := database.FollowUser(user.Id, userId, queryAction)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func setUserInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/user/info/set Set user info
	 * @apiName SetUserInfo
	 *
	 * @apiParam {String} nickname Nickname.
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
	queryNickname := r.PostFormValue("nickname")

	// set Nickname by userId
	ok := database.SetUserInfo(userId, queryNickname)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get user info
	user, _, _ := database.GetUserInfoById(userId, 0)
	data = map[string]interface{}{
		"user": user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoList(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/list Get video list
	 *                            including normal, user-(uploaded, liked, favorite, watched, followed)
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
	queryType, _ := strconv.Atoi(r.PostFormValue("type"))
	queryUserIdTmp, _ := strconv.ParseUint(r.PostFormValue("user_id"), 10, 32)
	queryUserId := uint(queryUserIdTmp)
	queryRelation := r.PostFormValue("relation")
	queryLimit, _ := strconv.Atoi(r.PostFormValue("limit"))
	queryStart, _ := strconv.Atoi(r.PostFormValue("start"))

	// for some bad parameter, strict limit
	if queryLimit > 9 {
		queryLimit = 9
	} else if queryLimit < 3 {
		queryLimit = 3
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check user
	tokenValid, userId, _, _ := findAndCheckToken(r)

	// get video list
	videoList := database.GetVideoList(queryType, queryUserId, queryRelation, queryLimit, queryStart, userId)
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
				videoList[i].IsUserWatched, videoList[i].IsUserLastPlay = database.CheckUserVideoAllRelation(userId, videoList[i].Id)
		}
	}

	// if queryUserId is not 0, get video numbers
	if queryUserId != 0 {
		videoNumbers := database.GetVideoNum(queryUserId, queryRelation)
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

func getVideoInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/info Get video info
	 * @apiName GetVideoInfo
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

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")

	// check user
	tokenValid, userId, _, _ := findAndCheckToken(r)

	// get video info
	video := database.GetVideoDetailByVideoUid(queryVideoUid, userId)
	if video.Id == 0 {
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
			video.IsUserLastPlay = database.CheckUserVideoAllRelation(userId, video.Id)
	}

	data = map[string]interface{}{
		"video": video,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func doVideoAction(w http.ResponseWriter, r *http.Request) {
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// do action: like, unlike, favorite, unfavorite, forward
	ok, errNo := database.GuestDoVideoAction(videoId, userId, queryAction)
	if !ok {
		if errNo == 1 { // already done
			errorMsg = "Already done."
			status = 0
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

func recordWatchedVideo(w http.ResponseWriter, r *http.Request) {
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
	_, userId, _, _ := findAndCheckToken(r)
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
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// record video history
	ok := database.GuestWatchedVideo(videoId, userId)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func guestForwardVideo(w http.ResponseWriter, r *http.Request) {
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
	_, userId, _, _ := findAndCheckToken(r)
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
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// record video history
	ok := database.GuestForwardVideo(videoId, userId)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/comment/list Get video comment list
	 * @apiName GetVideoComment
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {Number} limit Max number of comments.
	 * @apiParam {Number} start Start at.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check token
	tokenValid, _, _, _ := findAndCheckToken(r)
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
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")
	queryLimit, _ := strconv.Atoi(r.PostFormValue("limit"))
	queryStart, _ := strconv.Atoi(r.PostFormValue("start"))

	// for some bad parameter, strict limit to 20 per page
	if queryLimit > 20 {
		queryLimit = 20
	} else if queryLimit < 10 {
		queryLimit = 10
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check video (lighter)
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check user
	tokenValid, userId, _, _ := findAndCheckToken(r)

	// get video comment
	commentList := database.GetVideoCommentList(videoId, queryLimit, queryStart, userId)
	if len(commentList) == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "No more comment found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"comment_list": commentList,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func makeVideoComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/comment/make Make video comment
	 * @apiName MakeVideoComment
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {String} content Comment content.
	 * @apiParam {Number} quote_comment_id Quote comment id.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")
	queryContent := r.PostFormValue("content")
	queryQuoteCommentIdTmp, _ := strconv.Atoi(r.PostFormValue("quote_comment_id"))
	queryQuoteCommentId := uint(queryQuoteCommentIdTmp)

	// check video (lighter)
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check duplicate comment
	if database.CheckVideoCommentDuplicate(videoId, userId, queryContent, queryQuoteCommentId) {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "You have post a same comment, please do not re-post."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if quote comment is not exist
	if queryQuoteCommentId != 0 && !database.CheckVideoCommentExist(queryQuoteCommentId) {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Quote comment not exist."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// make comment
	ok := database.UserMakeVideoComment(videoId, userId, queryContent, queryQuoteCommentId)
	if !ok {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Cannot make comment."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func deleteVideoComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/comment/delete Delete video comment
	 * @apiName DeleteVideoComment
	 *
	 * @apiParam {Number} comment_id Comment id.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check token
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryCommentIdTmp, _ := strconv.Atoi(r.PostFormValue("comment_id"))
	queryCommentId := uint(queryCommentIdTmp)

	// check comment
	videoComment := database.GetVideoComment(queryCommentId, userId)
	if videoComment.Id == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Comment not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check user
	if videoComment.UserId != userId {
		fmt.Println(userId, videoComment.UserId)
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Permission denied."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// delete comment
	ok := database.GuestDeleteVideoComment(queryCommentId)

	if !ok {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Cannot delete comment."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func uploadVideo(w http.ResponseWriter, r *http.Request) {
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
	if database.CheckVideoDuplicate(queryVideoContent, queryVideoKeyword, queryVideoType, userId) {
		status = 0
		errorMsg = "Video already exists."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check video type
	if queryVideoType < 1 || !database.CheckVideoType(queryVideoType) {
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
	video, ok := database.UserCreateVideo(userId, videoUid, queryVideoType, queryVideoContent, queryVideoKeyword)
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

func uploadVideoRemote(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/uploadRemote Upload video remote (create first, upload later)
	 * @apiName UploadVideoRemote
	 *
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

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

	// check file title, type, and description
	if queryVideoContent == "" {
		status = 0
		errorMsg = "Video content cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video duplicated
	if database.CheckVideoDuplicate(queryVideoContent, queryVideoKeyword, queryVideoType, userId) {
		status = 0
		errorMsg = "Video already exists."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check video type
	if queryVideoType < 1 || !database.CheckVideoType(queryVideoType) {
		status = 0
		errorMsg = "Invalid video type."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// generate video uid
	videoUid := uuid.New().String()

	// save video info to database
	video, ok := database.UserCreateVideo(userId, videoUid, queryVideoType, queryVideoContent, queryVideoKeyword)
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

func setVideoInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/info/set Set video info
	 * @apiName SetVideoInfo
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {String} video_content Video content.
	 * @apiParam {String} video_keyword Video keywords.
	 * @apiParam {Number} video_type Video type.
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
	queryVideoContent := r.PostFormValue("video_content")
	queryVideoKeyword := r.PostFormValue("video_keyword")
	queryVideoTypeTmp, _ := strconv.Atoi(r.PostFormValue("video_type"))
	queryVideoType := int8(queryVideoTypeTmp)

	// check video (lighter)
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	if database.CheckUserVideoRelation(userId, videoId, "uploaded") {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// edit video info
	ok := database.SetVideoInfo(queryVideoUid, queryVideoContent, queryVideoKeyword, queryVideoType)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get video info
	video := database.GetVideoDetailByVideoUid(queryVideoUid, userId)

	data = map[string]interface{}{
		"video": video,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func deleteVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/delete Delete video
	 * @apiName DeleteVideo
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
	tokenValid, userId, _, _ := findAndCheckToken(r)
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
	videoId := database.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if video is not uploaded by this user
	if database.CheckUserVideoRelation(userId, videoId, "uploaded") {
		status = 0
		errorMsg = "This video was not uploaded by this user."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// delete video
	ok := database.DeleteVideo(queryVideoUid)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getVideoTypes(w http.ResponseWriter, r *http.Request) {
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
	videoTypes := database.GetVideoTypes()
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

func searchVideo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get|post} /v1/video/search Search video
	 * @apiName SearchVideo
	 *
	 * @apiParam {String} keyword Keyword.
	 * @apiParam {Number} limit Max number of videos.
	 * @apiParam {Number} start Start at.
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
	queryKeyword := r.PostFormValue("keyword")
	queryLimit, _ := strconv.Atoi(r.PostFormValue("limit"))
	queryStart, _ := strconv.Atoi(r.PostFormValue("start"))
	if queryKeyword == "" {
		status = 0
		errorMsg = "Keyword cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// for some bad parameter, strict limit
	if queryLimit > 9 {
		queryLimit = 9
	} else if queryLimit < 3 {
		queryLimit = 3
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check if user logged in
	currentUserId := uint(0)
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if tokenValid {
		currentUserId = userId
	}

	// search video
	videoList := database.SearchVideo(queryKeyword, queryLimit, queryStart, currentUserId)
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
				videoList[i].IsUserWatched, videoList[i].IsUserLastPlay = database.CheckUserVideoAllRelation(userId, videoList[i].Id)
		}
	}

	data = map[string]interface{}{
		"video_list": videoList,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getSearchVideoHotkeys(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/search/hotkeys Get search video hotkeys
	 * @apiName GetSearchVideoHotkeys
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
	queryMaxCount, _ := strconv.Atoi(r.PostFormValue("max_count"))
	if queryMaxCount < 5 {
		queryMaxCount = 5
	}
	if queryMaxCount > 20 {
		queryMaxCount = 20
	}

	// get hotkeys
	hotkeys := database.GetSearchVideoHotkeys(queryMaxCount)
	if len(hotkeys) == 0 {
		status = 0
		errorMsg = "No hotkeys found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"hotkeys": hotkeys,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func searchUser(w http.ResponseWriter, r *http.Request) {
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
	queryKeyword := r.PostFormValue("keyword")
	queryLimit, _ := strconv.Atoi(r.PostFormValue("limit"))
	queryStart, _ := strconv.Atoi(r.PostFormValue("start"))
	if queryKeyword == "" {
		status = 0
		errorMsg = "Keyword cannot be empty."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// for some bad parameter, strict limit
	if queryLimit > 9 {
		queryLimit = 9
	} else if queryLimit < 3 {
		queryLimit = 3
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check if user logged in
	currentUserId := uint(0)
	tokenValid, userId, _, _ := findAndCheckToken(r)
	if tokenValid {
		currentUserId = userId
	}

	// search user
	userList := database.SearchUser(queryKeyword, queryLimit, queryStart, currentUserId)
	if len(userList) == 0 {
		status = 0
		errorMsg = "No more user found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	var visibleInfoUserList []model.User
	for i := 0; i < len(userList); i++ {
		visibleInfoUserList = append(visibleInfoUserList, userList[i])
	}
	data = map[string]interface{}{
		"user_list": visibleInfoUserList,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getSearchUserHotkeys(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get|post} /v1/user/search/hotkeys Get search user hotkeys
	 * @apiName GetSearchUserHotkeys
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

	// for some bad parameter, strict limit
	queryMaxCount, _ := strconv.Atoi(r.PostFormValue("max_count"))
	if queryMaxCount < 5 {
		queryMaxCount = 5
	}
	if queryMaxCount > 20 {
		queryMaxCount = 20
	}

	// get hotkeys
	hotkeys := database.GetSearchUserHotkeys(queryMaxCount)
	if len(hotkeys) == 0 {
		status = 0
		errorMsg = "No hotkeys found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"hotkeys": hotkeys,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func qiniuHlsCallback(w http.ResponseWriter, r *http.Request) {
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

	// check post data
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
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "failed to parse json:", string(postData))
		}
		status = 0
		errorMsg = "Failed to parse json."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	file := qiniuCallbackData.Input.KodoFile.Key

	// get video uid
	videoUid := strings.Split(file, "/")[2][:36]

	// get video info
	videoId := database.GetVideoIdByVideoUid(videoUid)
	if videoId == 0 {
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

	SendJSONResponse(w, status, data, errorMsg)
}

func qiniuScreenshotCallback(w http.ResponseWriter, r *http.Request) {
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

	// check post data
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
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "failed to parse json:", string(postData))
		}
		status = 0
		errorMsg = "Failed to parse json."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	file := qiniuCallbackData.Input.KodoFile.Key

	// get video uid
	videoUid := strings.Split(file, "/")[2][:36]

	// get video info
	videoId := database.GetVideoIdByVideoUid(videoUid)
	if videoId == 0 {
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

	SendJSONResponse(w, status, data, errorMsg)
}

func setToken(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
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
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"token":   token,
		"user_id": userId,
		"exp":     tool.UnixTimeToRFC3339(exp),
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
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
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	if !isExist {
		status = 0
		errorMsg = "Token not exist."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	userId, exp, err := auth.GetInfoFromToken(token)
	if err != nil {
		status = 0
		errorMsg = "Failed to parse token."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"user_id": userId,
		"exp":     tool.UnixTimeToRFC3339(exp),
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func main() {

	config.InitConfig()

	database.InitDBPool()
	defer database.CloseDBPool()
	database.InitRedisPool()
	defer database.CloseRedisPool()
	auth.InitEcdsa()

	http.HandleFunc("/", defaultPage)
	http.HandleFunc("/v1/user/login", login)
	http.HandleFunc("/v1/user/logout", logout)
	http.HandleFunc("/v1/user/signup", signup)
	http.HandleFunc("/v1/user/info", getUserInfo)
	http.HandleFunc("/v1/user/query", getOtherUserInfo)
	http.HandleFunc("/v1/user/follow", followUser)
	http.HandleFunc("/v1/user/info/set", setUserInfo)

	http.HandleFunc("/v1/video/list", getVideoList)
	http.HandleFunc("/v1/video/info", getVideoInfo)
	http.HandleFunc("/v1/video/action", doVideoAction)
	http.HandleFunc("/v1/video/watch", recordWatchedVideo)
	http.HandleFunc("/v1/video/forward", guestForwardVideo)

	http.HandleFunc("/v1/video/comment/list", getVideoComment)
	http.HandleFunc("/v1/video/comment/make", makeVideoComment)
	http.HandleFunc("/v1/video/comment/delete", deleteVideoComment)

	http.HandleFunc("/v1/video/upload", uploadVideo)
	http.HandleFunc("/v1/video/uploadRemote", uploadVideoRemote)
	http.HandleFunc("/v1/video/info/set", setVideoInfo)
	http.HandleFunc("/v1/video/delete", deleteVideo)

	http.HandleFunc("/v1/video/types", getVideoTypes)

	http.HandleFunc("/v1/video/search", searchVideo)
	http.HandleFunc("/v1/video/search/hotkeys", getSearchVideoHotkeys)
	http.HandleFunc("/v1/user/search", searchUser)
	http.HandleFunc("/v1/user/search/hotkeys", getSearchUserHotkeys)

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
