package main

import "fmt"

func main() {
	var input int
	fmt.Println("enter integer")
	fmt.Scanln(&input)
	if input > 0 {
		fmt.Println("positive")
	} else if input == 0 {
		fmt.Println("zero")
	} else {
		fmt.Println("negative")
	}

	for i := 1; i < 10; i++ {
		fmt.Println(i)
	}

	var number int
	fmt.Println("Введите число второй раз")
	fmt.Scanln(&number)
	switch number {
	case 1:
		fmt.Println("mon")
	case 2:
		fmt.Println("tue")
	case 3:
		fmt.Println("wed")
	case 4:
		fmt.Println("thu")
	case 5:
		fmt.Println("fri")
	case 6:
		fmt.Println("sat")
	case 7:
		fmt.Println("sun")
	}
}
