package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/genai"
)

var (
	// 포스트 개수 필터
	posting_cnt_filter = strings.NewReplacer("(", "", ")", " ")
)

// url주소 업데이트
func updateUrl(basic_url string) string {
	res, _ := http.Get(basic_url)
	doc, _ := goquery.NewDocumentFromResponse(res)
	posting_cnt_string := doc.Find("a.link_item").Find("span.c_cnt").Text()
	posting_cnt_int := strings.Fields(posting_cnt_filter.Replace(posting_cnt_string))
	posting_cnt := 0
	for i := 0; i < len(posting_cnt_int); i++ {
		c, _ := strconv.Atoi(posting_cnt_int[i])
		posting_cnt += c
	}
	recent_url := basic_url + strconv.Itoa(posting_cnt+2)
	return recent_url
}

// 제목, 작성시간, 태그, 내용 뽑아내는 함수(title, time, tag, content)
func exportTTC() (string, string, string) {
	recent_url := updateUrl(secu.Url) // recent_post_url로 기본 url을 보내고 반환된 url을 recent_url에 저장한다.
	res, _ := http.Get(recent_url)
	doc, _ := goquery.NewDocumentFromResponse(res)

	post := doc.Find("div#content-inner")
	header := post.Find("div#head")                             // 제목과 글 작성 시간이 들어있음.
	title := header.Find("h2")                                  // 제목 뽑아옴.
	date := header.Find("div.post-meta").Find("span.meta-date") // 시간 뽑아옴.
	//tag := post.Find("div#body").Find("div.tag_label").Find("span") // 태그 뽑아옴.

	contentRaw := post.Find("div#body").Find("div.article").Find("div.tt_article_useless_p_margin.contents_style") // 내용 뽑아옴.
	content, _ := contentRaw.Html()

	// 임의의 값이 필요할 때 사용할 것
	// var summition string = "Go 언어의 기초 조건문인 if-else와 switch-case를 설명하는 글입니다. if-else는 조건식 괄호 생략과 특정 중괄호 배치가 특징이며, switch-case는 단순 값 비교부터 조건식 활용, fallthrough 기능까지 다양한 형식을 지원합니다."

	// 제미나이가 에러나면 content에 임의의 값을 넣어줄 것
	/*
	   Go 언어의 기초 조건문인 if-else와 switch-case를 설명하는 글입니다.
	   if-else는 조건식 괄호 생략과 특정 중괄호 배치가 특징이며,
	   switch-case는 단순 값 비교부터 조건식 활용, fallthrough 기능까지 다양한 형식을 지원합니다.
	*/
	summition := contentSummation(content)

	return title.Text(), date.Text(), summition
}

// 글 내용부분을 gemini를 이용해 요약하고 반환함.
// 하루 총 20번만 가능하므로 주의할 것.
func contentSummation(content string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  secu.APIKEY,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		genai.Text("이 글은 html로 이루어져있는데 태그같은건 무시하고 내용을 100자 내외로 요약해줘. 온점과 반점 후에는 띄어쓰기를 꼭 해줘\n"+content),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	return result.Text()
}
