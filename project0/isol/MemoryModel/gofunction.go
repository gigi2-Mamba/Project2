package main

/*
User: Dpro
Date: 2025/9/10  周三
Time: 16:29
*/

var l string

func f1() {
	print(l)
}

func hello() {
	l = "hello, world"
	go f1()
}
func main() {
	hello()
}
