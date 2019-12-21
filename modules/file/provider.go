package file

import (
	"fmt"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/minio/minio-go/v6"
	localProvider "github.com/zeromake/spring-web-demo/modules/file/local"
	minioProvider "github.com/zeromake/spring-web-demo/modules/file/minio"
	"github.com/zeromake/spring-web-demo/types"
	"io"
	"os"
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
	Config   *ProviderConfig `autowire:""`
	provider types.FileProvider
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
		m.provider = &minioProvider.Provider{
			Client: client,
			Bucket: m.Config.Bucket,
		}
	} else {
		m.provider = &localProvider.Service{}
	}
}

func (m *ProviderStarter) OnStopApplication(ctx SpringBoot.ApplicationContext) {
}

func (m *ProviderStarter) PutObject(name string, r io.Reader, size int64) error {
	return m.provider.PutObject(name, r, size)
}

func (m *ProviderStarter) GetObject(name string) (types.ReadCloserAt, error) {
	return m.provider.GetObject(name)
}

func (m *ProviderStarter) RemoveObject(name string) error {
	return m.provider.RemoveObject(name)
}

func (m *ProviderStarter) StatObject(name string) (os.FileInfo, error) {
	return m.provider.StatObject(name)
}

func (m *ProviderStarter) ExistObject(name string) bool {
	return m.provider.ExistObject(name)
}
