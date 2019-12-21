package file

import (
	"fmt"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/minio/minio-go/v6"
	localProvider "github.com/zeromake/spring-web-demo/modules/file/local"
	minioProvider "github.com/zeromake/spring-web-demo/modules/file/minio"
	"github.com/zeromake/spring-web-demo/types"
)

type ProviderConfig struct {
	Enable bool   `value:"${minio.enable:=true}"`      // 是否启用 HTTP
	Host   string `value:"${minio.host:='127.0.0.1'}"` // HTTP host
	Port   int    `value:"${minio.port:=9000}"`        // HTTP 端口
	Access string `value:"${minio.access:=}"`          // Access
	Secret string `value:"${minio.secret:=}"`          // Secret
	Secure bool   `value:"${minio.secure:=true}"`      // Secure
	Bucket string `value:"${minio.bucket:=}"`
}

type ProviderStarter struct {
	types.FileProvider
	Config *ProviderConfig `autowire:""`
	//provider types.FileProvider
}

func init() {
	SpringBoot.RegisterBean(new(ProviderConfig))
	var starter types.FileProviderStarter = new(ProviderStarter)
	SpringBoot.RegisterBean(starter)
}

func (m *ProviderStarter) OnStartApplication(ctx SpringBoot.ApplicationContext) {
	if m.Config.Enable {
		client, err := minio.New(
			fmt.Sprintf("%s:%d", m.Config.Host, m.Config.Port),
			m.Config.Access,
			m.Config.Secret,
			m.Config.Secure,
		)
		if err != nil {
			panic(err)
		}
		ok, err := client.BucketExists(m.Config.Bucket)
		if err != nil {
			panic(err)
		}
		if !ok {
			err = client.MakeBucket(m.Config.Bucket, "")
			if err != nil {
				panic(err)
			}
		}
		m.FileProvider = &minioProvider.Provider{
			Client: client,
			Bucket: m.Config.Bucket,
		}
	} else {
		m.FileProvider = &localProvider.Service{}
	}
}

func (m *ProviderStarter) OnStopApplication(ctx SpringBoot.ApplicationContext) {
}
