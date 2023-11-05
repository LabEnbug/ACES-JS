package cmd

import (
	"backend/config"
	"backend/database/mysql"
	"net/http"
	"strconv"
)

func MakeVideoBulletComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/bullet_comment/make Make video bullet comment
	 * @apiName MakeVideoBulletComment
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {String} content Comment content.
	 * @apiParam {float64} comment_at Comment at the time inner the video.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")
	queryContent := r.PostFormValue("content")
	queryCommentAt, _ := strconv.ParseFloat(r.PostFormValue("comment_at"), 64)

	// check video (lighter)
	videoId := mysql.GetVideoIdByVideoUid(queryVideoUid)
	if videoId == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Video not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// make bullet comment
	bulletCommentId, ok := mysql.MakeVideoBulletComment(videoId, userId, queryContent, queryCommentAt)
	if !ok {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to make bullet comment."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// get bullet comment
	bulletComment := mysql.GetVideoBulletComment(uint(bulletCommentId), userId)
	data = map[string]interface{}{
		"bullet_comment": bulletComment,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func GetVideoBulletCommentList(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/bullet_comment/list Get video bullet comment list
	 * @apiName GetVideoBulletCommentList
	 *
	 * @apiParam {String} video_uid Video uid.
	 * @apiParam {Number} limit Max number of bullet comments.
	 * @apiParam {Number} start Start at.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryVideoUid := r.PostFormValue("video_uid")
	queryLimit, _ := strconv.Atoi(r.PostFormValue("limit"))
	queryStart, _ := strconv.Atoi(r.PostFormValue("start"))

	// for some bad parameter, strict limit to 500 per video
	if queryLimit > 500 {
		queryLimit = 500
	} else if queryLimit < 1 {
		queryLimit = 500
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

	// get bullet comment list
	bulletCommentList, bulletCommentCount := mysql.GetVideoBulletCommentList(videoId, queryLimit, queryStart, userId)
	data = map[string]interface{}{
		"bullet_comment_list":  bulletCommentList,
		"bullet_comment_count": bulletCommentCount,
	}

	SendJSONResponse(w, status, data, errorMsg)
}

func DeleteVideoBulletComment(w http.ResponseWriter, r *http.Request) {
	/*
	 * @api {post} /v1/video/bullet_comment/delete Delete video bullet comment
	 * @apiName DeleteVideoBulletComment
	 *
	 * @apiParam {int} bullet_comment_id Bullet comment id.
	 */
	status := 200
	data := map[string]interface{}{}
	errorMsg := ""

	// check method, only accept POST
	if r.Method != "POST" {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Invalid request method."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check token
	tokenValid, userId, _, _ := FindAndCheckToken(r)
	if !tokenValid {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not logged in."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// parse form
	err := r.ParseMultipartForm(config.MaxNormalPostSize64)
	if err != nil {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to parse form."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}
	queryBulletCommentIdTmp, _ := strconv.Atoi(r.PostFormValue("bullet_comment_id"))
	queryBulletCommentId := uint(queryBulletCommentIdTmp)

	// check bullet comment
	bulletComment := mysql.GetVideoBulletComment(queryBulletCommentId, userId)
	if bulletComment.Id == 0 {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Bullet comment not found."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// check if not owner
	if bulletComment.UserId != userId {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Not owner."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	// delete bullet comment
	if !mysql.DeleteVideoBulletComment(queryBulletCommentId) {
		status := 0
		data := map[string]interface{}{}
		errorMsg := "Failed to delete bullet comment."
		SendJSONResponse(w, status, data, errorMsg)
		return
	}

	SendJSONResponse(w, status, data, errorMsg)
}
