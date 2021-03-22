package main

import "fmt"

type pos struct {
	pawns      uint64
	knights    uint64
	bishops    uint64
	rooks      uint64
	queens     uint64
	kings      uint64
	white      uint64
	black      uint64
	enPassant  uint64
	castle     int
	moveNumber int
	moveRule   int
	color      bool
}

var position pos

func initAll() {
	initLeaperAttacks()
	initSliderAttacks(true)
	initSliderAttacks(false)
}

func main() {
	uci("position startpos moves e4")
	printBitboard(position.pawns | position.white)
	fmt.Printf("%#v", position)
}
