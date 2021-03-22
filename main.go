package main

import (
	"fmt"
)

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

var game struct {
	wtime float64
	btime float64
	winc  float64
	binc  float64
}

var position pos

func initAll() {
	initLeaperAttacks()
	initSliderAttacks(true)
	initSliderAttacks(false)
}

func main() {
	initAll()
	uci("position startpos")
	uci("go wtime 10 btime 10 winc 5 binc 5 depth 4")
	fmt.Printf("%v %d", game, depth)
}
