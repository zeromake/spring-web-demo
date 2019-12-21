# go-spring web service demo

## 一、目标

1. 尝试使用 `go-spring` 做到 `controller` 注入 `service` 的引用。
2. 做到导入 `controller` 即可自动注册路由，初始化 `controller` 上的 `service` 实例。
3. 尝试两个不同文件访问逻辑实现通过配置切换做到实例化不同文件存储到对应的 `controller`。


## 二、实现方案

controller 注入在 `types` 包的 `FileProvider` 接口
```go
package upload

import (
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	"github.com/zeromake/spring-web-demo/types"
)

type Controller struct {
	File types.FileProvider `autowire:""`
}

func (c *Controller) InitWebBean(wc SpringWeb.WebContainer) {
	wc.POST("/upload", c.Upload)
}

func (c *Controller) Upload(ctx SpringWeb.WebContext) {
    // ……
}
```

types 声明需要的文件操作接口类型
```go
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
```

module 部分为通过配置实现对应的文件操作接口对象，自行到代码中查看

main 把以上的 controller 和 module 导入即可自动 ioc。
```go
package main

import (
	SpringBoot "github.com/go-spring/go-spring/spring-boot"

	_ "github.com/go-spring/go-spring-boot-starter/starter-gin"
	_ "github.com/go-spring/go-spring-boot-starter/starter-web"
	_ "github.com/zeromake/spring-web-demo/controllers"
	_ "github.com/zeromake/spring-web-demo/modules"
)

func main() {
	SpringBoot.RunApplication("config/")
}
```

## 二、运行实例验证

```bash
# 启动一个 minio 服务作为远端文件支持
docker-compose up -d minio
```
修改 `config/application.toml` 里的 `[minio]` 下的 `enable` 来切换文件的存储能力。

demo 里默认使用 `minio` 存储文件。

运行 `go run main.go` 启动服务

通过 `curl` 模拟上传

```bash
curl -F "file=@./README.md" http://127.0.0.1:8080/upload
```

可以通过 [http://127.0.0.1:9000](http://127.0.0.1:9000) 登录查看是否把文件存到了 `minio`。

修改 `config/application.toml` 里的 `[minio]` 下的 `enable = false` 重新运行再次上传会发现存到了本地
