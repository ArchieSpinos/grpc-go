package main

import "fmt"

type printer interface {
	print()
}

type list []printer

func (l list) print() {
	if len(l) == 0 {
		fmt.Println("Zero length list")
		return
	}

	for _, it := range l {
		it.print()
	}
}

func (l list) discount() {
	for _, item := range l {
		g, isGame := item.(*game)
		if isGame != true {
			continue
		}
		g.discount()
	}
}
