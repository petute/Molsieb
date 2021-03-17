package main

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

// initPosition sets the starting position for all 8 bitboards.
func initPosition(fen string) {
	if fen != "" {
		parseFenString(fen)
	} else {
		position.pawns = 71776119061282560
		position.knights = 4755801206503243842
		position.bishops = 2594073385365405732
		position.rooks = 9295429630892703873
		position.kings = 1152921504606846992
		position.queens = 576460752303423496
		position.black = 65535
		position.white = 18446462598732840960
		position.castle = 15
		position.moveNumber = 0
		position.moveRule = 0
		position.color = true
	}
}

func initAll() {
	initLeaperAttacks()
	initSliderAttacks(true)
	initSliderAttacks(false)
	//initPosition("r3kb1r/ppp2ppp/4p3/8/1n1PNB1q/4P3/PPP2PPP/R2Q1RK1 w kq - 2 11")
	initPosition("")
}

func main() {
	initAll()
	moveList := getPseudoLegalMoves(true)
	p := makeMove(moveList[0], true, position)

	printBitboard(p.white | p.black)
	printBitboard(p.enPassant)
	printBitboard(p.pawns)

	position = p

	moveList = getPseudoLegalMoves(false)
	p = makeMove(moveList[0], false, position)

	printBitboard(p.white | p.black)
	printBitboard(p.enPassant)
	printBitboard(p.pawns)
}
