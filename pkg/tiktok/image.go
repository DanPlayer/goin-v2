package tiktok

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

// ImageUploadParams 图片上传参数
type ImageUploadParams struct {
	Image []byte `json:"image"`
}

// ImageUploadSchema 图片上传响应参数
type ImageUploadSchema struct {
	Extra struct {
		ErrorCode      int    `json:"error_code"`
		Description    string `json:"description"`
		SubErrorCode   int    `json:"sub_error_code"`
		SubDescription string `json:"sub_description"`
		Logid          string `json:"logid"`
		Now            int64  `json:"now"`
	} `json:"extra"`
	Data struct{
		ErrorCode     int64  `json:"error_code"`
		Image struct{
			Height 	   int64 `json:"height"`
			ImageId    string `json:"image_id"`
			Width      int64 `json:"width"`
		} `json:"image"`
		Description   string `json:"description"`
	} `json:"data"`
}

// ImageUpload 上传图片到文件服务器
func (tk *TikTok) ImageUpload(openID, accessToken string, image []byte, fileName string) (res ImageUploadSchema, err error) {
	url := fmt.Sprintf("%s/image/upload?open_id=%s&access_token=%s", baseUrl, openID, accessToken)

	b := bytes.Buffer{}
	writer := multipart.NewWriter(&b)

	head := make(textproto.MIMEHeader)
	head.Set("Content-Type", "image/"+strings.Replace(fileName, ".", "", -1))
	head.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="%s"`, fileName))
	if _, err = writer.CreatePart(head); err != nil {
		return res, err
	}

	lastLine := fmt.Sprintf("\r\n--%s--\r\n", writer.Boundary())
	r := strings.NewReader(lastLine)

	bodyLen := int64(b.Len()) + int64(len(image)) + int64(len(lastLine))
	mr := io.MultiReader(&b, bytes.NewReader(image), r)
	contentType := writer.FormDataContentType()
	headers := http.Header{}
	headers.Add("Content-Type", contentType)

	defer writer.Close()

	req, err := http.NewRequest("POST", url, mr)
	if err != nil {
		fmt.Println("req err: ", err)
		return
	}
	req.ContentLength = bodyLen
	req.Header = headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("resp err: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("抖音接口请求失败，错误码:" + strconv.Itoa(resp.StatusCode))
		return
	}

	data, _ := ioutil.ReadAll(resp.Body)

	_ = json.Unmarshal(data, &res)
	return
}