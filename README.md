# URL 단축 서비스 (Golang)

[한국어](README.ko.md) | [English](README.md)

🚀 **데모 사이트:** [https://url.lou2.kr](https://url.lou2.kr)

![demo site](docs/screenshot.png)

## 👋 소개

Go 언어로 개발된 현대적이고 효율적인 URL 단축 서비스입니다. 본 서비스는 단순한 URL 단축 기능을 넘어, 딥 링크 처리, 플랫폼별 리디렉션, JWT 기반 API 인증 등 다양한 고급 기능을 제공하여 사용자 경험과 보안을 모두 만족시킵니다.

## ✨ 주요 기능

| 기능                 | 설명                                                                 |
| -------------------- | -------------------------------------------------------------------- |
| **효율적인 URL 생성** | 충돌 없는 고유 키 생성 알고리즘을 통해 빠르고 안정적으로 단축 URL을 생성합니다.      |
| **딥 링크 처리** | iOS/Android 플랫폼을 감지하여 앱 딥 링크로, 앱 미설치 시에는 지정된 폴백 URL로 리디렉션합니다. |
| **JWT 인증** | 초기 웹 UI 접근 시 게스트용 JWT 토큰을 발급하며, 특정 API 엔드포인트는 JWT를 통해 인증된 요청만 허용합니다. |
| **폴백 URL 지원** | 각 플랫폼(iOS, Android) 및 기본 폴백 URL을 설정하여, 앱 미설치 등 다양한 상황에 유연하게 대응합니다. |
| **웹훅(Webhook) 연동** | 단축 URL 접근 시 지정된 URL로 실시간 알림(POST 요청)을 전송하여 데이터 수집 및 분석이 용이합니다. |
| **OG 태그 맞춤 설정** | 링크 미리보기에 사용될 OG (Open Graph) 태그 (제목, 설명, 이미지)를 직접 설정할 수 있습니다. |

## 🛠️ 핵심 기술 및 설계

### 단축키 생성 방식 ([`pkg/short_key.go`](https://www.google.com/search?q=url-shortener-go/pkg/short_key.go), [`pkg/rand.go`](https://www.google.com/search?q=url-shortener-go/pkg/rand.go) 참고)

본 서비스는 안전하고 예측 불가능한 단축키 생성을 위해 다음과 같은 방식을 사용합니다:

1.  **고유 ID 생성**: URL 정보가 데이터베이스에 저장될 때 고유한 숫자 ID (Auto Increment)를 할당받습니다.
2.  **Base62 인코딩**: 발급된 숫자 ID를 Base62로 인코딩하여 짧은 문자열로 변환합니다. (`github.com/jxskiss/base62` 라이브러리 사용)
3.  **랜덤 문자열 추가**: 2자리의 랜덤 문자열을 생성하여 Base62로 인코딩된 문자열 앞뒤에 각각 1자리씩 추가합니다.
      * 예: 랜덤 문자열 "ab", ID 인코딩 값 "cde" → 최종 단축키 "acdeb"

이 방식은 다음과 같은 장점을 가집니다:

  * **고유성 보장**: 데이터베이스 ID를 기반으로 하므로 충돌이 발생하지 않습니다.
  * **보안성 강화**: 랜덤 문자열 추가로 인해 순차적인 키 추측이 어렵습니다.
  * **일관된 성능**: 데이터베이스 크기에 관계없이 안정적인 키 생성 속도를 보장합니다.

### 시스템 아키텍처

  * **API 서버**: Fiber 웹 프레임워크를 사용하여 구현되었습니다. ([`api/server.go`](https://www.google.com/search?q=url-shortener-go/api/server.go), [`main.go`](https://www.google.com/search?q=url-shortener-go/main.go))
  * **데이터베이스**: PostgreSQL 사용을 기본으로 하며, GORM을 통해 상호작용합니다. ([`config/db.go`](https://www.google.com/search?q=url-shortener-go/config/db.go), [`model/url.go`](https://www.google.com/search?q=url-shortener-go/model/url.go))
  * **캐시**: Redis를 사용하여 자주 접근되는 URL 정보를 캐싱하여 응답 속도를 향상시킵니다. ([`config/cache.go`](https://www.google.com/search?q=url-shortener-go/config/cache.go), [`api/handler.go`](https://www.google.com/search?q=url-shortener-go/api/handler.go))
  * **설정 관리**: `.env` 파일을 통해 환경 변수를 관리합니다. ([`config/env.go`](https://www.google.com/search?q=url-shortener-go/config/env.go))
  * **오류 추적**: Sentry를 연동하여 실시간으로 오류를 모니터링하고 분석합니다. ([`main.go`](https://www.google.com/search?q=url-shortener-go/main.go))

## ⚙️ 기술 스택

  * **언어**: Go
  * **웹 프레임워크**: [Fiber](https://gofiber.io/)
  * **데이터베이스**: [PostgreSQL](https://www.postgresql.org/) (GORM 라이브러리 사용)
  * **캐시**: [Redis](https://redis.io/)
  * **JSON 처리**: [jsoniter](https://github.com/json-iterator/go)
  * **환경 변수 관리**: [godotenv](https://github.com/joho/godotenv)
  * **JWT (JSON Web Token)**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
  * **Base62 인코딩**: [jxskiss/base62](https://github.com/jxskiss/base62)
  * **유효성 검사**: [go-playground/validator](https://github.com/go-playground/validator)

## 🚀 시작하기

### 사전 준비 사항

  * Go (버전 1.18 이상 권장)
  * PostgreSQL
  * Redis
  * Git

### 설치 및 실행

1.  **저장소 복제 (Clone)**:

    ```bash
    git clone https://github.com/lee-lou2/url-shortener-go.git
    cd url-shortener-go
    ```

2.  **환경 변수 설정**:
    루트 디렉터리에 `.env` 파일을 생성하고 다음 내용을 프로젝트 환경에 맞게 수정합니다.

    ```env
    # SERVER
    SERVER_PORT=
    JWT_SECRET=
    ENCRYPT_COOKIE_KEY=
    RUN_MIGRATIONS=true  # Control with environment variables to skip automatic migrations in test environments

    # DATABASE
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=
    DB_PASSWORD=
    DB_NAME=

    # REDIS
    REDIS_HOST=localhost
    REDIS_PORT=6379
    REDIS_PASSWORD=

    # SENTRY
    SENTRY_DSN=
    ```

3.  **필요한 Go 패키지 다운로드**:

    ```bash
    go mod tidy
    ```

4.  **데이터베이스 마이그레이션**:
    애플리케이션 실행 시 GORM의 `AutoMigrate` 기능이 자동으로 `urls` 테이블을 생성하거나 업데이트합니다. ([`model/url.go`](https://www.google.com/search?q=url-shortener-go/model/url.go)의 `init` 함수 참고)

5.  **애플리케이션 실행**:

    ```bash
    go run main.go
    ```

    서버가 시작되면 기본적으로 `http://localhost:3000` 에서 접속할 수 있습니다.

## 📖 API 참조

주요 API 엔드포인트는 다음과 같습니다: ([`api/route.go`](https://www.google.com/search?q=url-shortener-go/api/route.go), [`api/handler.go`](https://www.google.com/search?q=url-shortener-go/api/handler.go) 참고)

  * **`GET /`**:

      * 설명: URL 단축 서비스를 위한 웹 UI를 제공합니다.
      * 응답: HTML 페이지 (`views/index.html`) 및 API 호출을 위한 게스트 JWT 토큰 (`{{.Token}}`).

  * **`POST /v1/urls`**:

      * 설명: 새로운 단축 URL을 생성합니다.
      * 인증: `Authorization: Bearer <token>` 헤더에 유효한 JWT 토큰 필요. ([`api/middlewares.go`](https://www.google.com/search?q=url-shortener-go/api/middlewares.go)의 `jwtAuth` 미들웨어 적용)
      * 요청 본문 (JSON): ([`api/schemas.go`](https://www.google.com/search?q=url-shortener-go/api/schemas.go)의 `createShortUrlRequest` 구조체 참고)
        ```json
        {
            "iosDeepLink": "myapp://path",
            "iosFallbackUrl": "https://apps.apple.com/app/myapp",
            "androidDeepLink": "myapp://path",
            "androidFallbackUrl": "https://play.google.com/store/apps/details?id=com.myapp",
            "defaultFallbackUrl": "https://myapp.com", // 필수
            "webhookUrl": "https://your-server.com/webhook",
            "ogTitle": "링크 미리보기 제목",
            "ogDescription": "링크 미리보기 설명",
            "ogImageUrl": "https://example.com/image.jpg"
        }
        ```
        각 필드에 대한 유효성 검사가 수행됩니다.
      * 성공 응답 (200 OK):
        ```json
        {
            "message": "URL created successfully",
            "short_key": "생성된단축키"
        }
        ```
        이미 동일한 설정의 URL이 존재하면 해당 정보를 반환합니다.
      * 실패 응답: 400 (잘못된 요청), 401 (인증 실패), 500 (서버 오류)

  * **`GET /{short_key}`**:

      * 설명: 단축 URL 키에 해당하는 원본 URL로 사용자를 리디렉션합니다.
      * 처리 과정:
        1.  Redis 캐시에서 `short_key` 조회.
        2.  캐시 미스 시 데이터베이스에서 조회.
        3.  플랫폼(iOS/Android) 감지 및 딥 링크/폴백 URL 처리. ([`views/redirect.html`](https://www.google.com/search?q=url-shortener-go/views/redirect.html)의 클라이언트 사이드 로직 참고)
        4.  웹훅 URL이 설정된 경우 비동기로 웹훅 호출. ([`model/url.go`](https://www.google.com/search?q=url-shortener-go/model/url.go)의 `SendWebHook` 메소드 참고)
        5.  최종적으로 사용자를 해당 URL로 리디렉션하거나 OG 태그가 포함된 HTML 페이지(`views/redirect.html`)를 반환합니다.
      * 실패 응답: 400 (잘못된 키 형식), 404 (URL을 찾을 수 없음)

## 🌱 기여하기 (Contributing)

이 프로젝트에 기여하고 싶으신가요? 언제든지 환영합니다\! 다음 단계를 따라주세요:

1.  이 저장소를 Fork 하세요.
2.  새로운 기능이나 버그 수정을 위한 브랜치를 생성하세요. (`git checkout -b feature/새로운기능` 또는 `bugfix/버그수정`)
3.  코드를 수정하고, 변경 사항에 대한 테스트를 추가하세요.
4.  커밋 메시지를 명확하게 작성하세요.
5.  원본 저장소로 Pull Request를 보내주세요.

## 📜 라이선스 (License)

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 `LICENSE` 파일을 참고해주세요.

## 💬 지원 및 문의

문제가 발생하거나 궁금한 점이 있다면, GitHub Issues를 통해 언제든지 문의해주세요.
