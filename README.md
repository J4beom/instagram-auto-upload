# 📸 instagram-auto-upload: AI 기반 이미지 생성 및 포스팅 봇

**instagram-auto-upload**는 Golang을 기반으로 블로그 콘텐츠를 수집하고, 구글의 Gemini AI를 활용한 요약 및 커스텀 이미지 생성을 거쳐 인스타그램 업로드까지 전 과정을 자동화하는 프로그램입니다.

## 🚀 주요 특징
- **스마트 크롤링**: 블로그 포스트의 핵심 제목과 본문을 자동으로 추출합니다.
- **AI 본문 요약**: Google Gemini API를 사용하여 본문을 인스타그램 최적화 문구(100자 내외)로 변환합니다.
- **자동 이미지 캔버스**: 추출된 텍스트를 바탕으로 랜덤 그라데이션 배경과 가독성 높은 텍스트 레이아웃을 생성합니다.
- **인스타그램 자동화**: `playwright-go`를 사용하여 실제 브라우저 환경에서 안정적인 포스팅을 수행합니다.

## 📂 프로젝트 구조 및 파일 설명
- **`main.go`**: 전체 워크플로우를 제어하는 메인 실행 파일입니다.
- **`text.go`**: 웹 크롤러 및 Gemini AI 인터페이스를 포함하며, 콘텐츠 요약을 담당합니다.
- **`image.go`**: 그라데이션 배경 생성 및 텍스트 렌더링 로직이 포함된 이미지 엔진입니다.
- **`upload.go`**: Playwright를 이용한 인스타그램 원격 제어 및 멀티 이미지 업로드 로직입니다.
- **`util.go`**: 폰트 설정, 색상 대비 계산 등 이미지 생성을 보조하는 유틸리티 함수 모음입니다.
- **`export/`**: 생성된 이미지가 실행 날짜별 폴더(`YYMMDD`)에 맞춰 저장되는 공간입니다.
```
.
├── export/                      # 생성된 이미지가 날짜별 폴더로 자동 저장되는 공간
└── gen-auto-image-to-instagram/ # 프로젝트 폴더
    ├── font/S-Core_Dream_OTF    # 에스코어 드림 폰트 및 라이선스가 들어있는 파일
    ├── go.mod                   # module path를 정의
    ├── go.sum                   # 디펜던시가 사용될 때마다 기존 항목과 비교하여 체크섬 추가
    └── pkg/release              # 코드가 들어있는 폴더
        ├── main.go              # 프로그램 진입점 및 전체 워크플로우 제어
        ├── text.go              # 웹 크롤링 및 Gemini AI 요약 로직
        ├── image.go             # 그라데이션 배경 및 텍스트 렌더링 (이미지 생성)
        ├── upload.go            # Playwright 기반 인스타그램 자동화 업로드
        ├── util.go              # 폰트 설정, 사이즈 계산 등 이미지 생성 보조 도구
        └── security.json        # 계정 정보 및 설정 파일 (사용자 생성 필요)
```

## 🛠️ 기술 스택 및 환경 설정 (M3 MacBook 기준)

### 필수 요구 사항
- **Language**: Go 1.20+
- **Browser Automation**: Playwright for Go
- **AI SDK**: Google Generative AI SDK
- **OS**: macOS (Apple Silicon / M3 아키텍처 최적화)

### 설치 방법(사용자마다 다를 수 있음)
1. 저장소 클론:
   ```bash
   git clone [https://github.com/your-username/InstaAuto-Go.git](https://github.com/your-username/InstaAuto-Go.git)
   cd InstaAuto-Go
   ```

2. 의존성 설치 및 Playwright 드라이버 세팅 (M3 ARM64 필수):
   ```bash
   go mod tidy
   go run [github.com/playwright-community/playwright-go/cmd/playwright@latest](https://github.com/playwright-community/playwright-go/cmd/playwright@latest) install --with-deps
   ```

3. security.json 설정: gen-auto-image-to-instagram 폴더에 아래 형식으로 계정 및 설정 파일을 생성합니다.
   ```JSON
   {
       "api" : "API Key",
       "user" : {
           "id" : "Instagram ID",
           "pwd" : "Instagram PW"
       },
       "url" : "Blog URL",
       "fontroot" : "프로젝트가 들어있는 상대경로 + /gen-auto-image-to-instagram/font/S-Core_Dream_OTF/"
   }
   ```

## 💡 주요 구현 상세
- **`UserAgent 우회`**: 인스타그램의 모바일 업로드 제한을 피하기 위해 데스크톱 UserAgent를 시뮬레이션하여 멀티 업로드를 활성화했습니다.
- **`동적 컬러 대비`**: 배경색의 밝기(Luminance)를 계산하여 글자색을 흰색/검은색으로 자동 전환함으로써 가독성을 극대화했습니다.

## ⚠️ 주의사항
- security.json 파일에는 민감한 정보가 포함되어 있으므로 절대 퍼블릭 레포지토리에 커밋하지 마세요. (.gitignore 권장)
- 인스타그램의 자동화 정책을 준수하여 적절한 시간 간격(Delay)을 두고 실행하는 것을 권장합니다.