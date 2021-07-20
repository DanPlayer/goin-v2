package aliyun

import "github.com/aliyun/alibaba-cloud-sdk-go/sdk"

type Options struct {
	RegionId        string `json:"region_id"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type AliYun struct {
	Options Options     `json:"options"`
	Client  *sdk.Client `json:"client"`
}

// New 初始化阿里云服务
func New(options Options) (*AliYun, error) {
	var y AliYun
	client, err := sdk.NewClientWithAccessKey(options.RegionId, options.AccessKeyId, options.AccessKeySecret)
	y.Client = client
	return &y, err
}
