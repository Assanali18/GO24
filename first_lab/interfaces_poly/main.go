package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width float64
}
type Circle struct {
	Radius float64
}

func (C Circle) Area() float64 {
	return math.Pi * C.Radius * C.Radius
}
func (R Rectangle) Area() float64 {
	return R.Width * R.Width
}

func PrintArea(s Shape) {
	fmt.Println(s.Area())
}

func main() {
	c := Circle{
		Radius: 5,
	}
	r := Rectangle{
		Width: 5,
	}
	PrintArea(c)
	PrintArea(r)
}
