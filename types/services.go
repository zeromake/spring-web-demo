package types

import (
	"io"
)

type FileProvider interface {
	PutObject(name string, r io.Reader, size int64) error
	ExistsObject(name string) bool
}
