package cmd

import (
	"backend/auth"
	"backend/config"
	"backend/database"
	"backend/database/mysql"
	"backend/tool"
	"net/http"
	"strconv"
)

func Login(w http.ResponseWriter, r *http.Request) {
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
	user, ok, errNo := mysql.CheckUserPassword(queryUsername, queryPassword)
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
	mysql.SetLastLogin(user.Id)
	status = 200
	data = map[string]interface{}{
		"token": token,
		"exp":   tool.UnixTimeToRFC3339(exp),
		"user":  user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func Logout(w http.ResponseWriter, r *http.Request) {
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
	tokenValid, _, _, token := FindAndCheckToken(r)
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

func Signup(w http.ResponseWriter, r *http.Request) {
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
	user, ok, errNo := mysql.CreateUser(queryUsername, queryPassword, queryNickname)
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
	user, _, _ = mysql.GetUserInfoById(user.Id, 0)
	data = map[string]interface{}{
		"user": user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
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
	tokenValid, userId, exp, token := FindAndCheckToken(r)
	if !tokenValid {
		status = 0
		errorMsg = "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get user info
	user, ok, errNo := mysql.GetUserInfoById(userId, 0)
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

func GetOtherUserInfo(w http.ResponseWriter, r *http.Request) {
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
	_, userId, _, _ := FindAndCheckToken(r)

	// get user info
	user, userExist := mysql.GetUserInfoByUsername(queryUsername, userId)
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

func FollowUser(w http.ResponseWriter, r *http.Request) {
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
	queryUserIdTmp, _ := strconv.ParseUint(r.PostFormValue("user_id"), 10, 32)
	queryUserId := uint(queryUserIdTmp)
	queryAction := r.PostFormValue("action")

	// check user
	queryUser, userExist, _ := mysql.GetUserInfoById(queryUserId, userId)
	if !userExist {
		status = 0
		errorMsg = "User not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	if queryUser.Id == userId {
		status = 0
		errorMsg = "Cannot follow yourself."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if already followed
	isFollowed := mysql.CheckUserFollowed(queryUserId, userId)
	if queryAction == "follow" && isFollowed {
		status = 200
		errorMsg = "Already followed."
		SendJSONResponse(w, status, data, errorMsg)
		return
	} else if queryAction == "unfollow" && !isFollowed {
		status = 200
		errorMsg = "Already unfollowed."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// follow or unfollow user
	ok := mysql.FollowUser(queryUserId, userId, queryAction)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func SetUserInfo(w http.ResponseWriter, r *http.Request) {
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
	queryNickname := r.PostFormValue("nickname")

	// set Nickname by userId
	ok := mysql.SetUserInfo(userId, queryNickname)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get user info
	user, _, _ := mysql.GetUserInfoById(userId, 0)
	data = map[string]interface{}{
		"user": user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}
