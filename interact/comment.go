package interact

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Mrhunderb/douyin/database"
	"github.com/gin-gonic/gin"
)

type CommentRespon struct {
	StatusCode int64            `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string           `json:"status_msg"`  // 返回状态描述
	Comment    database.Comment `json:"comment"`     // 评论内容
}

type CommentListRespon struct {
	StatusCode  int64              `json:"status_code"`  // 状态码，0-成功，其他值-失败
	StatusMsg   string             `json:"status_msg"`   // 返回状态描述
	CommentList []database.Comment `json:"comment_list"` // 评论列表
}

func ReturnErrMsg(c *gin.Context) {
	c.JSON(http.StatusOK, Respon{
		StatusCode: 1,
		StatusMsg:  "Error occurred",
	})
}

func CommentAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	id, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "Failed to convert video id",
		})
		return
	}
	_, err = database.QueryUserToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "user doesn't exsit",
		})
		return
	}
	action_type := c.Query("action_type")
	if strings.Compare(action_type, "1") == 0 {
		content := c.Query("comment_text")
		comment_id, err := database.InsertComment(token, id, content)
		if err != nil {
			ReturnErrMsg(c)
			return
		}
		if err := database.IncVideoComment(id); err != nil {
			ReturnErrMsg(c)
			return
		}
		user, _ := database.QueryUserToken(token)
		c.JSON(http.StatusOK, CommentRespon{
			StatusCode: 0,
			StatusMsg:  "comment successful",
			Comment: database.Comment{
				Id:         comment_id,
				User:       *user,
				Content:    content,
				CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			},
		})
	} else if strings.Compare(action_type, "2") == 0 {
		comment_id := c.Query("comment_id")
		commentID, err := strconv.ParseInt(comment_id, 10, 64)
		if err != nil {
			ReturnErrMsg(c)
			return
		}
		if err := database.DeleteComment(commentID); err != nil {
			ReturnErrMsg(c)
			return
		}
		if err := database.DecVideoComment(id); err != nil {
			ReturnErrMsg(c)
			return
		}
		c.JSON(http.StatusOK, CommentRespon{
			StatusCode: 0,
			StatusMsg:  "delete successful",
			Comment:    database.Comment{},
		})
	}
}

func CommentList(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	id, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Respon{
			StatusCode: 1,
			StatusMsg:  "Failed to convert video id",
		})
		return
	}
	commentList, err := database.QueryCommentsByVideoID(id)
	if err != nil {
		c.JSON(http.StatusOK, CommentListRespon{
			StatusCode:  1,
			StatusMsg:   "query failed",
			CommentList: nil,
		})
	} else {
		c.JSON(http.StatusOK, CommentListRespon{
			StatusCode:  1,
			StatusMsg:   "query failed",
			CommentList: *commentList,
		})
	}
	fmt.Println(token)
}
