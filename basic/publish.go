package basic

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	user, err := database.QueryUserToken(token)
	if err != nil {
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
	finalname := fmt.Sprintf("%d_%s", user.ID, filename)
	creatFolder("./public/")
	savefile := filepath.Join("./public/", finalname)
	if err := c.SaveUploadedFile(data, savefile); err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	database.UpdateUserWorkcount(user.Name)
	url := "http://" + c.Request.Host + "/static/" + finalname
	err = database.InsertVideo(user.ID, title, url)
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  filename + " uploaded failed",
		})
	} else {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 0,
			StatusMsg:  filename + " uploaded successfully",
		})
	}
}

func creatFolder(folderPath string) {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
	}
}

type PublishRespon struct {
	StatusCode int64            `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string           `json:"status_msg"`  // 返回状态描述
	VideoList  []database.Video `json:"video_list"`  // 用户发布的视频列表
}

/*
发布列表

用户的视频发布列表，直接列出用户所有投稿过的视频
*/
func PublishList(c *gin.Context) {
	token := c.Query("token")
	id_str := c.Query("user_id")
	_, err := database.QueryUserToken(token)
	if err != nil {
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
	publish := getPublishList(token, id)
	c.JSON(http.StatusOK, PublishRespon{
		StatusCode: 0,
		StatusMsg:  "",
		VideoList:  publish,
	})
}

func getPublishList(token string, ID int64) []database.Video {
	list, err := database.QueryVideoID(token, ID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return *list
}
