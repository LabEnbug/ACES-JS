package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

var Test bool
var ShowLog bool

type DatabaseConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

var DefaultDatabase DatabaseConfig

type RedisConfig struct {
	Host    string
	Port    int
	Pass    string
	Channel int
}

var DefaultRedis RedisConfig

var QiniuStorageCfg storage.Config
var QiniuPutPolicy storage.PutPolicy
var QiniuMac *auth.Credentials

var BaseLocalFileDir string
var BaseRemoteFileDir string

var VideoProcessType int

var MaxUploadVideoSize int

func InitConfig() {
	Test = true // set to false to disable test api
	ShowLog = true
	VideoProcessType = 1 // 1: wait for qiniu transcode callback; 2: transcode by ffmpeg locally

	DefaultDatabase = DatabaseConfig{
		Host: "127.0.0.1",
		Port: 3306,
		User: "root",
		Pass: "acesaces",
		Name: "qiniu_sv",
	}
	DefaultRedis = RedisConfig{
		Host:    "127.0.0.1",
		Port:    6379,
		Pass:    "acesaces",
		Channel: 0,
	}
	QiniuStorageCfg = storage.Config{
		Region:        &storage.ZoneHuadongZheJiang2,
		UseHTTPS:      true,
		UseCdnDomains: false,
	}

	QiniuPutPolicy = storage.PutPolicy{
		Scope: "aces-js",
	}
	QiniuMac = qbox.NewMac("mWIWThPUX4LBqsEz8UvduI9DUhmfSEfuVPQ2VzPG", "s5JKrAkJWl4XQpNnP42qHAeEmtfZ75gRPYkSyAkd")
	BaseLocalFileDir = "/root/project/user_upload_files"
	BaseRemoteFileDir = ""

	MaxUploadVideoSize = 50 // MB
}
