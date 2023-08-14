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
