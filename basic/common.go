package basic

import "fmt"

type Respon struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type Video struct {
	Author        User   `json:"author"`              // 视频作者信息
	CommentCount  int64  `json:"comment_count"`       // 视频的评论总数
	CoverURL      string `json:"cover_url,omitempty"` // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"`      // 视频的点赞总数
	ID            int64  `json:"id"`                  // 视频唯一标识
	IsFavorite    bool   `json:"is_favorite"`         // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url"`            // 视频播放地址
	Title         string `json:"title"`               // 视频标题
}

func (v Video) String() string {
	return fmt.Sprintf("ID: %d\nAuther: %s\nTitle: %s\nURL:%s", v.ID, v.Author.Name, v.Title, v.PlayURL)
}

// 视频作者信息
type User struct {
	Avatar          string `json:"avatar,omitempty"`           // 用户头像
	BackgroundImage string `json:"background_image,omitempty"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count,omitempty"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`               // 关注总数
	FollowerCount   int64  `json:"follower_count"`             // 粉丝总数
	ID              int64  `json:"id"`                         // 用户id
	IsFollow        bool   `json:"is_follow"`                  // true-已关注，false-未关注
	Name            string `json:"name"`                       // 用户名称
	Signature       string `json:"signature,omitempty"`        // 个人简介
	TotalFavorited  string `json:"total_favorited,omitempty"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`                 // 作品数
}

func (u User) String() string {
	return fmt.Sprintf("ID: %d\nName: %s\nWorkCount%d", u.ID, u.Name, u.WorkCount)
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}
