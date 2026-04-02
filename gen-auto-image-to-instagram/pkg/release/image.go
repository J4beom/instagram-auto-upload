package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	route = ""                          // 파일 경로 찾기용 변수
	day   = time.Now().Format("060102") // 폴더 생성용 변수

	flag                        = 0   // 제목+시간 인지 내용인지 구분하기 위한 플래그
	titlesize           float64 = 100 // 제목 폰트 사이즈
	timesize            float64 = 50  // 시간 폰트 사이즈
	contentsize         float64 = 80  // 내용 폰트 사이즈
	titletimetextheight         = 0   // 제목을 쓴 후 시간을 쓸 때 위치를 확인하기 위한 변수
)

type ColorPair struct {
	Start color.RGBA
	End   color.RGBA
}

func generateImage(imageinfo ImageInfo) {
	// 이미지를 저장할 경로 추출
	route = ""
	a, _ := os.Getwd()
	for _, x := range strings.Split(a, "/") {
		if x != "GoLang" {
			route += x
		} else if x == "GoLang" {
			route += x + "/export/"
			break
		}
		route += "/"
	}

	v := reflect.ValueOf(imageinfo)
	t := reflect.TypeOf(imageinfo)

	// 배경이미지 생성
	img1 := image.NewRGBA(image.Rect(0, 0, 1080, 1080))
	img2 := image.NewRGBA(image.Rect(0, 0, 1080, 1080))

	// 5가지의 팔레트 생성
	palette := []ColorPair{
		{color.RGBA{131, 58, 180, 255}, color.RGBA{253, 29, 29, 255}},    // 인스타 시그니처
		{color.RGBA{236, 233, 230, 255}, color.RGBA{255, 255, 255, 255}}, // 모던 미니멀
		{color.RGBA{43, 192, 228, 255}, color.RGBA{234, 236, 198, 255}},  // 오션 딥
		{color.RGBA{65, 41, 90, 255}, color.RGBA{47, 7, 67, 255}},        // 미드나잇 퍼플
		{color.RGBA{253, 200, 48, 255}, color.RGBA{243, 115, 53, 255}},   // 선셋 오렌지
	}

	// 랜덤을 이용해 사용할 팔레트 정함
	randomIndex := rand.Intn(5)
	selected := palette[randomIndex]

	// 팔레트에 저장된 색상으로 이미지 배경색을 색칠함
	bounds := img1.Bounds()
	h := bounds.Dx()
	for x := 0; x < h; x++ {
		// 0.0에서 1.0 사이의 비율 계산
		ratio := float64(x) / float64(h)

		// 각 채널별로 색상 보간(Interpolation)
		currColor := color.RGBA{
			R: uint8(float64(selected.Start.R)*(1-ratio) + float64(selected.End.R)*ratio),
			G: uint8(float64(selected.Start.G)*(1-ratio) + float64(selected.End.G)*ratio),
			B: uint8(float64(selected.Start.B)*(1-ratio) + float64(selected.End.B)*ratio),
			A: 255,
		}

		// 해당 행(x) 전체에 현재 색상을 채움
		draw.Draw(img1, image.Rect(x, bounds.Min.Y, x+1, bounds.Max.Y), &image.Uniform{currColor}, image.Point{}, draw.Src) // 해당 행(y) 전체에 현재 색상을 채움
		draw.Draw(img2, image.Rect(x, bounds.Min.Y, x+1, bounds.Max.Y), &image.Uniform{currColor}, image.Point{}, draw.Src) // 해당 행(y) 전체에 현재 색상을 채움
	}

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name // 필드 이름 (Title, Time 등)
		//fieldValue := v.Field(i).Interface() // 필드에 담긴 값

		if fieldName == "Title" {
			// 폰트 로드
			fontBytes, err := loadFont(fieldName)
			if err != nil {
				log.Fatal(err)
			}
			useFont, err := opentype.Parse(fontBytes)
			if err != nil {
				log.Fatal(err)
			}
			// 폰트 사이즈 및 해상도와 기울기등 특수효과
			face, err := opentype.NewFace(useFont, &opentype.FaceOptions{
				Size:    titlesize,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				log.Fatal(err)
			}
			// img, face, text, 라인높이(글씨크기), 배경색인덱스 순서로 들어감
			drawTitleTimeText(img1, face, imageinfo.Title, int(titlesize), randomIndex)

		} else if fieldName == "Time" {
			// 폰트 로드
			fontBytes, err := loadFont(fieldName)
			if err != nil {
				log.Fatal(err)
			}
			useFont, err := opentype.Parse(fontBytes)
			if err != nil {
				log.Fatal(err)
			}
			// 폰트 사이즈 및 해상도와 기울기등 특수효과
			face, err := opentype.NewFace(useFont, &opentype.FaceOptions{
				Size:    timesize,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				log.Fatal(err)
			}
			// img, face, text, 라인높이(글씨크기), 배경색인덱스 순서로 들어감
			drawTitleTimeText(img1, face, imageinfo.Time, int(timesize), randomIndex)

		} else if fieldName == "Content" {
			// 폰트 로드
			fontBytes, err := loadFont(fieldName)
			if err != nil {
				log.Fatal(err)
			}
			useFont, err := opentype.Parse(fontBytes)
			if err != nil {
				log.Fatal(err)
			}
			// 폰트 사이즈 및 해상도와 기울기등 특수효과
			face, err := opentype.NewFace(useFont, &opentype.FaceOptions{
				Size:    contentsize,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				log.Fatal(err)
			}
			// img, face, text, 라인높이(글씨크기), 배경색인덱스 순서로 들어감
			drawContentText(img2, face, imageinfo.Content, int(contentsize)+10, randomIndex)
		}
	}

	if err := os.MkdirAll("../../../export/"+day, 0777); err != nil {
		log.Fatal(err)
	}

	// 파일 출력
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name // 필드 이름 (Title, Time 등)
		if (fieldName == "Title" || fieldName == "Tag") && flag == 0 {
			outFile, err := os.Create(route + "/" + day + "/" + fieldName + "_output_wrapped.jpg")
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()
			png.Encode(outFile, img1)
			flag = 1
			log.Println("제목 및 태그 이미지파일 생성 완료")
		} else if fieldName == "Content" && flag == 1 {
			outFile, err := os.Create(route + "/" + day + "/" + fieldName + "_output_wrapped.jpg")
			if err != nil {
				log.Fatal(err)
			}
			defer outFile.Close()
			png.Encode(outFile, img2)
			flag = 0
			log.Println("내용 이미지파일 생성 완료")
		}
	}
}

// 텍스트 줄바꿈 & 출력
// 제목 출력(제목 + 시간)
func drawTitleTimeText(img *image.RGBA, face font.Face, text string, lineHeight int, fontflag int) {
	// 크기는 1080 x 1080 사이즈
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
	}

	//
	if fontflag == 0 || fontflag == 3 {
		d.Src = image.NewUniform(color.White)
	}

	width := d.MeasureString(text).Ceil()
	words := strings.Fields(text) // 공백 단위로 나누기
	line := ""

	if width <= 1000 {
		if titletimetextheight != 0 {
			d.Dot = fixed.P((1080-width)/2, titletimetextheight)
			d.DrawString(text)
		} else {
			d.Dot = fixed.P((1080-width)/2, 550)
			d.DrawString(text)
		}
	} else if width > 1000 {
		var apply = []string{}
		for _, w := range words {
			testLine := line
			if testLine != "" {
				testLine += " "
			}
			testLine += w

			tladv := d.MeasureString(testLine).Ceil()

			if tladv >= 1000 {
				apply = append(apply, line)
				line = w
			} else {
				line = testLine
			}

			if words[len(words)-1] == w {
				apply = append(apply, line)
			}
		}

		// 구한 텍스트 구분을 토대로 그림을 그림
		// apply가 몇 줄인지 확인하고 y축을 미리 구함.
		y := 550 - 50*len(apply)
		for i, printLine := range apply {
			x := d.MeasureString(printLine).Ceil()
			curX := (1080 - x) / 2
			curY := y + lineHeight*i
			d.Dot = fixed.P(curX, curY)
			d.DrawString(printLine)
			titletimetextheight = curY + lineHeight // 제목을 다 작성하고 나서 시간을 그릴 때 사용함
		}
	}
}

