package local

import (
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/pkg/errors"
	"github.com/zeromake/spring-web-demo/types"
	"io"
	"os"
	"path"
)

const (
	DIR_MARK  os.FileMode = 0755
	FILE_MAEK os.FileMode = 0644
	FILE_FLAG             = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
)

type Service struct {
}

func init() {
	var s types.FileProvider = new(Service)
	SpringBoot.RegisterBean(s).ConditionOnMissingBean("minioClient")
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (s *Service) PutObject(name string, r io.Reader, size int64) error {
	dir := path.Dir(name)
	if !PathExists(dir) {
		err := os.MkdirAll(dir, DIR_MARK)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	outFile, err := os.OpenFile(name, FILE_FLAG, FILE_MAEK)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		_ = outFile.Close()
	}()
	_, err = io.Copy(outFile, r)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *Service) GetObject(name string) (types.ReadCloserAt, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return r, nil
}

func (s *Service) RemoveObject(name string) error {
	if PathExists(name) {
		return errors.WithStack(os.Remove(name))
	}
	return nil
}

func (s *Service) StatObject(name string) (os.FileInfo, error) {
	stat, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (s *Service) ExistObject(name string) bool {
	return PathExists(name)
}
