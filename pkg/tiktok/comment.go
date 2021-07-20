package tiktok

import (
	"encoding/json"
	"flying-star/utils"
	"fmt"
	"net/url"
)

type CommentListResponse struct {
	Data struct {
		HasMore bool `json:"has_more"`
		List    []struct {
			CreateTime        int64  `json:"create_time"`
			DiggCount         int32  `json:"digg_count"`
			ReplyCommentTotal int32  `json:"reply_comment_total"`
			Top               bool   `json:"top"`
			CommentId         string `json:"comment_id"`
			CommentUserId     string `json:"comment_user_id"`
			Content           string `json:"content"`
		} `json:"list"`
		Cursor int64 `json:"cursor"`
	} `json:"data"`
	Extra struct {
		Description    string `json:"description"`
		SubErrorCode   int    `json:"sub_error_code"`
		SubDescription string `json:"sub_description"`
		Logid          string `json:"logid"`
		Now            int64  `json:"now"`
		ErrorCode      int    `json:"error_code"`
	} `json:"extra"`
}

// CommentList 评论列表
func (tk *TikTok) CommentList(openID, accessToken, videoID string, cursor, count int64) (response CommentListResponse, err error) {
	videoID = url.QueryEscape(videoID)
	url := fmt.Sprintf("%s/item/comment/list/?open_id=%s&item_id=%s&access_token=%s&cursor=%d&count=%d&sort_type=time", baseUrl, openID, videoID, accessToken, cursor, count)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return
	}
	_ = json.Unmarshal(body, &response)
	return
}

type ReplyCommentRequest struct {
	CommentID string `json:"comment_id"`
	Content   string `json:"content"`
	ItemID    string `json:"item_id"`
}

type ReplyCommentResponse struct {
	Data struct {
		CommentId   string `json:"comment_id"`
		Description string `json:"description"`
		ErrorCode   string `json:"error_code"`
	} `json:"data"`
	Extra struct {
		Description    string `json:"description"`
		ErrorCode      string `json:"error_code"`
		Logid          string `json:"logid"`
		Now            string `json:"now"`
		SubDescription string `json:"sub_description"`
		SubErrorCode   string `json:"sub_error_code"`
	} `json:"extra"`
	Message string `json:"message"`
}

// ReplyComment 回复评论
func (tk *TikTok) ReplyComment(openID, accessToken string, request ReplyCommentRequest) (ret ReplyCommentResponse, err error) {
	url := fmt.Sprintf("%s/item/comment/reply/?open_id=%s&access_token=%s", baseUrl, openID, accessToken)
	body, _ := json.Marshal(request)
	data, err := utils.HttpPostBody(url, body)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &ret)
	return
}
