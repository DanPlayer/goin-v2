package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApiResponseData struct {
	Rtn  int                    `json:"rtn"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

type ApiResponseMsg struct {
	Rtn int    `json:"rtn"`
	Msg string `json:"msg"`
}

// HttpGet 根据struct返回不同形式定义的接口数据
func (s *ApiResponseData) HttpGet(url string) (resp ApiResponseData, err error) {
	apiResponse, err := HttpGetBody(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(apiResponse, &resp)
	if err != nil {
		err = errors.New("Unmarshal failed " + err.Error())
		return
	}

	return resp, err
}

// HttpGet 根据struct返回不同形式定义的接口数据
func (s *ApiResponseMsg) HttpGet(url string) (resp ApiResponseMsg, err error) {
	apiResponse, err := HttpGetBody(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(apiResponse, &resp)
	if err != nil {
		err = errors.New("Unmarshal failed " + err.Error())
		return
	}

	return resp, err
}

// HttpGetBody 通用的获取http的Body
func HttpGetBody(url string) (body []byte, err error) {
	res, _ := http.Get(url)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = errors.New("响应错误:" + res.Status)
		return
	}

	return ioutil.ReadAll(res.Body)
}

// HttpPostBody http的post body json请求获取响应
func HttpPostBody(url string, body []byte) (info []byte, err error) {
	res, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("响应失败，检查网络状况")
		return
	}

	return ioutil.ReadAll(res.Body)
}

// HttpPostWithForm http的post body form请求获取响应
func HttpPostWithForm(url string, body *bytes.Buffer, contentType string) (info []byte, err error) {
	res, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode)
		err = errors.New("响应失败，检查网络状况")
		return
	}

	return ioutil.ReadAll(res.Body)
}
