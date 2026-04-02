package main

import (
	"fmt"
	"os"
)

const (
	SCDream1 = "SCDream1.otf"
	SCDream2 = "SCDream2.otf"
	SCDream3 = "SCDream3.otf"
	SCDream4 = "SCDream4.otf"
	SCDream5 = "SCDream5.otf"
	SCDream6 = "SCDream6.otf"
	SCDream7 = "SCDream7.otf"
	SCDream8 = "SCDream8.otf"
	SCDream9 = "SCDream9.otf"
)

// 작동 순서
// 1. loadFont(A)
// 2. getFontMap(A)
func getFontMap(f string) string {
	if f == "Title" {
		return SCDream7
	} else if f == "Time" {
		return SCDream5
	} else if f == "Content" {
		return SCDream3
	} else {
		return SCDream5
	}
}

func loadFont(fn string) ([]byte, error) {
	fontFile := fmt.Sprintf(secu.Fontroot+"%s", getFontMap(fn))
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	return fontBytes, err
}
