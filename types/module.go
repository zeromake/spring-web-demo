package types

import (
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"io"
	"os"
)

type ReadCloserAt interface {
	io.Reader
	io.Closer
	io.ReaderAt
}

type FileProvider interface {
	PutObject(name string, r io.Reader, size int64) error
	GetObject(name string) (ReadCloserAt, error)
	RemoveObject(name string) error
	StatObject(name string) (os.FileInfo, error)
	ExistObject(name string) bool
}
type FileProviderStarter interface {
	FileProvider
	OnStartApplication(ctx SpringBoot.ApplicationContext)
	OnStopApplication(ctx SpringBoot.ApplicationContext)
}
