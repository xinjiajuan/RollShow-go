package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func MakeClient(server Config.Server) *minio.Client {

	// Initialize minio client object.
	minioClient, err := minio.New(server.Host, &minio.Options{
		Creds:        credentials.NewStaticV4(server.AccessKeyID, server.SecretAccessKey, ""),
		Secure:       server.Options.UseSSL,
		Region:       server.Options.Region,
		BucketLookup: server.Options.BucketLookupType,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%#v\n", minioClient)
	
	return minioClient
}
