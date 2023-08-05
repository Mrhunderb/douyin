package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

/*
将数据库中指定用户的work_count字段加1
*/
func UpdateUserWorkcount(username string) error {
	db, err := connect()

	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE users SET work_count = work_count + 1 WHERE name = ?", username)
	if err != nil {
		return err
	}
	return nil
}

func InsertVideo(user_id int64, title, url string) error {
	db, err := connect()

	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO video (user_id, title, play_url) VALUE (?, ?, ?)", user_id, title, url)
	if err != nil {
		return err
	}
	return nil
}

func InsertUser(name, token string) (*User, error) {
	db, err := connect()

	if err != nil {
		return nil, err
	}
	defer db.Close()
	result, err := db.Exec("INSERT INTO users (name, token) VALUE (?, ?)", name, token)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return QueryUserName(name)
}

func QueryUserName(name string) (*User, error) {
	db, err := connect()

	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("SELECT * FROM users WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	return queryUser(db, row)
}

func QueryUserToken(token string) (*User, error) {
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

func QueryUserID(id int64) (*User, error) {
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

func queryUser(db *sql.DB, row *sql.Row) (*User, error) {
	var user User
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

func QueryVideo(last_time int64) (*[]Video, error) {
	db, err := connect()
	var videolist []Video
	if err != nil {
		return nil, err
	}
	defer db.Close()
	time_str := time.Unix(last_time, 0).Format("2006-01-02 15:04:05")
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if len(time_str) > 19 {
		rows, err = db.Query("SELECT * FROM video WHERE upload_time < ? ORDER BY upload_time DESC", time_str[1:])
	} else {
		rows, err = db.Query("SELECT * FROM video WHERE upload_time < ? ORDER BY upload_time DESC", time_str)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			video       Video
			author_id   int64
			cover_url   sql.NullString
			upload_time string
		)
		err = rows.Scan(&video.ID, &author_id, &video.Title,
			&video.IsFavorite, &video.CommentCount, &video.FavoriteCount,
			&video.PlayURL, &cover_url, &upload_time)
		if err != nil {
			return nil, err
		}
		author, _ := QueryUserID(author_id)
		video.Author = *author
		video.CoverURL = cover_url.String
		videolist = append(videolist, video)
	}
	return &videolist, nil
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
