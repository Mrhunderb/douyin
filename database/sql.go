package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
将数据库中指定用户的work_count字段加 1
*/
func IncUserWorkcount(token string) error {
	var user User
	result := DB.Where("token = ?", token).First(&user)
	if result.Error != nil {
		return result.Error
	}
	user.WorkCount += 1
	result = DB.Save(user)
	return result.Error
}

func IncUserFavorite(token string, n int64) error {
	var user User
	result := DB.Where("token = ?", token).First(&user)
	if result.Error != nil {
		return result.Error
	}
	user.FavoriteCount += n
	result = DB.Save(user)
	return result.Error
}

func InsertVideo(user_id int64, title, url string) error {
	result := DB.Create(&Video{
		Author:  user_id,
		Title:   title,
		PlayURL: url,
	})
	return result.Error
}

func InsertUser(name, token string) (int64, error) {
	user := User{
		Name:  name,
		Token: token,
	}
	result := DB.Create(&user)
	return int64(user.ID), result.Error
}

func IsUserExsit(name string) bool {
	if err := DB.Where("name = ?", name).First(&User{}).Error; err != nil {
		return false
	} else {
		return true
	}
}

func QueryUserName(name string) (*User, error) {
	var uesr User
	result := DB.Where("name = ?", name).First(&uesr)
	return &uesr, result.Error
}

func QueryUserToken(token string) (*User, error) {
	var user User
	result := DB.Where("token = ?", token).First(&user)
	return &user, result.Error
}

func QueryUserID(id int64) (*User, error) {
	var user User
	result := DB.First(&user, id)
	return &user, result.Error
}

func IncVideoFavorite(video_id, n int64) error {
	var video Video
	result := DB.First(&video, video_id)
	if result.Error != nil {
		return result.Error
	}
	video.FavoriteCount += n
	result = DB.Save(video)
	return result.Error
}

/*
返回数据库中所有由用户user_id上传的视频(按时间倒序)
*/
func QueryVideoID(token string, user_id int64) (*[]Video, error) {
	var videolist []Video
	result := DB.Where("author = ?", user_id).
		Order("created_at DESC").
		Find(&videolist)
	return &videolist, result.Error
}

/*
返回数据库中所有返回最新投稿时间小于last_time的视频(按时间倒序)
*/
func QueryVideoTime(last_time int64) (*[]Video, error) {
	var videolist []Video
	result := DB.Where("updated_at < ?", last_time).
		Limit(30).
		Order("created_at DESC").
		Find(&videolist)
	return &videolist, result.Error
}

func InsertFavorite(token string, video_id int64) error {
	result := DB.Create(&Favorite{
		UsrToken: token,
		VideoID:  video_id,
	})
	return result.Error
}

func DeletFavorite(token string, video_id int64) error {
	result := DB.Where("usr_token = ? AND video_id = ?", token, video_id).Delete(&Favorite{})
	return result.Error
}

func IsFavorite(token string, video_id int64) bool {
	var fav Favorite
	result := DB.Where("usr_token = ? AND video_id = ?", token, video_id).First(&fav)
	return result.Error == nil
}

func QueryFavorite(token string) (*[]Video, error) {
	var videolist []Video
	result := DB.Table("videos").
		Select("*").
		Where("id IN (?)",
			DB.Table("favorites").
				Select("video_id").
				Where("deleted_at IS NULL").
				Where("usr_token = ?", token).
				Order("created_at DESC")).
		Find(&videolist)
	return &videolist, result.Error
}

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Video{}, &User{}, &Favorite{})
}
