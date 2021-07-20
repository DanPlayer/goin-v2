package tiktok

import (
	"encoding/json"
	"flying-star/utils"
	"fmt"
)

type PoiSearchKeywordResponse struct {
	Data struct {
		Cursor      int64  `json:"cursor"`
		Description string `json:"description"`
		ErrorCode   int64  `json:"error_code"`
		HasMore     string `json:"has_more"`
		Pois        []struct {
			Address     string `json:"address"`
			City        string `json:"city"`
			CityCode    string `json:"city_code"`
			Country     string `json:"country"`
			CountryCode string `json:"country_code"`
			District    string `json:"district"`
			Location    string `json:"location"`
			PoiId       string `json:"poi_id"`
			PoiName     string `json:"poi_name"`
			Province    string `json:"province"`
		} `json:"pois"`
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

// PoiSearchKeyword 根据关键词搜索poi地址
func (tk *TikTok) PoiSearchKeyword(clientToken, keyword, city string, cursor, count int64) (response PoiSearchKeywordResponse, err error) {
	url := fmt.Sprintf("%s/poi/search/keyword?access_token=%s&keyword=%s&city=%s&cursor=%d&count=%d", baseUrl, clientToken, keyword, city, cursor, count)
	body, err := utils.HttpGetBody(url)
	if err != nil {
		return
	}
	_ = json.Unmarshal(body, &response)
	return
}
