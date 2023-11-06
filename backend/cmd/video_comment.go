package cmd

import (
	"backend/config"
	"backend/database/mysql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetVideoCommentList(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {get} /v1/videos/{videoUid}/comments Get video comment list
	 * @apiName GetVideoCommentList
	 *
	 * @apiParam {Number} limit Max number of comments.
	 * @apiParam {Number} start Start at.
	 * @apiParam {Number} comment_id Get child comment.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	queryParams := r.URL.Query()
	queryLimit, _ := strconv.Atoi(queryParams.Get("limit"))
	queryStart, _ := strconv.Atoi(queryParams.Get("start"))
	queryCommentIdTmp, _ := strconv.Atoi(queryParams.Get("comment_id"))
	queryRootCommentId := uint(queryCommentIdTmp)

	// for some bad parameter, strict limit to 20 per page
	if queryLimit > 24 {
		queryLimit = 24
	} else if queryLimit < 1 {
		queryLimit = 1
	}
	if queryStart < 0 {
		queryStart = 0
	}

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check user
	_, userId, _, _ := FindAndCheckToken(r)

	// get video comment
	if queryRootCommentId == 0 { // root comment
		rootCommentList, rootCommentCountLeft := mysql.GetVideoRootCommentList(videoId, queryLimit, queryStart, userId)
		if len(rootCommentList) == 0 {
			status := 0
			data := map[string]interface{}{}
			errorMsg := "No more comment found."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
		data = map[string]interface{}{
			"comment_list":       rootCommentList,
			"comment_count_left": rootCommentCountLeft,
		}
	} else { // child comment
		childCommentList, childCommentCountLeft := mysql.GetVideoChildCommentList(videoId, queryRootCommentId, queryLimit, queryStart, userId)
		if len(childCommentList) == 0 {
			status := 0
			data := map[string]interface{}{}
			errorMsg := "No more comment found."
			SendJSONResponse(w, status, data, errorMsg)
			return
		}
		data = map[string]interface{}{
			"child_comment_list":       childCommentList,
			"child_comment_count_left": childCommentCountLeft,
		}
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func MakeVideoComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/videos/{videoUid}/comments Make video comment
	 * @apiName MakeVideoComment
	 *
	 * @apiParam {String} content Comment content.
	 * @apiParam {Number} quote_comment_id Quote comment id.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	vars := mux.Vars(r)
	queryVideoUid := vars["videoUid"]

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryContent := r.FormValue("content")
	queryQuoteCommentIdTmp, _ := strconv.Atoi(r.FormValue("quote_comment_id"))
	queryQuoteCommentId := uint(queryQuoteCommentIdTmp)

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if quote comment is not exist
	quoteComment := mysql.GetVideoComment(queryQuoteCommentId, userId)
	if queryQuoteCommentId != 0 && quoteComment.Id == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Quote comment not exist."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check duplicate comment
	// todo: need to be fixed
	//if database.CheckVideoCommentDuplicate(videoId, userId, queryContent, queryQuoteCommentId) {
	//	status := 0
	//	data := map[string]interface{}{}
	//	errorMsg := "You have post a same comment, please do not re-post."
	//	SendJSONResponse(w, status, data, errorMsg)
	//	return
	//}

	// make comment, 1: root comment, 2: child comment, 3: grandchild comment
	commentId, ok := mysql.MakeVideoComment(videoId, userId, queryContent, quoteComment)
	if !ok {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Cannot make comment."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get comment, no need to check root or child
	comment := mysql.GetVideoComment(uint(commentId), userId)
	data = map[string]interface{}{
		"comment": comment,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func DeleteVideoComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {delete} /v1/videos/{videoUid}/comments/{commentId} Delete video comment
	 * @apiName DeleteVideoComment
	 *
	 * @apiParam {Number} comment_id Comment id.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	vars := mux.Vars(r)
	//queryVideoUid := vars["videoUid"]
	queryCommentIdTmp, _ := strconv.Atoi(vars["commentId"])
	queryCommentId := uint(queryCommentIdTmp)

	// check comment
	videoComment := mysql.GetVideoComment(queryCommentId, userId)
	if videoComment.Id == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Comment not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check user
	if videoComment.UserId != userId {
		fmt.Println(userId, videoComment.UserId)
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Permission denied."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// delete comment
	ok := mysql.DeleteVideoComment(queryCommentId)

	if !ok {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Cannot delete comment."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}
