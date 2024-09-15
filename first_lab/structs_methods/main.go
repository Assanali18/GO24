package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() {
	fmt.Println("Hello, my name is", p.Name)
}

func main() {
	p := Person{Name: "Anelya", Age: 22}
	p.Greet()
}
