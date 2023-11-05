package mysql

import (
	"backend/config"
	"backend/model"
	"backend/tool"
	"database/sql"
	"log"
	"runtime"
)

//func CheckVideoCommentDuplicate(videoId uint, userId uint, content string, quoteCommentId uint) bool {
//	var count int
//	ERR := DB.QueryRow("SELECT COUNT(1) FROM video_comment WHERE user_id=? AND video_id=? AND content=? AND quote_comment_id=? LIMIT 1", userId, videoId, content, quoteCommentId).Scan(&count)
//	if ERR != nil || count == 0 {
//		return false
//	}
//	return true
//}

func MakeVideoComment(videoId uint, userId uint, content string, quoteComment model.VideoCommentChild) (int64, bool) {
	var res sql.Result
	var err error
	var commentId int64
	if quoteComment.Id == 0 { // 1: make root comment
		res, err = DB.Exec("INSERT INTO video_comment (video_id, user_id, content, comment_time) VALUES (?, ?, ?, NOW())", videoId, userId, content)
	} else {
		if quoteComment.QuoteRootCommentId == 0 { // 2: make child comment (quote root comment)
			res, err = DB.Exec("INSERT INTO video_comment (video_id, user_id, content, quote_root_comment_id, quote_user_id, comment_time) VALUES (?, ?, ?, ?, ?, NOW())", videoId, userId, content, quoteComment.Id, quoteComment.UserId)
		} else { // 3: grandchild comment (quote child comment)
			res, err = DB.Exec("INSERT INTO video_comment (video_id, user_id, content, quote_root_comment_id, quote_user_id, quote_child_comment_id, comment_time) VALUES (?, ?, ?, ?, ?, ?, NOW())", videoId, userId, content, quoteComment.QuoteRootCommentId, quoteComment.UserId, quoteComment.Id)
		}
	}

	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), "ERR: ", err, "videoId:", videoId, "userId:", userId, "content:", content, "quoteComment:", quoteComment)
		}
		return commentId, false
	}

	// get comment id
	commentId, err = res.LastInsertId()
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err, "videoId:", videoId, "userId:", userId, "content:", content, "quoteComment:", quoteComment)
		}
		return commentId, false
	}

	return commentId, true
}

func GetVideoComment(videoCommentId uint, currentUserId uint) model.VideoCommentChild {
	var videoComment model.VideoCommentChild
	err := DB.QueryRow("SELECT id, user_id, content, quote_root_comment_id, quote_user_id, quote_child_comment_id, comment_time FROM video_comment WHERE id=? AND delete_time IS NULL LIMIT 1", videoCommentId).
		Scan(&videoComment.Id, &videoComment.UserId, &videoComment.Content, &videoComment.QuoteRootCommentId,
			&videoComment.QuoteUserId, &videoComment.QuoteChildCommentId, &videoComment.CommentTime)
	if err != nil {
		return videoComment
	}
	videoComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(videoComment.CommentTime)
	user, _, _ := GetUserInfoById(videoComment.UserId, currentUserId)
	videoComment.User = user
	return videoComment
}

func DeleteVideoComment(videoCommentId uint) bool {
	if _, err := DB.Exec("UPDATE video_comment SET delete_time=NOW() WHERE id=?", videoCommentId); err != nil {
		return false
	}
	return true
}

func GetVideoChildCommentList(videoId uint, quoteRootCommentId uint, limit int, start int, currentUserId uint) ([]model.VideoCommentChild, int) {
	var childCommentList []model.VideoCommentChild
	var rows *sql.Rows
	var err error

	childCommentCountLeft := 0

	rows, err = DB.Query("SELECT id, user_id, content, quote_user_id, quote_child_comment_id, comment_time FROM video_comment WHERE video_id=? AND quote_root_comment_id=? AND delete_time IS NULL ORDER BY id DESC LIMIT ?, ?", videoId, quoteRootCommentId, start, limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return childCommentList, childCommentCountLeft
	}

	defer rows.Close()

	for rows.Next() {
		var childComment model.VideoCommentChild

		// for child comment, quoteRootCommentId is always the same as the parent comment
		childComment.QuoteRootCommentId = quoteRootCommentId

		err := rows.Scan(&childComment.Id, &childComment.UserId, &childComment.Content, &childComment.QuoteUserId, &childComment.QuoteChildCommentId, &childComment.CommentTime)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return childCommentList, childCommentCountLeft
		}
		childComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(childComment.CommentTime)
		user, _, _ := GetUserInfoById(childComment.UserId, currentUserId)
		childComment.User = user
		if childComment.QuoteUserId != 0 {
			quoteUser, _, _ := GetUserInfoById(childComment.QuoteUserId, currentUserId)
			childComment.QuoteUser = quoteUser
		}

		childCommentList = append(childCommentList, childComment)
	}

	if len(childCommentList) == 0 {
		return childCommentList, childCommentCountLeft
	}

	// get comment count left
	err = DB.QueryRow("SELECT COUNT(1) FROM video_comment WHERE video_id=? AND quote_root_comment_id=? AND id<? AND delete_time IS NULL", videoId, quoteRootCommentId, childCommentList[len(childCommentList)-1].Id).
		Scan(&childCommentCountLeft)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return childCommentList, childCommentCountLeft
	}

	return childCommentList, childCommentCountLeft
}

func GetVideoRootCommentList(videoId uint, limit int, start int, currentUserId uint) ([]model.VideoCommentRoot, int) {
	var rootCommentList []model.VideoCommentRoot
	var rows *sql.Rows
	var err error
	rootCommentCountLeft := 0

	rows, err = DB.Query("SELECT id, user_id, content, comment_time FROM video_comment WHERE video_id=? AND quote_root_comment_id=0 AND delete_time IS NULL ORDER BY id DESC LIMIT ?, ?", videoId, start, limit)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return rootCommentList, rootCommentCountLeft
	}

	defer rows.Close()

	for rows.Next() {
		var rootComment model.VideoCommentRoot
		err := rows.Scan(&rootComment.Id, &rootComment.UserId, &rootComment.Content, &rootComment.CommentTime)
		if err != nil {
			if config.ShowLog {
				funcName, _, _, _ := runtime.Caller(0)
				log.Println(runtime.FuncForPC(funcName).Name(), err)
			}
			return rootCommentList, rootCommentCountLeft
		}
		rootComment.CommentTimeStr = tool.DatabaseTimeToRFC3339(rootComment.CommentTime)
		user, _, _ := GetUserInfoById(rootComment.UserId, currentUserId)
		rootComment.User = user
		// for root comment, get 1 latest child comment
		rootComment.ChildCommentList, rootComment.ChildCommentCountLeft = GetVideoChildCommentList(videoId, rootComment.Id, 1, 0, currentUserId)
		rootCommentList = append(rootCommentList, rootComment)
	}

	if len(rootCommentList) == 0 {
		return rootCommentList, rootCommentCountLeft
	}

	// get comment count left
	err = DB.QueryRow("SELECT COUNT(1) FROM video_comment WHERE video_id=? AND id<? AND quote_root_comment_id=0 AND delete_time IS NULL", videoId, rootCommentList[len(rootCommentList)-1].Id).
		Scan(&rootCommentCountLeft)
	if err != nil {
		if config.ShowLog {
			funcName, _, _, _ := runtime.Caller(0)
			log.Println(runtime.FuncForPC(funcName).Name(), err)
		}
		return rootCommentList, rootCommentCountLeft
	}

	return rootCommentList, rootCommentCountLeft
}
