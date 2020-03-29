package home

import (
	SpringWeb "github.com/go-spring/go-spring-web/spring-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	"net/http"
)

func init() {
	SpringBoot.RegisterBean(new(Controller)).Init(func(c *Controller) {
		SpringBoot.GetMapping("/", c.Home)
	})
}

type Controller struct{}

func (c *Controller) Home(ctx SpringWeb.WebContext) {
	ctx.String(http.StatusOK, "OK!")
}
