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
	for i := uint64(0); i < 64; i++ {
		pawnAttacks[0][i] = setBit(0, i)
		pawnAttacks[0][i] = (pawnAttacks[0][i]&notAFile)>>7 ^ (pawnAttacks[0][i]&notHFile)>>9

		pawnAttacks[1][i] = setBit(0, i)
		pawnAttacks[1][i] = (pawnAttacks[1][i]&notHFile)<<7 ^ (pawnAttacks[1][i]&notAFile)<<9
	}
	return pawnAttacks
}

// maskKnightMoves generates all possible moves for knights.
func maskKnightMoves() (knightMoves [64]uint64) {
	for i := uint64(0); i < 64; i++ {
		knightMoves[i] = setBit(0, i)
		knightMoves[i] = (((knightMoves[i] >> 6) & notGHFile) ^
			((knightMoves[i] >> 10) & notABFile) ^
			((knightMoves[i] << 6) & notABFile) ^
			((knightMoves[i] << 10) & notGHFile) ^
			((knightMoves[i] >> 17) & notAFile) ^
			((knightMoves[i] >> 15) & notHFile) ^
			((knightMoves[i] << 17) & notHFile) ^
			((knightMoves[i] << 15) & notAFile))
	}
	return knightMoves
}

// maskKingMoves generates all possible moves for kings.
func maskKingMoves() (kingMoves [64]uint64) {
	for i := uint64(0); i < 64; i++ {
		kingMoves[i] = setBit(0, i)
		kingMoves[i] = (((kingMoves[i] >> 1) & notAFile) ^
			((kingMoves[i] >> 9) & notAFile) ^
			((kingMoves[i] << 7) & notAFile) ^
			(kingMoves[i] >> 8) ^
			(kingMoves[i] << 8) ^
			((kingMoves[i] >> 7) & notHFile) ^
			((kingMoves[i] << 1) & notHFile) ^
			((kingMoves[i] << 9) & notHFile))
	}
	return kingMoves
}

// maskRookMoves generates all relevant occupancy bits of rooks for magic bitboards.
func maskRookMoves() (rookMoves [64]uint64) {
	var rank, file int

	for i := 0; i < 64; i++ {
		rank = i / 8
		file = i % 8

		for r := rank + 1; r <= 6; r++ {
			rookMoves[i] |= (1 << (r*8 + file))
		}
		for r := rank - 1; r >= 1; r-- {
			rookMoves[i] |= (1 << (r*8 + file))
		}
		for f := file + 1; f <= 6; f++ {
			rookMoves[i] |= (1 << (rank*8 + f))
		}
		for f := file - 1; f >= 1; f-- {
			rookMoves[i] |= (1 << (rank*8 + f))
		}
	}

	return rookMoves
}

// maskBishopMoves generates all relevant occupancy bits of bishops for magic bitboards.
func maskBishopMoves() (bishopMoves [64]uint64) {
	var rank, file int

	for i := 0; i < 64; i++ {
		rank = i / 8
		file = i % 8

		for r, f := rank+1, file+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
			bishopMoves[i] |= (1 << (r*8 + f))
		}
		for r, f := rank-1, file+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
			bishopMoves[i] |= (1 << (r*8 + f))
		}
		for r, f := rank+1, file-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
			bishopMoves[i] |= (1 << (r*8 + f))
		}
		for r, f := rank-1, file-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
			bishopMoves[i] |= (1 << (r*8 + f))
		}
	}
	return bishopMoves
}

// generateRookMovesOnTheFly generates the rook moves for a certain blockboard (position).
func generateRookMovesOnTheFly(square int, blockboard uint64) (rookMoves uint64) {
	var rank, file int

	rank = square / 8
	file = square % 8

	for r := rank + 1; r <= 7; r++ {
		if (1<<(r*8+file))&blockboard != 0 {
			break
		} else {
			rookMoves |= (1 << (r*8 + file))
		}
	}
	for r := rank - 1; r >= 0; r-- {
		if (1<<(r*8+file))&blockboard != 0 {
			break
		} else {
			rookMoves |= (1 << (r*8 + file))
		}
	}
	for f := file + 1; f <= 7; f++ {
		if (1<<(rank*8+f))&blockboard != 0 {
			break
		} else {
			rookMoves |= (1 << (rank*8 + f))
		}
	}
	for f := file - 1; f >= 0; f-- {
		if (1<<(rank*8+f))&blockboard != 0 {
			break
		} else {
			rookMoves |= (1 << (rank*8 + f))
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
		if (1<<(r*8+f))&blockboard != 0 {
			break
		} else {
			bishopMoves |= (1 << (r*8 + f))
		}
	}
	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		if (1<<(r*8+f))&blockboard != 0 {
			break
		} else {
			bishopMoves |= (1 << (r*8 + f))
		}
	}
	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		if (1<<(r*8+f))&blockboard != 0 {
			break
		} else {
			bishopMoves |= (1 << (r*8 + f))
		}
	}
	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		if (1<<(r*8+f))&blockboard != 0 {
			break
		} else {
			bishopMoves |= (1 << (r*8 + f))
		}
	}
	return bishopMoves
}
