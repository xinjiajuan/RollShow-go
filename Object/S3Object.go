package Object

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"rollshow/Object/Config"
	"strconv"
)

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
	if size < 1048576 {
		unit := float64(size) / 1024
		return strconv.FormatFloat(unit, 'f', 2, 64) + " KiB"
	} else if size < 1073741824 {
		unit := float64(size) / 1048576
		return strconv.FormatFloat(unit, 'f', 2, 64) + " MiB"
	} else if size < 1099511627776 {
		unit := float64(size) / 1073741824
		return strconv.FormatFloat(unit, 'f', 2, 64) + " GiB"
	} else {
		unit := float64(size) / 1099511627776
		return strconv.FormatFloat(unit, 'f', 2, 64) + " TiB"
	}
}
