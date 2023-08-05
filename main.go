package main

import (
	"github.com/Mrhunderb/douyin/database"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	database.TestConnection()
	initRouter(r)
	r.Run()
}
