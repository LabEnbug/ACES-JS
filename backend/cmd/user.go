package cmd

import (
	"backend/auth"
	"backend/config"
	"backend/database"
	"backend/database/mysql"
	"backend/tool"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
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
	 * @api {get} /v1/user/info Get user info
	 * @apiName GetUserInfo
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	//// check method
	//if r.Method != "GET" && r.Method != "POST" {
	//	status = 0
	//	errorMsg = "Invalid request method."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

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

	user.IsSelf = true
	user.Balance = mysql.GetUserBalance(userId)

	data = map[string]interface{}{
		"token": token,
		"exp":   tool.UnixTimeToRFC3339(exp),
		"user":  user,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GetOtherUserInfo(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get} /v1/users/{username} Get other user info
	 * @apiName GetOtherUserInfo
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
	//queryUsername := r.PostFormValue("username")
	vars := mux.Vars(r)
	queryUsername := vars["username"]
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
	 * @api {post|delete} /v1/users/{username}/follow Follow user
	 * @apiName FollowUser
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
	queryUsername := vars["username"]

	// check user
	user, userExist := mysql.GetUserInfoByUsername(queryUsername, userId)
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

	// check if already followed
	isFollowed := mysql.CheckUserFollowed(user.Id, userId)
	if r.Method == "POST" && isFollowed {
		status = 200
		errorMsg = "Already followed."
		SendJSONResponse(w, status, data, errorMsg)
		return
	} else if r.Method == "DELETE" && !isFollowed {
		status = 200
		errorMsg = "Already unfollowed."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// follow or unfollow user
	var ok bool
	if r.Method == "POST" {
		ok = mysql.FollowUser(user.Id, userId, "follow")
	} else if r.Method == "DELETE" {
		ok = mysql.FollowUser(user.Id, userId, "unfollow")
	}

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
	 * @api {put} /v1/user/info/set Set user info
	 * @apiName SetUserInfoNickname
	 *
	 * @apiParam {String} type Type. (nickname || avatar)
	 * @apiParam {String} content Content of nickname.
	 * @apiParam {file} file File of avatar.
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

	// check avatar file size by MaxBytesReader, limit to {{config.MaxUploadAvatarSize64}}MB
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadAvatarSize64)
	err := r.ParseMultipartForm(config.MaxUploadAvatarSize64)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
		}
		status = 0
		errorMsg = fmt.Sprintf("File size limit to %dMB.", config.MaxUploadAvatarSize)
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryType := r.FormValue("type")
	if queryType == "nickname" {
		queryNickname := r.FormValue("nickname")
		// set Nickname by userId
		ok := mysql.SetUserInfoNickname(userId, queryNickname)
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
	} else if queryType == "avatar" {
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
		extensionList := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".ico"}
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

		// generate avatar uid
		avatarUid := uuid.New().String()

		// save file to {{BaseLocalFileDir}}/tmp/{{avatar_uid}}
		fileSavePath := path.Join(config.BaseLocalFileDir, "tmp", avatarUid+".origin.img")
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

		// transcode image to png
		file, err = os.Open(fileSavePath)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			status = 0
			errorMsg = "Failed to open file."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		var img image.Image
		img, _, err = image.Decode(file)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			status = 0
			errorMsg = "Failed to decode image."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		// resize to 128x128 for storage saving if over 128x128
		if img.Bounds().Dx() > 128 || img.Bounds().Dy() > 128 {
			resizedImg := image.NewRGBA(image.Rect(0, 0, 128, 128))
			draw.ApproxBiLinear.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
			img = resizedImg
		}

		var outFile *os.File
		outPath := path.Join(config.BaseLocalFileDir, "tmp", avatarUid+".png")
		outFile, err = os.Create(outPath)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			status = 0
			errorMsg = "Failed to create file."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		err = png.Encode(outFile, img)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "err: ", err)
			}
			status = 0
			errorMsg = "Failed to encode image."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		subDir := "avatar"

		// upload file to qiniu
		// try to upload 3 times
		uploadOk := false
		for i := 0; i < 3; i++ {
			ok := tool.UploadFileToQiniu(outPath, path.Join(config.BaseRemoteFileDir, subDir, avatarUid+".png"))
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

		ok := mysql.SetUserInfoAvatar(userId, avatarUid)
		if !ok {
			status = 0
			errorMsg = "Unknown error."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}

		user, _, _ := mysql.GetUserInfoById(userId, userId)

		// just return, wait for callback asynchronously
		data = map[string]interface{}{
			"user": user,
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func UserDeposit(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/user/deposit User deposit
	 * @apiName UserDeposit
	 *
	 * @apiParam {string} card_key Card key.
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
	queryCardKey := r.PostFormValue("card_key")

	// check card
	card := mysql.GetDepositCard(queryCardKey)
	if card.Id == 0 {
		status = 0
		errorMsg = "Card not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	if card.UsedTime.Valid {
		status = 0
		errorMsg = "Card already used."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// deposit
	ok := mysql.UserDeposit(userId, queryCardKey)
	if !ok {
		status = 0
		errorMsg = "Unknown error."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	data = map[string]interface{}{
		"deposit_amount": card.CardAmount,
	}

	SendJSONResponse(w, status, data, errorMsg)
}
