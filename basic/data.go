package basic

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
)

var videoList = []Video{}

var userInfoList = map[string]User{}

var userIdSeq int64 = 0

var videoIdSeq int64 = 0

func saveVideoList() {
	data, err := json.MarshalIndent(videoList, "", "   ")
	if err != nil {
		log.Fatal("Failed to conver json")
		return
	}
	err = os.WriteFile("table/video.json", data, fs.FileMode(os.O_RDWR|os.O_CREATE))
	if err != nil {
		log.Fatal("Failed to write file")
		return
	}
}

func ReadVideoList() {
	file, err := os.ReadFile("table/video.json")
	if err != nil {
		log.Fatal("Failed to read file")
		return
	}
	err = json.Unmarshal(file, &videoList)
	if err != nil {
		log.Fatal("Failed to unmarshal file")
		return
	}
	videoIdSeq = int64(len(videoList))
}

func saveUserInfo() {
	data, err := json.MarshalIndent(userInfoList, "", "   ")
	if err != nil {
		log.Fatal("Failed to conver json")
		return
	}
	err = os.WriteFile("table/user.json", data, fs.FileMode(os.O_RDWR|os.O_CREATE))
	if err != nil {
		log.Fatal("Failed to write file")
		return
	}
}

func ReadUserInfo() {
	file, err := os.ReadFile("table/user.json")
	if err != nil {
		log.Fatal("Failed to read file")
		return
	}
	err = json.Unmarshal(file, &userInfoList)
	if err != nil {
		log.Fatal("Failed to unmarshal file")
		return
	}
	userIdSeq = int64(len(userInfoList))
}
