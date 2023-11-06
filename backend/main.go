package main

import (
	"backend/auth"
	"backend/cmd"
	"backend/config"
	"backend/database"
	"backend/database/mysql"
	"fmt"
	"github.com/gorilla/mux"
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

	r := mux.NewRouter()

	// user
	r.HandleFunc("/v1/user/login", cmd.Login).Methods("POST")
	r.HandleFunc("/v1/user/logout", cmd.Logout).Methods("POST")
	r.HandleFunc("/v1/user/signup", cmd.Signup).Methods("POST")
	r.HandleFunc("/v1/user/info", cmd.GetUserInfo).Methods("GET")
	r.HandleFunc("/v1/user/info", cmd.SetUserInfo).Methods("PUT")
	r.HandleFunc("/v1/users/{username}", cmd.GetOtherUserInfo).Methods("GET")
	r.HandleFunc("/v1/user/follow", cmd.FollowUser).Methods("POST")
	r.HandleFunc("/v1/user/deposit", cmd.UserDeposit).Methods("POST")

	// video
	http.HandleFunc("/v1/video/userRecomList", cmd.GetRecommendVideoList)
	r.HandleFunc("/v1/videos", cmd.GetVideoList).Methods("GET")
	r.HandleFunc("/v1/video/{videoUid}", cmd.GetVideoInfo).Methods("GET")
	r.HandleFunc("/v1/video/{videoUid}", cmd.SetVideoInfo).Methods("PUT")
	r.HandleFunc("/v1/video/{videoUid}", cmd.DeleteVideo).Methods("DELETE")
	r.HandleFunc("/v1/video/action", cmd.DoVideoAction).Methods("POST")
	r.HandleFunc("/v1/video/watch", cmd.RecordWatchedVideo).Methods("POST")
	r.HandleFunc("/v1/video/forward", cmd.GuestForwardVideo).Methods("POST")

	r.HandleFunc("/v1/video/upload", cmd.UploadVideo).Methods("POST")                // [deprecated] file and info add together
	r.HandleFunc("/v1/video/uploadRemote", cmd.UploadVideoRemote).Methods("POST")    // add info first, upload to qiniu directly from source, admin usage
	r.HandleFunc("/v1/video/upload/file", cmd.UploadVideoFile).Methods("POST")       // upload file first, save in server, user usage step 1
	r.HandleFunc("/v1/video/upload/confirm", cmd.ConfirmVideoUpload).Methods("POST") // confirm publish, upload from server to qiniu, user usage step 2
	r.HandleFunc("/v1/video/top", cmd.TopVideo).Methods("POST")
	r.HandleFunc("/v1/video/private", cmd.PrivateVideo).Methods("POST")
	r.HandleFunc("/v1/video/types", cmd.GetVideoTypes).Methods("POST")

	r.HandleFunc("/v1/video/promote", cmd.PromoteVideo).Methods("POST")
	r.HandleFunc("/v1/video/advertise", cmd.AdvertiseVideo).Methods("POST")

	// video comment
	r.HandleFunc("/v1/video/comment/list", cmd.GetVideoCommentList).Methods("POST")
	r.HandleFunc("/v1/video/comment/make", cmd.MakeVideoComment).Methods("POST")
	r.HandleFunc("/v1/video/comment/delete", cmd.DeleteVideoComment).Methods("POST")

	// video bullet comment
	r.HandleFunc("/v1/video/bullet_comment/list", cmd.GetVideoBulletCommentList).Methods("POST")
	r.HandleFunc("/v1/video/bullet_comment/make", cmd.MakeVideoBulletComment).Methods("POST")
	r.HandleFunc("/v1/video/bullet_comment/delete", cmd.DeleteVideoBulletComment).Methods("POST")

	// search
	r.HandleFunc("/v1/search/video", cmd.SearchVideo).Methods("GET")
	r.HandleFunc("/v1/search/video/hotkeys", cmd.GetSearchVideoHotkeys).Methods("GET")
	r.HandleFunc("/v1/search/user", cmd.SearchUser).Methods("GET")
	r.HandleFunc("/v1/search/user/hotkeys", cmd.GetSearchUserHotkeys).Methods("GET")

	// callback
	r.HandleFunc("/callback/qiniu/hls", cmd.QiniuHlsCallback).Methods("POST")
	r.HandleFunc("/callback/qiniu/screenshot", cmd.QiniuScreenshotCallback).Methods("POST")

	http.Handle("/", r)

	err := http.ListenAndServe(":8051", nil)
	if err != nil {
		log.Fatal("Startup err: ", err)
	}

}
