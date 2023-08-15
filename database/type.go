package database

import (
	"gorm.io/gorm"
)

type Video struct {
	ID            uint           `gorm:"primarykey;autoIncrement"`
	CreatedAt     int64          `gorm:"autoCreateTime"`
	UpdatedAt     int64          `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Author        int64          `gorm:"foreignKey:UserID"` // 视频作者信息
	CommentCount  int64          `gorm:"default:0"`         // 视频的评论总数
	CoverURL      string         // 视频封面地址
	FavoriteCount int64          `gorm:"default:0"` // 视频的点赞总数
	PlayURL       string         `gorm:"not null"`  // 视频播放地址
	Title         string         `gorm:"not null"`  // 视频标题
}

type User struct {
	gorm.Model
	Name            string `gorm:"not null"` // 用户名称
	Token           string `gorm:"not null"`
	FollowCount     int64  `gorm:"default:0"` // 关注总数
	FollowerCount   int64  `gorm:"default:0"` // 粉丝总数
	FavoriteCount   int64  `gorm:"default:0"` // 喜欢数
	WorkCount       int64  `gorm:"default:0"` // 作品数
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	Signature       string // 个人简介
	TotalFavorited  string // 获赞数量
}

type Favorite struct {
	gorm.Model
	UsrToken string `gorm:"foreignKey:Token"`
	VideoID  int64  `gorm:"foreignKey:VideoID"`
}

type VideoJSON struct {
	Author        UserJSON `json:"author,omitempty"`      // 视频作者信息
	CommentCount  int64    `json:"comment_count"`         // 视频的评论总数
	CoverURL      string   `json:"cover_url"`             // 视频封面地址
	FavoriteCount int64    `json:"favorite_count"`        // 视频的点赞总数
	ID            int64    `json:"id"`                    // 视频唯一标识
	IsFavorite    bool     `json:"is_favorite,omitempty"` // true-已点赞，false-未点赞
	PlayURL       string   `json:"play_url"`              // 视频播放地址
	Title         string   `json:"title"`                 // 视频标题
}

func ConvertVideo(video *Video, token string) *VideoJSON {
	user, _ := QueryUserID(video.Author)
	isfav := IsFavorite(token, int64(video.ID))
	return &VideoJSON{
		Author:        *ConvertUser(user),
		CommentCount:  video.CommentCount,
		CoverURL:      video.CoverURL,
		FavoriteCount: video.FavoriteCount,
		ID:            int64(video.ID),
		IsFavorite:    isfav,
		PlayURL:       video.PlayURL,
		Title:         video.Title,
	}
}

// 视频作者信息
//
// User
type UserJSON struct {
	Avatar          string `json:"avatar,omitempty"`           // 用户头像
	BackgroundImage string `json:"background_image,omitempty"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`             // 喜欢数
	FollowCount     int64  `json:"follow_count"`               // 关注总数
	FollowerCount   int64  `json:"follower_count"`             // 粉丝总数
	ID              int64  `json:"id"`                         // 用户id
	IsFollow        bool   `json:"is_follow"`                  // true-已关注，false-未关注
	Name            string `json:"name"`                       // 用户名称
	Signature       string `json:"signature,omitempty"`        // 个人简介
	TotalFavorited  string `json:"total_favorited,omitempty"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`                 // 作品数
}

func ConvertUser(user *User) *UserJSON {
	return &UserJSON{
		ID:              int64(user.ID),
		Avatar:          "",
		BackgroundImage: "",
		FavoriteCount:   user.FavoriteCount,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        false,
		Name:            user.Name,
		Signature:       "",
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
	}
}
