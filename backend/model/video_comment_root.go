package model

type VideoCommentRoot struct {
	// Id, UserId, Content, QuoteRootCommentId, QuoteUserId, QuoteChildCommentId, CommentTime
	VideoCommentChild

	ChildCommentList      []VideoCommentChild `default:"[]" json:"child_comment_list"`
	ChildCommentCountLeft int                 `default:"0" json:"child_comment_count_left"`
}
