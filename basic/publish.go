package basic

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Mrhunderb/douyin/database"
	"github.com/gin-gonic/gin"
)

type Respon struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

/*
投稿接口

登录用户选择视频上传
*/
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	defer saveVideoList()
	defer saveUserInfo()
	if _, exist := userInfoList[token]; !exist {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "formfile failed",
		})
		return
	}
	filename := data.Filename
	user := userInfoList[token]
	finalname := fmt.Sprintf("%d_%s", user.ID, filename)
	savefile := filepath.Join("./public/", finalname)
	if err := c.SaveUploadedFile(data, savefile); err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	user.WorkCount++
	userInfoList[token] = user
	id := atomic.AddInt64(&videoIdSeq, 1)
	videoList = append(videoList, database.Video{
		Author:        user,
		CommentCount:  0,
		FavoriteCount: 0,
		ID:            id,
		IsFavorite:    false,
		Title:         title,
		PlayURL:       "http://" + c.Request.Host + "/static/" + finalname,
	})
	c.JSON(http.StatusOK, Respon{
		StatusCode: 0,
		StatusMsg:  filename + " uploaded successfully",
	})
}

type PublishRespon struct {
	StatusCode int64            `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string          `json:"status_msg"`  // 返回状态描述
	VideoList  []database.Video `json:"video_list"`  // 用户发布的视频列表
}

/*
发布列表

用户的视频发布列表，直接列出用户所有投稿过的视频
*/
func PublishList(c *gin.Context) {
	token := c.Query("token")
	id_str := c.Query("user_id")
	if _, exist := userInfoList[token]; !exist {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}
	id, err := strconv.ParseInt(id_str, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "ID error",
		})
		return
	}
	publish := getPublishList(id)
	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: 0,
		StatusMsg:  "",
		VideoList:  publish,
	})
}

func getPublishList(ID int64) []database.Video {
	list, err := database.QueryVideo(time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return *list
}
