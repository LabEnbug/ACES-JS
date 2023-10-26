package model

type QiniuCallbackData struct {
	Version  string `json:"version"`
	Id       string `json:"id"`
	Reqid    string `json:"reqid"`
	Pipeline string `json:"pipeline"`
	Input    struct {
		KodoFile struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		} `json:"kodo_file"`
	} `json:"input"`
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Ops  []struct {
		Id  string `json:"id"`
		Fop struct {
			Cmd       string `json:"cmd"`
			InputFrom string `json:"input_from"`
			Result    struct {
				Code      int    `json:"code"`
				Desc      string `json:"desc"`
				HasOutput bool   `json:"has_output"`
				KodoFile  struct {
					Bucket string `json:"bucket"`
					Key    string `json:"key"`
					Hash   string `json:"hash"`
				} `json:"kodo_file"`
			} `json:"result"`
		} `json:"fop"`
	} `json:"ops"`
	CreatedAt int64 `json:"created_at"`
}
