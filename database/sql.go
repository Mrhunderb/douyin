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
	CreateTable(userTable)
	CreateTable(videoTable)
	CreateTable(favoriteTable)
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
	update := `
	UPDATE users 
	SET work_count = work_count + 1 
	WHERE name = ?
	`
	_, err = db.Exec(update, username)
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
	insert := `
	INSERT INTO video (user_id, title, play_url) 
	VALUE (?, ?, ?)
	`
	_, err = db.Exec(insert, user_id, title, url)
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

func DecVideoFavorite(video_id int64) error {
	db, err := connect()

	if err != nil {
		return err
	}
	defer db.Close()
	update := `
	UPDATE video 
	SET favorite_count = favorite_count - 1 
	WHERE id = ?
	`
	_, err = db.Exec(update, video_id)
	if err != nil {
		return err
	}
	return nil
}

func IncVideoFavorite(video_id int64) error {
	db, err := connect()

	if err != nil {
		return err
	}
	defer db.Close()
	update := `
	UPDATE video 
	SET favorite_count = favorite_count + 1 
	WHERE id = ?
	`
	_, err = db.Exec(update, video_id)
	if err != nil {
		return err
	}
	return nil
}

/*
返回数据库中所有由用户user_id上传的视频(按时间倒序)
*/
func QueryVideoID(token string, user_id int64) (*[]Video, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var rows *sql.Rows
	query := `
	SELECT * FROM video 
	WHERE user_id = ? 
	ORDER BY upload_time DESC
	`
	rows, err = db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	return queryVideo(db, rows, token)
}

/*
返回数据库中所有返回最新投稿时间小于last_time的视频(按时间倒序)
*/
func QueryVideoTime(token string, last_time int64) (*[]Video, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	time_str := time.Unix(last_time, 0).Format("2006-01-02 15:04:05")
	if err != nil {
		return nil, err
	}
	var rows *sql.Rows
	query := `
	SELECT * FROM video 
	WHERE upload_time < ? 
	ORDER BY upload_time DESC
	LIMIT 30
	`
	rows, err = db.Query(query, time_str)
	if err != nil {
		return nil, err
	}
	return queryVideo(db, rows, token)
}

func queryVideo(db *sql.DB, rows *sql.Rows, token string) (*[]Video, error) {
	var videolist []Video
	defer rows.Close()
	for rows.Next() {
		var (
			video       Video
			author_id   int64
			cover_url   sql.NullString
			upload_time string
		)
		err := rows.Scan(&video.ID, &author_id, &video.Title,
			&video.CommentCount, &video.FavoriteCount,
			&video.PlayURL, &cover_url, &upload_time)
		if err != nil {
			return nil, err
		}
		author, _ := QueryUserID(author_id)
		video.Author = *author
		video.CoverURL = cover_url.String
		video.IsFavorite = IsFavorite(token, video.ID)
		videolist = append(videolist, video)
	}
	return &videolist, nil
}

func InsertFavorite(token string, video_id int64) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO favorite (token, video_id) VALUES (?, ?)", token, video_id)
	if err != nil {
		return err
	}
	return nil
}

func DeletFavorite(token string, video_id int64) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM favorite WHERE video_id = ? && token = ?", video_id, token)
	if err != nil {
		return err
	}
	return nil
}

func IsFavorite(token string, video_id int64) bool {
	db, err := connect()
	if err != nil {
		return false
	}
	defer db.Close()
	query := "SELECT COUNT(*) FROM favorite WHERE token = ? && video_id = ?"
	var count int
	err = db.QueryRow(query, token, video_id).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

func QueryFavorite(token string) (*[]Video, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var rows *sql.Rows
	query := `
	SELECT * FROM video WHERE id IN (
		SELECT video_id FROM favorite WHERE token = ? ORDER BY favorite_time DESC 
	)
	`
	rows, err = db.Query(query, token)
	if err != nil {
		return nil, err
	}
	return queryVideo(db, rows, token)
}

func CreateTable(table string) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(table)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table created successfully.")
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
