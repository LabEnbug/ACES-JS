package tool

import (
	"backend/config"
	"context"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

func UploadFileToQiniu(localPath string, remotePath string) bool {
	upToken := config.QiniuPutPolicy.UploadToken(config.QiniuMac)
	formUploader := storage.NewFormUploader(&config.QiniuStorageCfg)
	ret := storage.PutRet{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, remotePath, localPath, nil)
	if err != nil {
		funcName, _, _, _ := runtime.Caller(0)
		log.Println(runtime.FuncForPC(funcName).Name(), err)
		return false
	}
	funcName, _, _, _ := runtime.Caller(0)
	log.Println(runtime.FuncForPC(funcName).Name(), "upload success", ret.Key, ret.Hash)
	return true
}

func UploadDirToQiniu(dirName string) bool {
	// {{video_uid}} -> /video/{{video_uid}}/
	localDirPath := path.Join(config.BaseLocalFileDir, dirName)
	//remoteDirPath := path.Join(config.BaseRemoteFileDir, dirName) // seems cannot make dir by api?
	remoteDirPath := path.Join(config.BaseRemoteFileDir, "") // just put in root dir
	// search all files in localDirPath
	var files []string
	err := filepath.Walk(localDirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return false
	}
	// upload all files
	for _, file := range files {
		remoteFilePath := path.Join(remoteDirPath, file[len(localDirPath):])
		// try upload 3 times
		uploadOk := false
		for i := 0; i < 3; i++ {
			if UploadFileToQiniu(file, remoteFilePath) {
				uploadOk = true
				break
			}
		}
		if !uploadOk {
			return false
		}
	}
	return true
}
