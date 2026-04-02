package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type ImageInfo struct {
	Title   string
	Time    string
	Content string
}

type Security struct {
	APIKEY string `json:"api"`
	User   struct {
		Id  string `json:"id"`
		Pwd string `json:"pwd"`
	} `json:"user"`
	Url      string `json:"url"`
	Fontroot string `json:"fontroot"`
}

var secu Security

func LoadJson() {
	// 1. 파일 경로 설정 (이전에 배운 절대 경로 방식)
	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Join(filepath.Dir(filename), "../../security.json")

	// 2. 파일 읽기
	fileData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("파일을 읽을 수 없습니다: %v\n", err)
		return
	}

	// 3. 구조체 객체 생성 및 데이터 채우기
	err = json.Unmarshal(fileData, &secu)
	if err != nil {
		fmt.Printf("JSON 파싱 실패: %v\n", err)
		return
	}
}

// 메인 함수
func main() {
	LoadJson()

	title, time, content := exportTTC()

	imageinfo := ImageInfo{Title: title, Time: time, Content: content}
	fmt.Println(imageinfo)

	// 이미지 생성
	generateImage(imageinfo)

	// 인스타그램 포스트
	upload(imageinfo)
}
