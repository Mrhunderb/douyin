package main

import (
	"github.com/Mrhunderb/douyin/basic"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", basic.Feed)
	apiRouter.GET("/user/", basic.UserInfo)
	apiRouter.POST("/user/register/", basic.Register)
	apiRouter.POST("/user/login/", basic.Login)
	apiRouter.POST("/publish/action/", basic.Publish)
	apiRouter.GET("/publish/list/", basic.PublishList)

	// extra apis - I
	// apiRouter.POST("/favorite/action/", interact.FavoriteAction)
	// apiRouter.GET("/favorite/list/", interact.FavoriteList)
	// apiRouter.POST("/comment/action/", controller.CommentAction)
	// apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	// apiRouter.POST("/relation/action/", controller.RelationAction)
	// apiRouter.GET("/relation/follow/list/", controller.FollowList)
	// apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	// apiRouter.GET("/relation/friend/list/", controller.FriendList)
	// apiRouter.GET("/message/chat/", controller.MessageChat)
	// apiRouter.POST("/message/action/", controller.MessageAction)
}
