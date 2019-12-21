package upload

import (
	"github.com/gin-gonic/gin"
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/zeromake/spring-web-demo/types"
	"net/http"
	"path"
)

type Controller struct {
	File types.FileProvider `autowire:""`
}

func init() {
	SpringBoot.RegisterBean(new(Controller))
}

func (c *Controller) InitWebBean(wc SpringWeb.WebContainer) {
	wc.POST("/upload", c.Upload)
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
	size := file.Size
	out := path.Join("temp", file.Filename)
	if !c.File.ExistObject(out) {
		err = c.File.PutObject(out, w, size)
	}
	defer func() {
		if err != nil {
			_ = c.File.RemoveObject(out)
		}
	}()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "文件保存失败",
			"error":   err.Error(),
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
