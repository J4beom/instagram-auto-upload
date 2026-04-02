package main

import (
	"log"
	"time"

	"github.com/playwright-community/playwright-go"
)

// 브라우저를 띄워서 업로드 해주는 방식
// 1. 자동 로그인
// 2. 자동 사진 추가
// 3. 설명 추가 및 공유하기
func upload(imageinfo ImageInfo) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("Playwright 실행 실패: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless:       playwright.Bool(false),
		ExecutablePath: playwright.String("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
		Args: []string{
			"--disable-blink-features=AutomationControlled", // 자동화 도구 감지 차단
			"--start-maximized",
		},
	})

	if err != nil {
		log.Fatalf("브라우저 실행 실패: %v", err)
	}

	context, _ := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.7559.133 Safari/537.36"),
		Viewport: &playwright.Size{
			Width:  1200,
			Height: 800,
		},
		IsMobile: playwright.Bool(false),
		HasTouch: playwright.Bool(false),
	})

	page, _ := context.NewPage()
	page.Goto("https://www.instagram.com/accounts/login/")

	// 1. 로그인 하기
	// 아이디, 패스워드를 입력하고 로그인버튼을 누름
	idField := page.Locator("input[name='email'], input[name='username']")
	pwField := page.Locator("input[name='pass'], input[name='password']")
	loginBtn := page.Locator("div[role='none'], button[type='submit']").Filter(playwright.LocatorFilterOptions{
		HasText: "로그인",
	}).First()

	if err := idField.Fill(secu.User.Id); err != nil {
		log.Fatal(err)
	}
	if err := pwField.Fill(secu.User.Pwd); err != nil {
		log.Fatal(err)
	}

	loginBtn.Click()

	// 인스타그램 메인 URL로 이동할 때까지 기다립니다.
	page.WaitForURL("https://www.instagram.com/**")

	// 2. 게시물 추가하기
	// + 버튼 클릭 -> 게시물 버튼 클릭
	// 파일을 선택하고 다음버튼 클릭
	// 게시물 설명 적고 공유하기 클릭

	// 2-1. '만들기' 또는 '새 게시물' 단어가 포함된 라벨을 가진 SVG의 부모 링크를 찾습니다.
	plusBtn := page.Locator("div[aria-selected='false']").Filter(playwright.LocatorFilterOptions{
		Has: page.Locator("svg[aria-label*='새로운 게시물']"),
	})
	plusBtn.Click()

	// 2-2. 요소가 나타날 때까지 대기
	if err := plusBtn.WaitFor(playwright.LocatorWaitForOptions{
		State: playwright.WaitForSelectorStateVisible,
	}); err != nil {
		log.Fatal("버튼을 찾을 수 없습니다.")
	}

	// 2-3. 업로드 버튼 클릭 (모바일 뷰에서는 하단 중앙에 위치)
	// 선택자는 인스타그램 업데이트에 따라 변동될 수 있습니다.
	postBtn := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "컴퓨터에서 선택",
	})

	// 2-4. 파일 선택창 이벤트 예약
	fileChooser, err := page.ExpectFileChooser(func() error {
		// 2. 이제 '게시물' 버튼 클릭
		return postBtn.Click()
	})
	if err != nil {
		log.Fatalf("파일 선택창 대기 실패: %v", err)
	}

	// 3. 업로드할 이미지 파일 경로들을 슬라이스(리스트)로 만듭니다.
	files := []string{
		"/Users/jabeom/Documents/education/GoLang/export/" + day + "/Title_output_wrapped.jpg",   // 제목 사진
		"/Users/jabeom/Documents/education/GoLang/export/" + day + "/Content_output_wrapped.jpg", // 내용 사진
	}

	// 3-1. 코드로 파일을 직접 주입합니다. (실제 Finder 창이 뜨지 않고 바로 올라갑니다.)
	if err := fileChooser.SetFiles(files); err != nil {
		log.Fatalf("파일 업로드 실패: %v", err)
	}

	// 다음 버튼이 두 번 떠서 두 번 호출함
	nextBtn := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "다음",
	})
	nextBtn.Click()
	nextBtn.Click()

	shareBtn := page.GetByRole("button", playwright.PageGetByRoleOptions{
		Name: "공유하기",
	})

	// 4. 편집 및 캡션 입력 (Next 버튼 클릭 루틴)
	postingField := page.Locator("textarea[aria-label='문구를 입력하세요...'], div[aria-label='문구를 입력하세요...']")
	postingText := "\n" + imageinfo.Title + "\n" + imageinfo.Time + "\n\n" + imageinfo.Content

	if err := postingField.Fill(postingText); err != nil {
		log.Fatal(err)
	}
	shareBtn.Click()

	time.Sleep(20 * time.Second) // 업로드 완료 대기
	browser.Close()
	pw.Stop()
}
