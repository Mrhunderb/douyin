package basic

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mrhunderb/douyin/database"
	"github.com/gin-gonic/gin"
)

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	id, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "Failed to convert video id",
		})
		return
	}
	action := c.Query("action_type")
	_, err = database.QueryUserToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "user doesn't exsit",
		})
		return
	}
	if strings.Compare(action, "1") == 0 {
		err = database.InsertFavorite(token, id)
		if err != nil {
			database.IncVideoFavorite(id, 1)
			fmt.Println(err.Error())
			c.JSON(http.StatusOK, Respon{
				StatusCode: 1,
				StatusMsg:  "favorite failed",
			})
			return
		}
	} else {
		err = database.DeletFavorite(token, id)
		if err != nil {
			database.IncVideoFavorite(id, -1)
			fmt.Println(err.Error())
			c.JSON(http.StatusOK, Respon{
				StatusCode: 1,
				StatusMsg:  "cancel failed",
			})
			return
		}
	}
	c.JSON(http.StatusOK, Respon{
		StatusCode: 0,
		StatusMsg:  "",
	})
}

type FavoriteListRepson struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string  `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 用户点赞视频列表
}

func FavoriteList(c *gin.Context) {
	id_str := c.Query("user_id")
	token := c.Query("token")
	user_id, err := strconv.ParseInt(id_str, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = database.QueryUserID(user_id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, FavoriteListRepson{
			StatusCode: 1,
			StatusMsg:  "user doesn't exsit",
			VideoList:  nil,
		})
		return
	}
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, FavoriteListRepson{
			StatusCode: 1,
			StatusMsg:  "query failed",
			VideoList:  nil,
		})
	} else {
		c.JSON(http.StatusOK, FavoriteListRepson{
			StatusCode: 0,
			StatusMsg:  "",
			VideoList:  *getFavoritList(token),
		})
	}
}

func getFavoritList(token string) *[]Video {
	favorite, err := database.QueryFavorite(token)
	if err != nil {
		fmt.Println(err)
	}
	var list []Video
	for _, vidoe := range *favorite {
		list = append(list, *ConvertVideo(&vidoe, token))
	}
	return &list
}
