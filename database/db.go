package database

import "gorm.io/gorm"

var dbUser string = "root"
var dbPass string = "886364"
var dbName string = "dou"
var dbHost string = "localhost"
var dbPort string = "3306"

var DB *gorm.DB
