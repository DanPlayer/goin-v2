package tiktok

import (
	"encoding/json"
	"flying-star/utils"
	"fmt"
)

type VideoListResponse struct {
	Data struct {
		Cursor      int64           `json:"cursor"`
		Description string          `json:"description"`
		ErrorCode   int64           `json:"error_code"`
		HasMore     bool            `json:"has_more"`
		List        []VideoListItem `json:"list"`
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

type VideoListItem struct {
	Cover      string `json:"cover"`
	CreateTime int64  `json:"create_time"`
	IsReviewed bool   `json:"is_reviewed"`
	IsTop      bool   `json:"is_top"`
	ItemId     string `json:"item_id"`
	ShareUrl   string `json:"share_url"`
	Statistics struct {
		CommentCount  int64 `json:"comment_count"`
		DiggCount     int64 `json:"digg_count"`
		DownloadCount int64 `json:"download_count"`
		ForwardCount  int64 `json:"forward_count"`
		PlayCount     int64 `json:"play_count"`
		ShareCount    int64 `json:"share_count"`
	} `json:"statistics"`
	Title       string `json:"title"`
	VideoStatus int32  `json:"video_status"`
}

// VideoList 查询授权账号视频数据
func (tk *TikTok) VideoList(openID, accessToken string, cursor, count int64) (response VideoListResponse, err error) {
	url := fmt.Sprintf("%s/video/list?open_id=%s&access_token=%s&cursor=%d&count=%d", baseUrl, openID, accessToken, cursor, count)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return
	}
	_ = json.Unmarshal(body, &response)
	return
}

type CreateVideoRequest struct {
	VideoId             string   `json:"video_id"`               // video_id, 通过/video/upload/接口得到。注意每次调用/video/create/都要调用/video/upload/生成新的video_id。
	Text                string   `json:"text"`                   // 视频标题， 可以带话题,@用户。注意：话题审核依旧遵循抖音的审核逻辑，强烈建议第三方谨慎拟定话题名称，避免强导流行为。 title1#话题1 #话题2 @openid1
	PoiId               string   `json:"poi_id"`                 // 地理位置id，poi_id可通过"查询POI信息"能力获取
	PoiName             string   `json:"poi_name"`               // 地理位置名称
	MicroAppId          string   `json:"micro_app_id"`           // 小程序id
	MicroAppTitle       string   `json:"micro_app_title"`        // 小程序标题
	CoverTsp            float64  `json:"cover_tsp"`              // 将传入的指定时间点对应帧设置为视频封面（单位：秒）
	AtUsers             []string `json:"at_users"`               // 如果需要at其他用户。将text中@nickname对应的open_id放到这里。
	MicroAppUrl         string   `json:"micro_app_url"`          // 开发者在小程序中生成该页面时写的path地址
	CustomCoverImageUrl string   `json:"custom_cover_image_url"` // 自定义封面图片,参数为接口/image/upload/ 返回的image_id
}

type CreateVideoResponse struct {
	Data struct {
		ItemId string `json:"item_id"`
	} `json:"data"`
	Extra struct {
		ErrorCode      int    `json:"error_code"`
		Description    string `json:"description"`
		SubErrorCode   int    `json:"sub_error_code"`
		SubDescription string `json:"sub_description"`
		Logid          string `json:"logid"`
		Now            int64  `json:"now"`
	} `json:"extra"`
}

// CreateVideo 创建视频
func (tk *TikTok) CreateVideo(openID, accessToken string, body CreateVideoRequest) (ret CreateVideoResponse, err error) {
	url := fmt.Sprintf("%s/video/create?open_id=%s&access_token=%s", baseUrl, openID, accessToken)
	marshal, _ := json.Marshal(body)
	data, err := utils.HttpPostBody(url, marshal)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &ret)
	return
}

type VideosByIdsRequest struct {
	ItemIds []string `json:"item_ids"`
}

type VideosByIdsResponse struct {
	Data struct {
		List []struct {
			Cover      string `json:"cover"`
			CreateTime int64  `json:"create_time"`
			IsReviewed bool   `json:"is_reviewed"`
			IsTop      bool   `json:"is_top"`
			Statistics struct {
				CommentCount  int `json:"comment_count"`
				DiggCount     int `json:"digg_count"`
				DownloadCount int `json:"download_count"`
				ForwardCount  int `json:"forward_count"`
				PlayCount     int `json:"play_count"`
				ShareCount    int `json:"share_count"`
			} `json:"statistics"`
			VideoStatus int32  `json:"video_status"`
			ItemId      string `json:"item_id"`
			ShareUrl    string `json:"share_url"`
			Title       string `json:"title"`
		} `json:"list"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	} `json:"data"`
	Extra struct {
		Logid          string `json:"logid"`
		ErrorCode      int    `json:"error_code"`
		Description    string `json:"description"`
		SubErrorCode   int    `json:"sub_error_code"`
		SubDescription string `json:"sub_description"`
		Now            int    `json:"now"`
	} `json:"extra"`
}

// VideosByIds 获取视频通过ids
func (tk *TikTok) VideosByIds(openID, accessToken string, ids []string) (ret VideosByIdsResponse, err error) {
	var body VideosByIdsRequest
	body.ItemIds = ids
	url := fmt.Sprintf("%s/video/data?open_id=%s&access_token=%s", baseUrl, openID, accessToken)
	marshal, _ := json.Marshal(body)
	data, err := utils.HttpPostBody(url, marshal)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &ret)
	return
}
