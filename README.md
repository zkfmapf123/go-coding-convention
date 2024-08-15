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