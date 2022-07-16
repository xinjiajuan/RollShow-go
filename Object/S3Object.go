package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func MakeClientObject(config Config.Yaml) {
	var serverList []Config.Server
	for _, list := range config.ServerList {
		serverList = append(serverList, list)
	}
	for _, server := range serverList {
		go MakeClient(server)
	}
}

func GetObjectInfo(ObjectClient minio.Client, objectSize bool) []string {
	endpoint := "192.168.2.220:9000"
	accessKeyID := "API user"
	secretAccessKey := "16885886hzq"
	opt := "china-jx-gz-01"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
		Region: opt,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%#v\n", minioClient)
	return nil
}
