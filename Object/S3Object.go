package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func GetObject(ObjectClient *minio.Client, bucket string, prefix string, objectSize bool, objectETage bool, objectLastModified bool) (Config.ObjectInfoList, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := ObjectClient.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	objectList := Config.ObjectInfoList{}
	for object := range objectCh {
		objectinfo := Config.ObjectInfo{}
		if object.Err != nil {
			return objectList, object.Err
		}
		objectinfo.Key = object.Key
		if objectSize {
			objectinfo.Size = object.Size
		}
		if objectETage {
			objectinfo.ETag = object.ETag
		}
		if objectLastModified {
			objectinfo.LsatModified = object.LastModified
		}
		objectList.Info = append(objectList.Info, objectinfo)
	}
	return objectList, nil
}
func MakeClient(server Config.Server) (*minio.Client, error) {

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
	//log.Printf("%#v\n", minioClient)

	return minioClient, err
}
