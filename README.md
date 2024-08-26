# Go Convention 

## 폴더구조 (선택사항)

- 서비스 / 필요에 따라 추가하거나, 제외해도 됨

```sh
## 라이브러리 / Tools
|- cmd                  ## main applications
|- internal             ## 내부에서만 사용하는 라이브러리 코드       
|- pkg                  ## 외부 App에서 사용되어도 좋은 라이브러리 
|- vendor               ## 종속성 파일 관리 go mod vendor

## Service Codes
|- api                  ## Swagger / JSON Schema / 프로토콜 정의 파일
|- web                  ## 웹 컴포넌트 / SPA 템플릿 => 서버만 사용한다면 필요없음
|- configs              ## 설정 파일 / 기본설정
|- build                ## build 결과물
|- infra                ## 인프라 파일 ( Terraform / Ansible / CDK )
|- utils                ## util 함수

Makefile                ## 빌드 / 분석에 필요한 스크립트
```

## Panic 잘 사용하기

- panic 자체는 Runtime 환경에서 위험한 놈 process.exit(1)
- panic을 발생시키는 함수라면 Prefix에 <b>Must</b> 를 붙혀야 함
- 굳이 활용한다면, recovery를 활용하는 것이 좋음
- defer는 함수 초반에 위치하게 하자...
- panic 과 fatal의 차이를 명확하게 하자
    - panic : stacktrace + stderr -> 디버깅에 유용
    - fatal : output (간결) -> 간결한 로그

```go
// 해당 방식으로 운용하면 서비스 자체는 죽지 않음
// 다만 아래와 같이 실행이 보장되지 않음
// recovery 를 하긴하지만 false를 뱉음 (return이 수행되지 않음)
func calculator() bool {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("에러는 발생했지만, 뭐 상관하지 않아...")
		}
	}()

	fmt.Println(add(10, 20)) // 정상 동작
	fmt.Println(MustMin(20, 10)) // 에러 발생 
    // ------------------- 여기서 동작이 끝남 ------------------- return false
	fmt.Println(MustMin(10, 20)) // 동작하지 않음 XXXXXX

	return true
}

func add(v1, v2 int) int {

	return v1 + v2
}

func MustMin(v1, v2 int) int {

	if v1 < v2 {
		panic("v2 bigger than v1")
	}

	return v1 - v2
}
```

```go
// 중간에 미들웨어 역할을 하는 기능을 추가
// processSafeCall
func calculator() bool {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("에러는 발생했지만, 뭐 상관하지 않아...")
		}
	}()
ll
	fmt.Println(add(10, 20))
	processSafeCall(func() {
		fmt.Println(MustMin(20, 10))
	})

	processSafeCall(func() {
		fmt.Println(MustMin(10, 20))
	})

	return true
}

func processSafeCall(fn func()) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover panic : ", r)
		}
	}()

	fn()
}

func add(v1, v2 int) int {

	return v1 + v2
}

func MustMin(v1, v2 int) int {

	if v1 < v2 {
		panic("v2 bigger than v1")
	}

	return v1 - v2
}
```

## 에러 잘 활용하기

- Error을 그냥 내도 좋지만, StackTrace를 나타내게 구성하는 것이 좋음

### 그냥 에러를 내고싶은 경우

```go
err := fmt.Errorf("%s world", "hello")
fmt.Println(err)
```

### Stack Trace를 사용하고 싶을 경우

- errors.wrap 같은 경우, 계속 쌓이면 가독성은 떨어지지만, 제일 쉽게 구현이 가능하다.

```go

// Install
// go get github.com/pkg/errors

func Test_errWithStack(t *testing.T) {

	err := func() error {
		return func() error {
			err := errors.Wrap(ErrRecordNotFound, "err-1")
			err = errors.Wrap(err, "err-2")
			err = errors.Wrap(err, "err-3")
			err = errors.Wrap(err, "err-4")
			err = errors.Wrap(err, "err-5")
			err = errors.Wrap(err, "err-6")
			return err
		}()
	}()

	fmt.Printf("%+v\n", err)
}
```

### 에러 네이밍 컨벤션

- 소문자로 진행
- 마침표 X
- 에러는 최대한 선언해서 사용하자

```go
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrUnknown        = errors.New("unknown error")
)
```

### 에러 로깅
- 핸들로 내부에서 발생한 에러는 인터셉터 혹은 미들웨어에 의해 로깅처리 진행
- <b>로깅은 최대한 다른 로직에 맡기자</b>

## Slice / Map 선언

- slice / map 에 추가될 아이템은 최대한 len, cap 을 설정하자

```go
ids := make([]string, len(users))
for i, v := range users {
	ids[i] = v.id
}
```

## Golang 에서의 파라미터 주입

- Golang에서는 옵셔널 파라미터가 없기때문에, GRPC 패턴으로 넣어줘야 함...

```go
NewAESCipher(key, WithGCM(nonce))
NewAESCipher(key, WithEncoding(euckr))
```

## Tools / Lint / Lib

- <a href="https://golangci-lint.run/"> golang-lint </a>

```sh
    ## install
    brew install golangci-lint
    brew upgrade golangci-lint

    ## run
    golangci-lint run                   ## ...
    golangci-lint run ./...             ## directories
    golangci-lint run dir1 dir2/...     ## directory
    golangci-lint run file.go           ## file

    ## .vscode setting
    .vscode/settings.json
```

- <a href="https://go.dev/doc/effective_go"> Effective Go </a>
- <a href="https://thanos.io/tip/contributing/coding-style-guide.md/#wrap-errors-for-more-context-dont-repeat-failed--there"> Thanos Golang Code Convention </a>
- <a href="https://github.com/uber-go/guide/blob/master/style.md#function-grouping-and-ordering"> 우버 코딩 가이드 </a>