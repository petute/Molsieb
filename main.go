package main

var position struct {
	pawns   uint64
	knights uint64
	bishops uint64
	rooks   uint64
	queens  uint64
	kings   uint64
	white   uint64
	black   uint64
}

// initPosition sets the starting position for all 8 bitboards.
func initPosition() {
	position.pawns = 71776119061282560
	position.knights = 4755801206503243842
	position.bishops = 2594073385365405732
	position.rooks = 9295429630892703873
	position.kings = 1152921504606846992
	position.queens = 576460752303423496
	position.black = 65535
	position.white = 18446462598732840960
}

func main() {
	var kingMoves = maskKingMoves()

	printBitboard(kingMoves[60])
}
