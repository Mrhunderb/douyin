package main

import (
	"github.com/Mrhunderb/douyin/basic"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	initRouter(r)
	basic.ReadUserInfo()
	basic.ReadVideoList()
	r.Run()
}
