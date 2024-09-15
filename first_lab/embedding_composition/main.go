package main

import "fmt"

type Employee struct {
	Name string
	ID   int
}

type Manager struct {
	Employee
	Department string
}

func (e Employee) Work() {
	fmt.Println(e.Name, e.ID)
}

func main() {
	m := Manager{
		Employee{
			Name: "Anelya",
			ID:   1,
		},
		"Engineering",
	}
	m.Work()
}
