package cmd

import (
	"backend/config"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
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
