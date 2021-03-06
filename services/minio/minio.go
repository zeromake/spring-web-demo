package minio

import (
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/minio/minio-go/v6"
	"github.com/zeromake/spring-web-demo/types"
	"io"
)

type Service struct {
	Client *minio.Client `autowire:""`
	Bucket string        `value:"${minio.bucket:=}"`
}

func init() {
	var s = new(Service)
	SpringBoot.RegisterBean(s).AsInterface((*types.FileProvider)(nil)).ConditionOnBean("minioClient")
}

func (s *Service) PutObject(name string, r io.Reader, size int64) error {
	_, err := s.Client.PutObject(
		s.Bucket,
		name,
		r,
		size,
		minio.PutObjectOptions{},
	)
	return err
}

func (s *Service) ExistsObject(name string) bool {
	_, err := s.Client.StatObject(s.Bucket, name, minio.StatObjectOptions{})
	return err == nil
}
