package Config

import (
	"github.com/minio/minio-go/v7"
	"time"
)

type Yaml struct {
	ServerList []Server `yaml:"server"`
}
type Server struct {
	Name            string `yaml:"name"`
	ListenPort      int    `yaml:"listenPort"`
	Enable          bool   `yaml:"enable"`
	Host            string `yaml:"host"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Bucket          string `yaml:"bucket"`
	Options         struct {
		UseSSL                   bool                   `yaml:"useSSL"`
		Region                   string                 `yaml:"region"`
		BucketLookupType         minio.BucketLookupType `yaml:"bucketLookupType"` //DNS,Path:1,Auto:0
		AccessControlAllowOrigin string                 `yaml:"access-control-allow-origin"`
		Favicon                  string                 `yaml:"favicon"`
		BeianMiit                string                 `yaml:"beianMiit"`
	} `yaml:"options"`
}

type ObjectInfo struct {
	Num          int
	Key          string
	Size         int64
	ETag         string
	LsatModified time.Time
}

var (
	Version string = "1.1.6"
	AppName string = "RollShow"
	Usage   string = "基于S3对象储存文件下载服务器"
)
