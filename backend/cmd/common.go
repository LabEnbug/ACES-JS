package cmd

import (
	"backend/auth"
	"backend/database"
	"github.com/golang-jwt/jwt/v5/request"
	"net/http"
)

func FindAndCheckToken(r *http.Request) (bool, uint, int64, string) {
	// find token
	// Authorization: Bearer xxx
	token, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		return false, 0, 0, ""
	}

	// check token
	isExist, _ := database.CheckTokenIsExist(token)
	if isExist {
		userId, exp, err := auth.GetInfoFromToken(token)
		if err == nil {
			return true, userId, exp, token
		}
	}
	return false, 0, 0, ""
}
