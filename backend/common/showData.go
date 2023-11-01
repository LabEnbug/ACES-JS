package common

import (
	"backend/model"
	"backend/tool"
	"fmt"
)

func MakeUserSupInfo(user model.User) model.User {
	user.RegTimeStr = tool.DatabaseTimeToRFC3339(user.RegTime)
	return user
}

func MakeVideoSupInfo(video model.Video, screenshotOk uint8, hlsOk uint8, user model.User) model.Video {
	if screenshotOk == 1 {
		video.CoverUrl = fmt.Sprintf("http://s348vstvo.bkt.clouddn.com/video/%s/cover.png", video.VideoUid)
	}
	if hlsOk == 1 {
		video.PlayUrl = fmt.Sprintf("http://s348vstvo.bkt.clouddn.com/video/%s/index.m3u8", video.VideoUid)
	}
	video.UploadTimeStr = tool.DatabaseTimeToRFC3339(video.UploadTime)
	video.User = user
	return video
}
