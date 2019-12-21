package minio

import (
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/minio/minio-go/v6"
	"github.com/pkg/errors"
	"github.com/zeromake/spring-web-demo/types"
	"io"
	"os"
	"time"
)

type Service struct {
	Client *minio.Client `autowire:""`
	Bucket string        `value:"${minio.bucket:=}"`
}

func init() {
	var s types.FileProvider = new(Service)
	SpringBoot.RegisterBean(s).ConditionOnBean("minioClient")
}

func (s *Service) PutObject(name string, r io.Reader, size int64) error {
	_, err := s.Client.PutObject(
		s.Bucket,
		name,
		r,
		size,
		minio.PutObjectOptions{},
	)
	return errors.WithStack(err)
}

func (s *Service) GetObject(name string) (types.ReadCloserAt, error) {
	r, err := s.Client.GetObject(
		s.Bucket,
		name,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r, nil
}

func (s *Service) RemoveObject(name string) error {
	return errors.WithStack(s.Client.RemoveObject(s.Bucket, name))
}

func (s *Service) StatObject(name string) (os.FileInfo, error) {
	stat, err := s.Client.StatObject(s.Bucket, name, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return ObjectInfoToStat(stat), nil
}

func (s *Service) ExistObject(name string) bool {
	_, err := s.StatObject(name)
	return err == nil
}

type stat struct {
	s *minio.ObjectInfo
}

func (s *stat) Size() int64 {
	return s.s.Size
}

func (s *stat) Name() string {
	return s.s.Key
}

func (s *stat) Mode() os.FileMode {
	return 0644
}

func (s *stat) ModTime() time.Time {
	return s.s.LastModified
}

func (s *stat) IsDir() bool {
	return false
}

func (s *stat) Sys() interface{} {
	return nil
}

func ObjectInfoToStat(s minio.ObjectInfo) os.FileInfo {
	return &stat{s: &s}
}