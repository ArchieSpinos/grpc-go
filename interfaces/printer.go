package main

func main() {
	var (
		mobydick  = book{title: "moby dick", price: 10}
		minecraft = game{title: "minecraft", price: 4}
		tetris    = game{title: "tetris", price: 8}
	)
	var store list
	store = append(store, &minecraft, &tetris, &mobydick)

	store.discount()

	// store.print()

}
