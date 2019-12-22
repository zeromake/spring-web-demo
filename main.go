package main

import (
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	_ "github.com/go-spring/go-spring/starter-gin"
	_ "github.com/go-spring/go-spring/starter-web"
	_ "github.com/zeromake/spring-web-demo/controllers"
	_ "github.com/zeromake/spring-web-demo/modules"
	_ "github.com/zeromake/spring-web-demo/services"
)

func main() {
	SpringBoot.RunApplication("config/")
}
