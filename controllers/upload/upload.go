package upload

import (
	"github.com/gin-gonic/gin"
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"io"
	"net/http"
	"os"
	"path"
)

type Controller struct{}

const (
	DIR_MARK  os.FileMode = 0755
	FILE_MAEK os.FileMode = 0644
	FILE_FLAG             = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
)

func init() {
	SpringBoot.RegisterBean(new(Controller)).Init(func(c *Controller) {
		SpringBoot.PostMapping("/upload", c.Upload)
	})
}

func (c *Controller) Upload(ctx SpringWeb.WebContext) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "没有读取到 file",
			"error":   err.Error(),
		})
		return
	}
	w, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "文件打开失败",
			"error":   err.Error(),
		})
		return
	}
	defer func() {
		_ = w.Close()
	}()
	out := path.Join("temp", file.Filename)
	if !PathExists(out) {
		dir := path.Dir(out)
		if !PathExists(dir) {
			err = os.MkdirAll(dir, DIR_MARK)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code":    1,
					"message": "文件夹创建失败",
					"error":   err.Error(),
				})
				return
			}
		}
		dst, err := os.OpenFile(out, FILE_FLAG, FILE_MAEK)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    1,
				"message": "文件创建失败",
				"error":   err.Error(),
			})
			return
		}
		defer func() {
			_ = dst.Close()
		}()
		_, err = io.Copy(dst, w)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    1,
				"message": "文件写入失败",
				"error":   err.Error(),
			})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "该文件已存在",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": http.StatusText(http.StatusOK),
		"data": map[string]string{
			"url": out,
		},
	})
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
