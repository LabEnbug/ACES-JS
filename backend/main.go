package main

import (
	"backend/auth"
	"backend/cmd"
	"backend/config"
	"backend/database"
	"backend/database/mysql"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

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

func main() {

	config.InitConfig()

	mysql.InitDBPool()
	defer mysql.CloseDBPool()

	database.InitRedisPool()
	defer database.CloseRedisPool()

	auth.InitEcdsa()

	http.HandleFunc("/", defaultPage)

	// user
	http.HandleFunc("/v1/user/login", cmd.Login)
	http.HandleFunc("/v1/user/logout", cmd.Logout)
	http.HandleFunc("/v1/user/signup", cmd.Signup)
	http.HandleFunc("/v1/user/info", cmd.GetUserInfo)
	http.HandleFunc("/v1/user/query", cmd.GetOtherUserInfo)
	http.HandleFunc("/v1/user/follow", cmd.FollowUser)
	http.HandleFunc("/v1/user/info/set", cmd.SetUserInfo)
	http.HandleFunc("/v1/user/deposit", cmd.UserDeposit)

	// video
	http.HandleFunc("/v1/video/list", cmd.GetVideoList)
	http.HandleFunc("/v1/video/info", cmd.GetVideoInfo)
	http.HandleFunc("/v1/video/action", cmd.DoVideoAction)
	http.HandleFunc("/v1/video/watch", cmd.RecordWatchedVideo)
	http.HandleFunc("/v1/video/forward", cmd.GuestForwardVideo)

	http.HandleFunc("/v1/video/upload", cmd.UploadVideo)                // [deprecated] file and info add together
	http.HandleFunc("/v1/video/uploadRemote", cmd.UploadVideoRemote)    // add info first, upload to qiniu directly from source, admin usage
	http.HandleFunc("/v1/video/upload/file", cmd.UploadVideoFile)       // upload file first, save in server, user usage step 1
	http.HandleFunc("/v1/video/upload/confirm", cmd.ConfirmVideoUpload) // confirm publish, upload from server to qiniu, user usage step 2
	http.HandleFunc("/v1/video/info/set", cmd.SetVideoInfo)
	http.HandleFunc("/v1/video/delete", cmd.DeleteVideo)
	http.HandleFunc("/v1/video/top", cmd.TopVideo)
	http.HandleFunc("/v1/video/private", cmd.PrivateVideo)
	http.HandleFunc("/v1/video/types", cmd.GetVideoTypes)

	http.HandleFunc("/v1/video/promote", cmd.PromoteVideo)
	http.HandleFunc("/v1/video/advertise", cmd.AdvertiseVideo)

	// video comment
	http.HandleFunc("/v1/video/comment/list", cmd.GetVideoCommentList)
	http.HandleFunc("/v1/video/comment/make", cmd.MakeVideoComment)
	http.HandleFunc("/v1/video/comment/delete", cmd.DeleteVideoComment)

	// video bullet comment
	http.HandleFunc("/v1/video/bullet_comment/list", cmd.GetVideoBulletCommentList)
	http.HandleFunc("/v1/video/bullet_comment/make", cmd.MakeVideoBulletComment)
	http.HandleFunc("/v1/video/bullet_comment/delete", cmd.DeleteVideoBulletComment)

	// search
	http.HandleFunc("/v1/video/search", cmd.SearchVideo)
	http.HandleFunc("/v1/video/search/hotkeys", cmd.GetSearchVideoHotkeys)
	http.HandleFunc("/v1/user/search", cmd.SearchUser)
	http.HandleFunc("/v1/user/search/hotkeys", cmd.GetSearchUserHotkeys)

	// callback
	http.HandleFunc("/callback/qiniu/hls", cmd.QiniuHlsCallback)
	http.HandleFunc("/callback/qiniu/screenshot", cmd.QiniuScreenshotCallback)

	err := http.ListenAndServe(":8051", nil)
	if err != nil {
		log.Fatal("Startup err: ", err)
	}

}