// 텍스트 줄바꿈 & 출력
// 내용 출력
func drawContentText(img *image.RGBA, face font.Face, text string, lineHeight int, fontflag int) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
	}

	if fontflag == 0 || fontflag == 3 {
		d.Src = image.NewUniform(color.White)
	}

	width := d.MeasureString(text).Ceil()
	linecount := math.Ceil(float64(width) / 1000.0) // 긴 글이 몇줄 나와야하는지 올림으로 카운팅

	curX, curY := 40, 550-50*int(linecount) // 제목인 경우

	if width > 1000 {
		words := strings.Fields(text) // 공백 단위로 나누기
		line := ""

		for _, w := range words {
			testLine := line
			if testLine != "" {
				testLine += " "
			}
			testLine += w

			// testLine의 길이 측정
			tladv := d.MeasureString(testLine).Ceil()

			if tladv > 1000 { // 가로폭 초과하면 줄바꿈
				d.Dot = fixed.P(curX, curY)
				d.DrawString(line)
				line = w
				curY += lineHeight
			} else if string(testLine[len(testLine)-1]) == "." { // 가로폭도 괜찮고, 맨 마지막에 온점이 온다면
				d.Dot = fixed.P(curX, curY)
				d.DrawString(testLine)
				line = ""
				curY += lineHeight
			} else {
				line = testLine
			}
		}
		// 남아있는 마지막 줄 그리기
		if line != "" {
			d.Dot = fixed.P(curX, curY)
			d.DrawString(line)
		}
	} else {
		// 텍스트가 maxWidth보다 짧은 경우 처리 (한 줄 출력)
		d.Dot = fixed.P(curX, curY)
		d.DrawString(text)
	}
}
