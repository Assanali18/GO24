package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	Name     string
	Price    int
	Quantity int
}

func productToJSON(p Product) (string, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
func jsonToProduct(jsonStr string) (Product, error) {
	var p Product
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

func main() {
	p := Product{Name: "Laptop", Price: 1000, Quantity: 5}
	jsonStr, err := productToJSON(p)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Product as JSON:", jsonStr)
	}

	decodedProduct, err := jsonToProduct(jsonStr)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Decoded Product: %+v\n", decodedProduct)
	}

	fmt.Println("Product Name:", decodedProduct.Name)
}
