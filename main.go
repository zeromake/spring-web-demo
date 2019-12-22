package main

import (
	_ "github.com/go-spring/go-spring-boot-starter/starter-gin"
	_ "github.com/go-spring/go-spring-boot-starter/starter-web"
	SpringBoot "github.com/go-spring/go-spring/spring-boot"
	_ "github.com/zeromake/spring-web-demo/controllers"
)

func main() {
	SpringBoot.RunApplication("config/")
}
