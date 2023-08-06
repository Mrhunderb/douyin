package database

var dbUser string = "root"
var dbPass string = "886364"
var dbName string = "DOUYIN"
var dbHost string = "localhost"
var dbPort string = "3306"

var userTable = `
CREATE TABLE IF NOT EXISTS video (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	user_id BIGINT NOT NULL,
	title VARCHAR(50) NOT NULL,
	is_favorite BOOL NOT NULL DEFAULT false,
	comment_count BIGINT NOT NULL DEFAULT 0,
	favorite_count BIGINT NOT NULL DEFAULT 0,
	play_url VARCHAR(500) NOT NULL,
	cover_url VARCHAR(500),
	upload_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)
`

var videoTable = `
CREATE TABLE IF NOT EXISTS users(
	id BIGINT AUTO_INCREMENT,
	name VARCHAR(50) NOT NULL,
	PRIMARY KEY (id, name),
	token VARCHAR(100) NOT NULL,
	follow_count BIGINT NOT NULL DEFAULT 0,
	follower_count BIGINT NOT NULL DEFAULT 0,
	is_follow BOOL NOT NULL DEFAULT false,
	work_count BIGINT NOT NULL DEFAULT 0,
	avatar VARCHAR(255),
	background_count VARCHAR(255),
	favorite_count BIGINT,
	signature VARCHAR(255),
	total_favotited VARCHAR(20)
)
`
