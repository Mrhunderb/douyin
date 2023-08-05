package basic

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Mrhunderb/douyin/database"
	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	NextTime   int64            `json:"next_time,omitempty"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64            `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string           `json:"status_msg"`           // 返回状态描述
	VideoList  []database.Video `json:"video_list,omitempty"` // 视频列表
}

/*
视频流接口

不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
*/
func Feed(c *gin.Context) {
	latest, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	// 限制返回视频的最新投稿时间戳，不填表示当前时间
	if err != nil {
		latest = time.Now().Unix()
	}
	// TODO
	videolist := getVideoList(latest)
	c.JSON(http.StatusOK, FeedResponse{
		NextTime:   latest,
		StatusCode: 0,
		StatusMsg:  "",
		VideoList:  videolist,
	})
}

func getVideoList(latest_time int64) []database.Video {
	// TODO
	videoL, err := database.QueryVideoTime(latest_time)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return *videoL
}
