package config

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"time"
)

type DatabaseConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string
}

type RedisConfig struct {
	Host    string
	Port    int
	Pass    string
	Channel int
}

var Test bool // set to false to disable test api
var ShowLog bool

var TokenExpireTime time.Duration // token expire time
var TokenMaxDevice int64          // max logged in devices

var DefaultDatabase DatabaseConfig

var DefaultRedis RedisConfig

var QiniuStorageCfg storage.Config
var QiniuPutPolicy storage.PutPolicy
var QiniuMac *auth.Credentials

var BaseLocalFileDir string  // local dir for temp file
var BaseRemoteFileDir string // remote dir

var VideoProcessType int // 1: wait for qiniu transcode callback; 2: transcode by ffmpeg locally

var MaxNormalPostSize int   // MB
var MaxUploadVideoSize int  // MB
var MaxUploadAvatarSize int // MB
var MaxNormalPostSize64 int64
var MaxUploadVideoSize64 int64
var MaxUploadAvatarSize64 int64

var PromotePrice float64
var AdvertisePrice float64

func InitConfig() {
	Test = true
	ShowLog = true

	TokenExpireTime = time.Hour * 24 * 3 // 3 days
	TokenMaxDevice = 2                   // 2 devices

	VideoProcessType = 1

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

	MaxNormalPostSize = 10
	MaxUploadVideoSize = 200
	MaxUploadAvatarSize = 2
	MaxNormalPostSize64 = int64(MaxNormalPostSize * 1024 * 1024)
	MaxUploadVideoSize64 = int64(MaxUploadVideoSize * 1024 * 1024)
	MaxUploadAvatarSize64 = int64(MaxUploadAvatarSize * 1024 * 1024)

	PromotePrice = 0.1
	AdvertisePrice = 0.3
}
