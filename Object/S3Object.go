package Object

import (
	"S3ObjectStorageFileBrowser/Object/Config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"strconv"
)

/*
func GetObject(ObjectClient *minio.Client, bucket string, prefix string, objectSize bool, objectETage bool, objectLastModified bool) ([]Config.ObjectInfo, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := ObjectClient.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	var objectList []Config.ObjectInfo
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
		objectList = append(objectList, objectinfo)
	}
	return objectList, nil
}*/
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
func getObjectSizeSuitableUnit(size int64) string {
	if size%125 >= 1 && size < 1000 {
		unit := float64(size) / 125
		return strconv.FormatFloat(unit, 'f', 2, 64) + " Kb"
	} else if size%1000 >= 1 && size < 1000000 {
		unit := float64(size) / 1000
		return strconv.FormatFloat(unit, 'f', 2, 64) + " KB"
	} else if size%1000000 > 1 && size < 1000000000 {
		unit := float64(size) / 1000000
		return strconv.FormatFloat(unit, 'f', 2, 64) + " MB"
	} else if size%1000000000 > 1 && size < 1000000000000 {
		unit := float64(size) / 1000000000
		return strconv.FormatFloat(unit, 'f', 2, 64) + " GB"
	} else {
		unit := float64(size) / 1000000000000
		return strconv.FormatFloat(unit, 'f', 2, 64) + " TB"
	}
}
