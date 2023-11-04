package cmd

import (
	"backend/config"
	"backend/database/mysql"
	"backend/model"
	"net/http"
	"strconv"
)

func SearchVideo(w http.ResponseWriter, r *http.Request) {
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
	if queryLimit > 24 {
		queryLimit = 24
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check if user logged in
	currentUserId := uint(0)
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if tokenValid {
		currentUserId = userId
	}

	// search video
	videoList := mysql.SearchVideo(queryKeyword, queryLimit, queryStart, currentUserId)
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

func GetSearchVideoHotkeys(w http.ResponseWriter, r *http.Request) {
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
	hotkeys := mysql.GetSearchVideoHotkeys(queryMaxCount)
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

func SearchUser(w http.ResponseWriter, r *http.Request) {
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
	if queryLimit > 24 {
		queryLimit = 24
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check if user logged in
	currentUserId := uint(0)
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if tokenValid {
		currentUserId = userId
	}

	// search user
	userList := mysql.SearchUser(queryKeyword, queryLimit, queryStart, currentUserId)
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

func GetSearchUserHotkeys(w http.ResponseWriter, r *http.Request) {
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
	hotkeys := mysql.GetSearchUserHotkeys(queryMaxCount)
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
