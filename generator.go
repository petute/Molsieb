package main

var (
	notAFile  uint64 = 9187201950435737471
	notHFile  uint64 = 18374403900871474942
	notABFile uint64 = 4557430888798830399
	notGHFile uint64 = 18229723555195321596
)

// <<-------------------------------- Masks -------------------------------->>

// maskPawnAttacks generates all possible attacks for pawns.
func maskPawnAttacks() (pawnAttacks [2][64]uint64) {
	for i := 0; i < 64; i++ {
		pawnAttacks[0][i] = setBit(0, i)
		pawnAttacks[0][i] = (pawnAttacks[0][i]&notAFile)>>7 ^ (pawnAttacks[0][i]&notHFile)>>9

		pawnAttacks[1][i] = setBit(0, i)
		pawnAttacks[1][i] = (pawnAttacks[1][i]&notHFile)<<7 ^ (pawnAttacks[1][i]&notAFile)<<9
	}
	return pawnAttacks
}

// maskKnightMoves generates all possible moves for knights.
func maskKnightMoves(square int) (knightMove uint64) {
	knightMove = setBit(0, square)
	knightMove = (((knightMove >> 6) & notGHFile) ^
		((knightMove >> 10) & notABFile) ^
		((knightMove << 6) & notABFile) ^
		((knightMove << 10) & notGHFile) ^
		((knightMove >> 17) & notAFile) ^
		((knightMove >> 15) & notHFile) ^
		((knightMove << 17) & notHFile) ^
		((knightMove << 15) & notAFile))

	return knightMove
}

// maskKingMoves generates all possible moves for kings.
func maskKingMoves(square int) (kingMove uint64) {
	kingMove = setBit(0, square)
	kingMove = (((kingMove >> 1) & notAFile) ^
		((kingMove >> 9) & notAFile) ^
		((kingMove << 7) & notAFile) ^
		(kingMove >> 8) ^
		(kingMove << 8) ^
		((kingMove >> 7) & notHFile) ^
		((kingMove << 1) & notHFile) ^
		((kingMove << 9) & notHFile))

	return kingMove
}

// maskRookMoves generates all relevant occupancy bits of rooks for magic bitboards.
func maskRookMoves(square int) (rookMove uint64) {
	rank := square / 8
	file := square % 8

	for r := rank + 1; r <= 6; r++ {
		rookMove |= (1 << (r*8 + file))
	}
	for r := rank - 1; r >= 1; r-- {
		rookMove |= (1 << (r*8 + file))
	}
	for f := file + 1; f <= 6; f++ {
		rookMove |= (1 << (rank*8 + f))
	}
	for f := file - 1; f >= 1; f-- {
		rookMove |= (1 << (rank*8 + f))
	}

	return rookMove
}

// maskBishopMoves generates all relevant occupancy bits of bishops for magic bitboards.
func maskBishopMoves(square int) (bishopMove uint64) {
	rank := square / 8
	file := square % 8

	for r, f := rank+1, file+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		bishopMove |= (1 << (r*8 + f))
	}
	for r, f := rank-1, file+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		bishopMove |= (1 << (r*8 + f))
	}
	for r, f := rank+1, file-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		bishopMove |= (1 << (r*8 + f))
	}
	for r, f := rank-1, file-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		bishopMove |= (1 << (r*8 + f))
	}

	return bishopMove
}

// generateRookMovesOnTheFly generates the rook moves for a certain blockboard (position).
func generateRookMovesOnTheFly(square int, blockboard uint64) (rookMoves uint64) {
	var rank, file int

	rank = square / 8
	file = square % 8

	for r := rank + 1; r <= 7; r++ {
		rookMoves |= (1 << (r*8 + file))
		if (1<<(r*8+file))&blockboard != 0 {
			break
		}
	}
	for r := rank - 1; r >= 0; r-- {
		rookMoves |= (1 << (r*8 + file))
		if (1<<(r*8+file))&blockboard != 0 {
			break
		}

	}
	for f := file + 1; f <= 7; f++ {
		rookMoves |= (1 << (rank*8 + f))
		if (1<<(rank*8+f))&blockboard != 0 {
			break
		}
	}
	for f := file - 1; f >= 0; f-- {
		rookMoves |= (1 << (rank*8 + f))
		if (1<<(rank*8+f))&blockboard != 0 {
			break
		}
	}

	return rookMoves
}

// generateBishopMovesOnTheFly generates the bishop moves for a certain blockboard (position).
func generateBishopMovesOnTheFly(square int, blockboard uint64) (bishopMoves uint64) {
	var rank, file int

	rank = square / 8
	file = square % 8

	for r, f := rank+1, file+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	return bishopMoves
}

// setOccupancy generates the relevant occupancy bitboard for a given rook or bishop moves bitboard.
func setOccupancy(bitsInMask, index int, moveMask uint64) (occupancy uint64) {
	for i := 0; i < bitsInMask; i++ {
		var square = getLS1BIndex(moveMask)
		moveMask = popBit(moveMask, square)

		if index&(1<<i) != 0 {
			occupancy |= (1 << square)
		}
	}
	return occupancy
}
