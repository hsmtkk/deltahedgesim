package main

import "fmt"

func main() {
	prices := []int{100, 200, 300, 400}
	fmt.Println(prices[len(prices)-1])
	fmt.Println(prices[:len(prices)-1])
}
