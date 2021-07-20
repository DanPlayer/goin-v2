package oss

import (
	"bytes"
	"flying-star/utils"
	OSS "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-sts-go-sdk/sts"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Options 初始化参数
type Options struct {
	//访问密钥
	AccessId     string
	AccessSecret string
	EndPoint     string
	RoleArn      string
	ExpireTIme   int64
}

type Client struct {
	opt Options
	oss *OSS.Client
}

type Credentials struct {
	Region        string `json:"region"` //访问区域
	BucketName    string `json:"bucketName"`
	AccessId      string `json:"accessId"`
	AccessSecret  string `json:"accessSecret"`
	SecurityToken string `json:"securityToken"`
	Expiration    uint   `json:"expiration"`
}

// NewClient 初始化客户端
func NewClient(options Options) (client *Client, err error) {
	oss, err := OSS.New(options.EndPoint, options.AccessId, options.AccessSecret)
	if err != nil {
		return nil, err
	}
	client = &Client{
		opt: options,
		oss: oss,
	}
	return client, nil
}

// GetSTS 获取STS授权证书
func (c *Client) GetSTS(bucketName string, sessionName string) (cred Credentials, err error) {
	region, err := c.oss.GetBucketLocation(bucketName)
	if err != nil {
		return cred, err
	}
	client := sts.NewClient(c.opt.AccessId, c.opt.AccessSecret, c.opt.RoleArn, sessionName)
	res, err := client.AssumeRole(uint(c.opt.ExpireTIme))
	if err != nil {
		return cred, err
	}
	cred = Credentials{
		Region:        region,
		BucketName:    bucketName,
		AccessId:      res.Credentials.AccessKeyId,
		AccessSecret:  res.Credentials.AccessKeySecret,
		SecurityToken: res.Credentials.SecurityToken,
		Expiration:    uint(res.Credentials.Expiration.Unix()),
	}
	return cred, nil
}

// GetObject 流式下载文件到本地
func (c *Client) GetObject(bucketName string, objectName string) (data []byte, err error) {
	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	body, err := bucket.GetObject(objectName)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	data, err = ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// UploadNetworkImage 上传网络图片到OSS
func (c *Client) UploadNetworkImage(bucketName string, link string) (file FileDetail, err error) {
	return c.UploadNetworkFile(bucketName, link, "image/jpg")
}

// UploadNetworkAmr 上传网络MP3到OSS
func (c *Client) UploadNetworkAmr(bucketName string, link string) (file FileDetail, err error) {
	return c.UploadNetworkFile(bucketName, link, "amr")
}

// UploadNetworkFile 上传网络文件到OSS
func (c *Client) UploadNetworkFile(bucketName string, link string, mimeType string) (file FileDetail, err error) {
	response, err := http.Get(link)
	if err != nil {
		return file, err
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	mimeList := strings.Split(mimeType, "/")
	if len(mimeList) == 2 {
		file.Format = mimeList[1]
	} else {
		file.Format = mimeType
	}
	//计算文件hash
	file.Name = utils.MD5(data) + "." + file.Format

	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return file, err
	}

	reader := bytes.NewReader(data)
	file.Size = int(reader.Size())
	osMeta := OSS.Meta("size", strconv.Itoa(int(reader.Size())))
	contentType := OSS.ContentType(mimeType)

	if header, _ := bucket.GetObjectDetailedMeta(file.Name); header != nil {
		sizeStr := header.Get("X-Oss-Meta-Size")
		if sizeStr == "" {
			_ = bucket.SetObjectMeta(file.Name, osMeta)
		}
		return file, nil
	}

	err = bucket.PutObject(file.Name, reader, osMeta, contentType)
	return file, err
}

// MultipartCallback 分片上传文件
type MultipartCallback func(mediaId string, indexBuf string) (MediaData, error)

func (c *Client) UploadMultipartFile(bucketName string, fileName string, mediaId string, fn MultipartCallback) error {
	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return err
	}

	//如果文件已存在，则直接返回当前文件信息
	if header, _ := bucket.GetObjectDetailedMeta(fileName); header != nil {
		return nil
	}

	isFinish := false
	indexBuf := ""
	repeat := 3
	chunkIndex := 1
	parts := make([]OSS.UploadPart, 0)

	storageType := OSS.ObjectStorageClass(OSS.StorageStandard)
	mulObj, err := bucket.InitiateMultipartUpload(fileName, storageType)
	if err != nil {
		return err
	}

	for !isFinish {
		mediaData, err := fn(mediaId, indexBuf)
		if err != nil {
			//重试3次
			if repeat > 0 {
				repeat--
				continue
			}
			return err
		}

		if mediaData.IsFinish {
			isFinish = mediaData.IsFinish
		}
		indexBuf = mediaData.OutIndexBuf

		reader := bytes.NewReader(mediaData.Data)
		part, err := bucket.UploadPart(mulObj, reader, reader.Size(), chunkIndex)
		if err != nil {
			return err
		}
		parts = append(parts, part)
		chunkIndex++
	}

	objectAcl := OSS.ObjectACL(OSS.ACLPublicRead)
	_, err = bucket.CompleteMultipartUpload(mulObj, parts, objectAcl)
	return err
}

// UploadFile 上传文件
func (c *Client) UploadFile(bucketName string, filePath string, mimeType string) (file FileDetail, err error) {
	data, _ := ioutil.ReadFile(filePath)
	file.Format = mimeType
	//计算文件hash
	file.Name = utils.MD5(data) + "." + file.Format

	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return file, err
	}

	reader := bytes.NewReader(data)
	file.Size = int(reader.Size())
	osMeta := OSS.Meta("size", strconv.Itoa(int(reader.Size())))
	contentType := OSS.ContentType(file.Format)

	if header, _ := bucket.GetObjectDetailedMeta(file.Name); header != nil {
		sizeStr := header.Get("X-Oss-Meta-Size")
		if sizeStr == "" {
			_ = bucket.SetObjectMeta(file.Name, osMeta)
		}
		return file, nil
	}

	err = bucket.PutObject(file.Name, reader, osMeta, contentType)
	return file, err
}
