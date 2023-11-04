package cmd

import (
	"backend/config"
	"backend/database/mysql"
	"backend/model"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func QiniuHlsCallback(w http.ResponseWriter, r *http.Request) {
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
	videoId := mysql.GetVideoIdByVideoUid(videoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// update video info
	ok := mysql.CallbackUpdateVideoHLS(videoUid)
	if !ok {
		status = 0
		errorMsg = "Failed to update video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func QiniuScreenshotCallback(w http.ResponseWriter, r *http.Request) {
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
	videoId := mysql.GetVideoIdByVideoUid(videoUid)
	if videoId == 0 {
		status = 0
		errorMsg = "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// update video info
	ok := mysql.CallbackUpdateVideoScreenshot(videoUid)
	if !ok {
		status = 0
		errorMsg = "Failed to update video info."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}
