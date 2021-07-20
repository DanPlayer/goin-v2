package oss

// FileDetail 文件上传类型
type FileDetail struct {
	Name   string
	Size   int
	Format string
}

type MediaData struct {
	OutIndexBuf string `json:"outindexbuf,omitempty"`
	IsFinish    bool   `json:"is_finish,omitempty"`
	Data        []byte `json:"data,omitempty"`
}
