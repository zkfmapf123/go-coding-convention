package main

import "fmt"

func main() {

	isOk := calculator()
	fmt.Println("isOk : ", isOk)

}

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
