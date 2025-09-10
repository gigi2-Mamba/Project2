package main

import "time"

var a, b int

func f() {
	a = 1
	b = 2
}

func g() {
	print(b)
	time.Sleep(20 * time.Millisecond)
	print(a)
}

func main() {
	go f() // goroutine 1
	time.Sleep(300)
	g() // goroutine 2
	// case1:   0,1   case 2: 2,1
	//case 3: 2,0     这里面有两个goroutine,实际上go没有要求一个goroutine能够看到另外一个goroutine对变量的写读，而且cpu和编译器会重排序

}
