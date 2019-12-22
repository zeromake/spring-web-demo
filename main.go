package main

import (
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	_ "github.com/go-spring/go-spring/starter-gin"
	_ "github.com/go-spring/go-spring/starter-web"
	"net/http"
)

func init() {
	SpringBoot.RegisterBean(new(Controller)).InitFunc(func(c *Controller) {
		SpringBoot.GetMapping("/", c.Home)
	})
}

type Controller struct{}


func (c *Controller) Home(ctx SpringWeb.WebContext) {
	ctx.String(http.StatusOK, "OK!")
}

func main() {
	SpringBoot.RunApplication("config/")
}
