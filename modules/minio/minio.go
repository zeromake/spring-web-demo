package minio

import (
	"fmt"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	SpringCore "github.com/go-spring/go-spring/spring-core"
	"github.com/minio/minio-go/v6"
)

type MinioConfig struct {
	Enable bool   `value:"${minio.enable:=true}"`    // 是否启用 HTTP
	Host   string `value:"${minio.host:=127.0.0.1}"` // HTTP host
	Port   int    `value:"${minio.port:=9000}"`      // HTTP 端口
	Access string `value:"${minio.access:=}"`        // Access
	Secret string `value:"${minio.secret:=}"`        // Secret
	Secure bool   `value:"${minio.secure:=true}"`    // Secure
	Bucket string `value:"${minio.bucket:=}"`
}

func init() {
	SpringBoot.RegisterNameBeanFn("minioClient", func() *minio.Client {
		var (
			config MinioConfig
		)
		SpringBoot.BindProperty("", &config)
		if config.Enable {
			client, err := minio.New(
				fmt.Sprintf("%s:%d", config.Host, config.Port),
				config.Access,
				config.Secret,
				config.Secure,
			)
			if err != nil {
				panic(err)
			}
			ok, err := client.BucketExists(config.Bucket)
			if err != nil {
				panic(err)
			}
			if !ok {
				err = client.MakeBucket(config.Bucket, "")
				if err != nil {
					panic(err)
				}
			}
			return client
		}
		return nil
	}).ConditionOnMatches(func(ctx SpringCore.SpringContext) bool {
		return SpringBoot.GetBoolProperty("minio.enable")
	})
}
