package main

import "fmt"

type book struct {
	title string
	price float64
}

func (b book) print() {
	fmt.Printf("%-15s: %s\n", b.title, b.price)
}

func (b book) test() {
	fmt.Printf("%-15s: %s\n", b.title, b.price)
}

type game struct {
	title string
	price float64
}

func (g *game) print() {
	fmt.Print("%-15s: %s\n", g.title, g.price)
}

func (g *game) discount() {
	fmt.Print("discount method")
}
