package basic

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	NextTime   int64   `json:"next_time,omitempty"`  // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64   `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`           // 返回状态描述
	VideoList  []Video `json:"video_list,omitempty"` // 视频列表
}

func Feed(c *gin.Context) {
	latest, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil {
		latest = time.Now().Unix()
	}
	videolist, err := getVideoList(latest)
	if err != nil {
		videolist = DemoVideos
	}
	c.JSON(http.StatusOK, FeedResponse{
		NextTime:   latest,
		StatusCode: 0,
		StatusMsg:  "",
		VideoList:  videolist,
	})
}

func getVideoList(latest_time int64) ([]Video, error) {
	// TODO
	return DemoVideos, nil
}
