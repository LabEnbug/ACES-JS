package mysql

import (
	"backend/config"
	"log"
	"runtime"
)

func CallbackUpdateVideoHLS(videoUid string) bool {
	if _, err := DB.Exec("UPDATE video SET hls=1 WHERE video_uid=?", videoUid); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func CallbackUpdateVideoScreenshot(videoUid string) bool {
	if _, err := DB.Exec("UPDATE video SET screenshot=1 WHERE video_uid=?", videoUid); err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}
