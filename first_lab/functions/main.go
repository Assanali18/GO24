package main

import "fmt"

func add(first, second int) int {
	return first + second
}
func swap(first, second string) (string, string) {
	return second, first
}
func reminder_and_quotient(first, second int) (int, int) {
	return first % second, first / second
}

func main() {
	fmt.Println(add(1, 2))
	fmt.Println(swap("hello", "world"))
	fmt.Println(reminder_and_quotient(5, 2))
}
