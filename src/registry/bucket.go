package registry

import (
	"log"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nbisso/storicard-challenge/infrastracture/conf"
)

var Instance minio.Client
var miniodbonce sync.Once

func (r *register) NewMinIOClient() minio.Client {

	miniodbonce.Do(func() {
		endpoint := conf.Instance.Minio.Endpoint
		accessKeyID := conf.Instance.Minio.AccessKey
		secretAccessKey := conf.Instance.Minio.SecretKey
		useSSL := conf.Instance.Minio.UseSSL
		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalln(err)
		}

		Instance = *minioClient
	})

	return Instance

}
