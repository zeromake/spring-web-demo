package main

import (
	_ "github.com/go-spring/go-spring-boot-starter/starter-gin"
	_ "github.com/go-spring/go-spring-boot-starter/starter-web"
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"net/http"
)

func init() {
	SpringBoot.RegisterBean(new(Controller))
}

type Controller struct{}

func (c *Controller) InitWebBean(wc SpringWeb.WebContainer) {
	wc.GET("/", c.Home)
}

func (c *Controller) Home(ctx SpringWeb.WebContext) {
	ctx.String(http.StatusOK, "OK!")
}

func main() {
	SpringBoot.RunApplication("config/")
}
