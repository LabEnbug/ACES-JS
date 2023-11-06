package mysql

import (
	"backend/auth"
	"backend/common"
	"backend/config"
	"backend/model"
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"runtime"
	"strconv"
)

func GetUserBalance(userId uint) float64 {
	var balance float64
	err := DB.QueryRow("SELECT balance FROM user WHERE id=? LIMIT 1", userId).Scan(&balance)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
	}
	return balance
}

func GetUserInfoById(userId uint, currentUserId uint) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	err := DB.QueryRow("SELECT id, username, password, nickname, avatar, follow_count, be_followed_count, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, reg_time, last_login_time FROM user WHERE id=? LIMIT 1", userId).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.Avatar, &user.FollowCount, &user.BeFollowedCount, &user.BeLikedCount, &user.BeFavoriteCount, &user.BeCommentedCount, &user.BeForwardedCount, &user.BeWatchedCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		if errors.Is(err, sql.ErrNoRows) { // no such user
			ok = false
			errNo = 1
		} else {
			ok = false
			errNo = 9999
		}
	}
	if currentUserId != 0 {
		user.BeFollowed = CheckUserFollowed(user.Id, currentUserId)
		user.IsSelf = user.Id == currentUserId
	}
	user = common.MakeUserSupInfo(user)
	return user, ok, errNo
}

func GetUserInfoByUsername(username string, currentUserId uint) (model.User, bool) {
	var user model.User
	err := DB.QueryRow("SELECT id, username, password, nickname, avatar, follow_count, be_followed_count, be_liked_count, be_favorite_count, be_commented_count, be_forwarded_count, be_watched_count, reg_time, last_login_time FROM user WHERE username=? LIMIT 1", username).
		Scan(&user.Id, &user.Username, &user.Password, &user.Nickname, &user.Avatar, &user.FollowCount, &user.BeFollowedCount, &user.BeLikedCount, &user.BeFavoriteCount, &user.BeCommentedCount, &user.BeForwardedCount, &user.BeWatchedCount, &user.RegTime, &user.LastLoginTime)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return user, false
	}

	if currentUserId != 0 {
		user.BeFollowed = CheckUserFollowed(user.Id, currentUserId)
		user.IsSelf = user.Id == currentUserId
	}
	user = common.MakeUserSupInfo(user)
	return user, true
}

func CheckUserPassword(username string, password string) (model.User, bool, int) {
	ok := true
	errNo := 0

	user, userExist := GetUserInfoByUsername(username, 0)

	if !userExist {
		ok = false
		errNo = 1
	} else if !auth.ValidatePassword(password, user.Password) {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR:", "wrong password")
		}
		ok = false
		errNo = 2
	}
	// no need to check be followed
	user = common.MakeUserSupInfo(user)
	return user, ok, errNo
}

func CreateUser(username string, password string, nickname string) (model.User, bool, int) {
	var user model.User
	ok := true
	errNo := 0

	_, userExist := GetUserInfoByUsername(username, 0)

	if !userExist { // no such user
		// create user
		if nickname == "" { // empty nickname
			nickname = "用户" + strconv.Itoa(rand.Intn(100000)+100000) // random nickname
		}
		// encode password
		encodedPassword := auth.MakePassword(password)
		res, err := DB.Exec("INSERT INTO user (username, password, nickname, reg_time) VALUES (?, ?, ?, NOW())",
			username, encodedPassword, nickname)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
			}
			ok = false
			errNo = 9999
		} else {
			lastId, err := res.LastInsertId()
			if err != nil {
				if config.ShowLog {
					funcName, _, _, _ := runtime.Caller(0)
					log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
				}
				ok = false
				errNo = 9999
			} else {
				if config.ShowLog {
					funcName, _, _, _ := runtime.Caller(0)
					log.Println(runtime.FuncForPC(funcName).Name(), "signup lastId: ", lastId)
				}
				user.Id = uint(lastId)
				user.Username = username
				user.Password = password
				user.Nickname = nickname
			}
		}
	} else { // username exist
		ok = false
		errNo = 1
	}
	// no need to check be followed
	user = common.MakeUserSupInfo(user)
	return user, ok, errNo
}

func SetLastLogin(userId uint) {
	_, err := DB.Exec("UPDATE user SET last_login_time=NOW() WHERE id=?", userId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
	}
}

func CheckUserFollowed(followUserId uint, currentUserId uint) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(1) FROM user_follow WHERE follow_user_id=? AND user_id=? AND unfollow_time IS NULL LIMIT 1", followUserId, currentUserId).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

func FollowUser(followUserId uint, currentUserId uint, action string) bool {
	if action == "follow" {
		_, err := DB.Exec("INSERT INTO user_follow (follow_user_id, user_id, follow_time) VALUES (?, ?, NOW())", followUserId, currentUserId)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return false
		}
	} else if action == "unfollow" {
		_, err := DB.Exec("UPDATE user_follow SET unfollow_time=NOW() WHERE follow_user_id=? AND user_id=? AND unfollow_time IS NULL", followUserId, currentUserId)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return false
		}
	} else {
		return false
	}
	return true
}

func SetUserInfoNickname(userId uint, nickname string) bool {
	_, err := DB.Exec("UPDATE user SET nickname=? WHERE id=?", nickname, userId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func SetUserInfoAvatar(userId uint, avatar string) bool {
	_, err := DB.Exec("UPDATE user SET avatar=? WHERE id=?", avatar, userId)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err)
		}
		return false
	}
	return true
}

func GetDepositCard(cardKey string) model.UserDepositCard {
	var card model.UserDepositCard
	err := DB.QueryRow("SELECT id, card_key, card_amount, used_user_id, used_time FROM user_deposit_card WHERE card_key=?", cardKey).
		Scan(&card.Id, &card.CardKey, &card.CardAmount, &card.UsedUserId, &card.UsedTime)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
	}
	return card
}

func UserDeposit(userId uint, cardKey string) bool {
	_, err := DB.Exec("UPDATE user_deposit_card SET used_user_id=?, used_time=NOW() WHERE card_key=?", userId, cardKey)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return false
	}

	// todo: test only
	if cardKey == "ACES-AAAA-AAAA-AAAA" {
		_, err := DB.Exec("UPDATE user_deposit_card SET used_user_id=NULL, used_time=NULL WHERE card_key=?", cardKey)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return false
		}
	}
	return true
}
