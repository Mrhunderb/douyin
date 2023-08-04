package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Mrhunderb/douyin/basic"
	_ "github.com/go-sql-driver/mysql"
)

func TestConnection() {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL database!")
	db.Close()
}

func QueryUserToken(token string) (*basic.User, error) {
	db, err := connect()

	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM users WHERE token = ?", token)
	if err != nil {
		return nil, err
	}
	return queryUser(db, row)
}

func QueryUserID(id int64) (*basic.User, error) {
	db, err := connect()

	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return queryUser(db, row)
}

func queryUser(db *sql.DB, row *sql.Row) (*basic.User, error) {
	var user basic.User
	var token string
	// 这些字段在数据库中可以为空，导致row.Scan出错
	var (
		Avatar          sql.NullString
		BackgroundImage sql.NullString
		FavoriteCount   sql.NullInt64
		Signature       sql.NullString
		TotalFavorited  sql.NullString
	)
	err := row.Scan(&user.ID, &user.Name, &token, &user.FollowCount,
		&user.FollowerCount, &user.IsFollow, &user.WorkCount,
		&Avatar, &BackgroundImage, &FavoriteCount,
		&Signature, &TotalFavorited)
	if err != nil {
		return nil, err
	}
	user.Avatar = Avatar.String
	user.BackgroundImage = BackgroundImage.String
	user.FavoriteCount = FavoriteCount.Int64
	user.Signature = Signature.String
	user.TotalFavorited = TotalFavorited.String
	return &user, nil
}

func connect() (*sql.DB, error) {
	// 构建数据库连接字符串
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// 打开数据库连接
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	// 测试连接是否成功
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
