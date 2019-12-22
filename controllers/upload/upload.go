package upload

import (
	"github.com/gin-gonic/gin"
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"github.com/zeromake/spring-web-demo/services/file"
	"net/http"
	"path"
)

type Controller struct {
	File *file.Service `autowire:""`
}

func init() {
	SpringBoot.RegisterBean(new(Controller)).InitFunc(func(c *Controller) {
		SpringBoot.PostMapping("/upload", c.Upload)
	})
}

func (c *Controller) Upload(ctx SpringWeb.WebContext) {
	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "没有读取到 file",
			"error":   err.Error(),
		})
		return
	}
	w, err := f.Open()
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
	out := path.Join("temp", f.Filename)

	if !c.File.ExistsObject(out) {
		err = c.File.PutObject(out, w, f.Size)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    1,
				"message": "保存失败",
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
