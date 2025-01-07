# URL 단축 서비스

[한국어](README.ko.md) | [English](README.md)

![alt text](docs/screenshot.png)

🚀 데모: [https://f-it.kr](https://f-it.kr)

## 소개

언제나 빠르고 안정적인 처리가 가능한 새로운 방식의 URL Shortener 입니다.
또한 플랫폼별 처리나 og 태그 적용, 데이터 수집 등이 가능하도록 구성하였습니다.
보안을 위해 개인 정보는 취급하지 않지만 이메일 검증을 통해 안전하게 Short URL 을 생성하고있습니다.

### 주요 기능

| 기능 | 설명                                   |
|---------|--------------------------------------|
| 효율적인 URL 생성 | 언제나 빠르고 안전하게 URL 을 생성하는 특별한 알고리즘을 사용 |
| 딥링크 처리 | iOS/Android 딥링크를 위한 플랫폼 감지 및 리다이렉션   |
| 이메일 인증 | 계정 생성 없이 이메일 인증을 통한 간단한 URL 보안       |
| 대체 URL 지원 | 앱이 설치되지 않은 경우를 위한 대체 URL 설정 가능       |
| 웹훅 연동 | 실시간 액세스 로그를 통한 링크 사용 추적              |

### 핵심 기술

이 URL Shortener 는 여러 가지 기술적 이점을 가진 효율적인 시스템으로 구현되었습니다.

#### 키 생성 시스템

##### 순차적인 문자열 생성과 그 인덱스를 이용한 핵심 함수
조금 특별한 `id_to_key`와 `key_to_id` 함수를 사용합니다:

```
const CHARS: &str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
const BASE: i64 = 62;

/// ID To Key
/// Converts the primary key of the ShortURL table into a string.
/// Indexes when English lowercase/uppercase letters and numbers are sequentially combined.
pub fn id_to_key(mut id: i64) -> Option<String> {
    if id < 1 {
        return None;
    }
    let mut key = Vec::new();
    while id > 0 {
        id -= 1;
        let digit = (id % BASE) as usize;
        key.push(CHARS.as_bytes()[digit] as char);
        id /= BASE;
    }
    key.reverse();
    Some(key.iter().collect())
}

/// Key To ID
/// Converts arbitrary characters into a number.
/// Indexes when English lowercase/uppercase letters and numbers are sequentially combined.
pub fn key_to_id(key: &str) -> Option<i64> {
    let mut result = 0i64;
    for c in key.chars() {
        let digit = CHARS.find(c)? as i64;
        result = result * BASE + (digit + 1);
    }
    Some(result)
}

/// Split Short Key
/// Extracts the unique key in the middle and the random keys at the front and back.
pub fn split_short_key(short_key: &str) -> (String, String) {
    let front_random_key = short_key[..2].to_string();
    let back_random_key = short_key[short_key.len() - 2..].to_string();
    let random_key = &(front_random_key + &back_random_key);
    let unique_key = short_key[2..short_key.len() - 2].to_string();
    (unique_key.to_string(), random_key.to_string())
}
```

(코드 설명)
- 영어 소문자/대문자, 숫자(a-z, A-Z, 0-9)를 순차적으로 생성 했을 때 나오는 인덱스를 이용해 변환하는 방식
- 순서는 a->z, A->z, 0->9 순서로 생성하고 9까지 모두 생성되면 aa 로 자리 수를 늘려나감
- 예시:
    - 1 → "a"
    - 2 → "b"
    - 27 → "A"
    - 28 → "B"
      등
- 양방향 변환:
    - ID는 고유한 문자열로 변환
    - 문자열은 원래 ID로 역변환
    - 변환은 결정적이며 충돌이 없음

##### 단축 키 생성 과정
1. **랜덤 접두사/접미사 생성**
    - 4자리 랜덤 문자열 생성
    - 2자리씩 두 부분으로 분할

2. **핵심 키 생성**
    - 사용자가 요청한 정보와 위 4자리 랜덤한 키로 ShortURL 객체 생성
    - 생성 시 발급되는 ID 를 `id_to_key` 함수를 이용해 문자열로 변환
    - 조합: `{2자리-접두사}{변환된-id}{2자리-접미사}`

3. **키 구조 예시**
   ```
   랜덤 키가 "ABCD"이고 ID가 12345인 경우:
   - 접두사 = "AB"
   - 변환된 ID = id_to_key(12345)
   - 접미사 = "CD"
   최종 키 = "AB{변환된-id}CD"
   ```

#### 기술적 이점

1. **고유성**
    - 각 데이터베이스 ID는 고유한 키 생성
    - 충돌 검사 불필요
    - 랜덤 접두사/접미사로 보안 강화

2. **성능**
    - 일관된 키 생성 시간
    - 데이터베이스 증가에도 안정적인 성능
    - 재시도 루프 불필요

3. **확장성**
    - 대량의 URL을 효율적으로 처리
    - 실질적인 고유 키 제한 없음
    - 고트래픽 사용에 적합

4. **구현 이점**
    - 메모리 효율적
    - 최소한의 데이터베이스 쿼리
    - 쉬운 유지보수
    - 명확한 디버깅 프로세스

5. **보안**
    - 랜덤 요소를 통한 예측 불가능한 키
    - 순차적 추측 방지
    - 보안과 URL 길이의 균형

## 기술 스택

### 개발 스택
- Rust (백엔드)
- HTML (프론트엔드)

### 라이브러리

| 라이브러리 | 버전 | 용도 |
|---------|---------|---------|
| axum | 0.8.1 | 웹 프레임워크 |
| tokio | 1.0 | 비동기 런타임 |
| rusqlite | 0.32.1 | SQLite 연동 |
| lettre | 0.11 | 이메일 처리 |
| serde | 1.0 | JSON 처리 |

## API 참조

### 엔드포인트

#### 1. `GET /`
- URL 생성을 위한 웹 인터페이스
- URL 생성 요청 및 커스터마이징 UI

#### 2. `POST /v1/urls`

다음 과정으로 Short URL 생성:

1. **입력 유효성 검사**
    - URL 형식 및 접근성 확인

2. **URL 처리**
    - 플랫폼별 URL, 기본 URL 등 모든 정보를 결합하여 해시 생성
    - 이미 존재하는 해시 값인지 데이터베이스에 확인

3. **이메일 인증** *(JWT 토큰 사용 시 선택사항)*
    - JWT 토큰 제공 시:
        - 이메일 인증 생략
        - 즉시 Short URL 생성 및 반환
    - 기존 URL의 경우:
        - 인증 상태 반환
        - 대기 중인 경우 인증 재전송
    - 새로운 URL (JWT 토큰 없는 경우):
        - 미인증 URL 생성
        - 인증 이메일 전송
        - og 태그가 제공되지 않은 경우 http client 로 가져오기

#### 3. `GET /v1/verify/{code}`
인증 처리:
1. **코드 유효성 검사**
    - 코드 유효성 확인
    - 만료 상태 확인

2. **상태 업데이트**
    - 인증 상태 업데이트
    - 사용된 코드 제거
    - 인증 확인

#### 4. `GET /{short_key}`
리다이렉트 관리:
1. **캐시 확인**
    - Redis에서 URL 데이터 조회

2. **캐시 히트 처리**
    - 직접 URL 리다이렉트
    - 분석을 위한 비동기 웹훅
    - 플랫폼별 라우팅

3. **캐시 미스 처리**
    - 데이터베이스 조회
    - 캐시 업데이트
    - 사용자 리다이렉트
    - 분석 웹훅

## 시작하기

1. 저장소 클론:
   ```bash
   git clone https://github.com/lee-lou2/rust-url-shortener
   ```

2. 환경 변수 설정:
   ```env
   SERVER_PROTOCOL=http
   SERVER_HOST=127.0.0.1
   SERVER_PORT=3000
   DATABASE_URL=sqlite://sqlite3.db

   EMAIL_ADDRESS=
   EMAIL_USER_NAME=
   EMAIL_PASSWORD=
   EMAIL_HOST=
   EMAIL_PORT=
   
   JWT_SECRET=
   ```

3. 이메일 템플릿 설정:
    - `templates/verify/error.html` 구성
    - `templates/verify/failed.html` 구성
    - 기본값 `your@email.com`에서 이메일 주소 업데이트
    - `templates/index.html` 의 head 태그 수정

4. 데이터베이스 초기화:
   ```bash
   sh init_database.sh
   ```

5-1. 프로젝트 직접 실행:
   ```bash
   cargo run
   ```

5-2. 프로젝트 도커 실행:
   ```bash
   sh deploy.sh
   ```

## 향후 개발 계획

- 플랫폼별 처리 개선
- 이메일 템플릿 개선
- 관리자 대시보드
- 분석 기능
- 테스트 커버리지 확대
- Docker 지원

## 기여하기

기여를 환영합니다. 기여 방법:

1. 저장소 포크
2. 기능 브랜치 생성
3. 풀 리퀘스트 제출

## 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다.

## 지원

도움이 필요한 경우:

1. 기존 이슈 검토
2. 상세 내용과 함께 새 이슈 생성
3. 커뮤니티 토론 참여

URL 단축 서비스에 관심을 가져주셔서 감사합니다 🙇‍♂️