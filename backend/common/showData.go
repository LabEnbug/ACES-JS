package common

import (
	"backend/model"
	"backend/tool"
)

func GetVisibleUserInfo(user model.User) map[string]interface{} {
	return map[string]interface{}{
		"user_id":  user.Id,
		"username": user.Username,
		"nickname": user.Nickname,
		"reg_time": tool.DatabaseTimeToRFC3339(user.RegTime),
	}
}
