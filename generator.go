package main

var (
	notAFile  uint64 = 9187201950435737471
	notHFile  uint64 = 18374403900871474942
	notABFile uint64 = 4557430888798830399
	notGHFile uint64 = 18229723555195321596
)

// <<-------------------------------- Masks -------------------------------->>

// MaskPawnAttacks generates all possible attacks for pawns.
func maskPawnAttacks() (pawnAttacks [2][64]uint64) {
	for i := uint64(0); i < 64; i++ {
		pawnAttacks[0][i] = setBit(0, i)
		pawnAttacks[0][i] = (pawnAttacks[0][i]&notAFile)>>7 ^ (pawnAttacks[0][i]&notHFile)>>9

		pawnAttacks[1][i] = setBit(0, i)
		pawnAttacks[1][i] = (pawnAttacks[1][i]&notHFile)<<7 ^ (pawnAttacks[1][i]&notAFile)<<9
	}
	return pawnAttacks
}

// MaskKnightMoves generates all possible moves for knights.
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

// MaskKingMoves generates all possible moves for kings.
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
