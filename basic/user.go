package basic

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type UserRespon struct {
	StatusCode int64  `json:"status_code"`       // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`        // 返回状态描述
	Token      string `json:"token,omitempty"`   // 用户鉴权token
	UserID     int64  `json:"user_id,omitempty"` // 用户id
}

type InfoRespon struct {
	StatusCode int64  `json:"status_code"`    // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`     // 返回状态描述
	User       User   `json:"user,omitempty"` // 用户信息
}

/*
用户注册

新用户注册时提供用户名，密码即可，用户名需要保证唯一。创建成功后返回用户 id 和权限token
*/
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if isExist(username) {
		c.JSON(http.StatusOK, UserRespon{
			StatusCode: 1,
			StatusMsg:  "User already exist",
		})
	} else {
		id := atomic.AddInt64(&userIdSeq, 1)
		token := username + password
		userInfoList[token] = User{
			ID:            id,
			Name:          username,
			IsFollow:      false,
			FollowCount:   0,
			FollowerCount: 0,
			WorkCount:     0,
		}
		c.JSON(http.StatusOK, UserRespon{
			StatusCode: 0,
			StatusMsg:  "",
			UserID:     id,
			Token:      token,
		})
	}
}

func isExist(username string) bool {
	// TODO
	for _, user := range userInfoList {
		if strings.Compare(username, user.Name) == 0 {
			return true
		}
	}
	return false
}

/*
用户登录

通过用户名和密码进行登录，登录成功后返回用户 id 和权限 token
*/
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	if user, exist := userInfoList[token]; exist {
		c.JSON(http.StatusOK, UserRespon{
			StatusCode: 0,
			StatusMsg:  "",
			UserID:     user.ID,
			Token:      token,
		})
	} else {
		c.JSON(http.StatusOK, UserRespon{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
	}
}

/*
用户信息

获取用户的 id、昵称，还会返回关注数和粉丝数(社交部分)
*/
func UserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user, exist := userInfoList[token]; exist {
		c.JSON(http.StatusOK, InfoRespon{
			StatusCode: 0,
			StatusMsg:  "",
			User:       user,
		})
	} else {
		msg := "User " + user_id + " doesn't exist"
		c.JSON(http.StatusOK, InfoRespon{
			StatusCode: 1,
			StatusMsg:  msg,
		})
	}
}
