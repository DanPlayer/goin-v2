package tiktok

import (
	"encoding/json"
	"flying-star/utils"
	"fmt"
)

type FansResponse struct {
	Data struct {
		Cursor      int64  `json:"cursor"`
		Description string `json:"description"`
		ErrorCode   string `json:"error_code"`
		HasMore     bool   `json:"has_more"`
		List        []struct {
			Avatar   string `json:"avatar"`
			City     string `json:"city"`
			Country  string `json:"country"`
			Gender   int64  `json:"gender"`
			Nickname string `json:"nickname"`
			OpenId   string `json:"open_id"`
			Province string `json:"province"`
			UnionId  string `json:"union_id"`
		} `json:"list"`
		Total int64 `json:"total"`
	} `json:"data"`
	Extra struct {
		Description    string `json:"description"`
		ErrorCode      int64  `json:"error_code"`
		Logid          string `json:"logid"`
		Now            int64  `json:"now"`
		SubDescription string `json:"sub_description"`
		SubErrorCode   int64  `json:"sub_error_code"`
	} `json:"extra"`
}

// FansList 粉丝列表
func (tk *TikTok) FansList(openID, accessToken string, cursor, count int64) (response FansResponse, err error) {
	url := fmt.Sprintf("%s/fans/list/?open_id=%s&access_token=%s&cursor=%d&count=%d", baseUrl, openID, accessToken, cursor, count)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return
	}
	_ = json.Unmarshal(body, &response)
	return
}

type FollowingResponse struct {
	Data struct {
		Cursor      int64  `json:"cursor"`
		Description string `json:"description"`
		ErrorCode   string `json:"error_code"`
		HasMore     bool   `json:"has_more"`
		List        []struct {
			Avatar   string `json:"avatar"`
			City     string `json:"city"`
			Country  string `json:"country"`
			Gender   int64  `json:"gender"`
			Nickname string `json:"nickname"`
			OpenId   string `json:"open_id"`
			Province string `json:"province"`
			UnionId  string `json:"union_id"`
		} `json:"list"`
		Total int64 `json:"total"`
	} `json:"data"`
	Extra struct {
		Description    string `json:"description"`
		ErrorCode      int64  `json:"error_code"`
		Logid          string `json:"logid"`
		Now            int64  `json:"now"`
		SubDescription string `json:"sub_description"`
		SubErrorCode   int64  `json:"sub_error_code"`
	} `json:"extra"`
}


// FollowingList 关注列表
func (tk *TikTok) FollowingList(openID, accessToken string, cursor, count int64) (response FollowingResponse, err error) {
	url := fmt.Sprintf("%s/following/list/?open_id=%s&access_token=%s&cursor=%d&count=%d", baseUrl, openID, accessToken, cursor, count)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return
	}
	_ = json.Unmarshal(body, &response)
	return
}