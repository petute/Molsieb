package main

var (
	notAFile  uint64 = 9187201950435737471
	notHFile  uint64 = 18374403900871474942
	notABFile uint64 = 4557430888798830399
	notGHFile uint64 = 18229723555195321596
)

// Amount oft relevant occupancy bits per square for the bishop.
var relevantBitsBishop = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// Amount oft relevant occupancy bits per square for the rook.
var relevantBitsRook = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// <<--------------------------------- Masks --------------------------------->>

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

// <<--------------------------------- Magic --------------------------------->>

// setOccupancy generates the relevant occupancy bitboard for a given rook or bishop moves bitboard.
func setOccupancy(bitsInMask, index int, moveMask uint64) (occupancy uint64) {
	for i := 0; i < bitsInMask; i++ {
		square := getLS1BIndex(moveMask)
		moveMask = popBit(moveMask, square)

		if index&(1<<i) != 0 {
			occupancy |= (1 << square)
		}
	}
	return occupancy
}

// state is used to generate the random numbers.
var state uint32 = 1804289383

// getRandom32BitNumber generates pseudoRandom numbers (XORSHIFT32).
func getRandom32BitNumber() uint32 {
	number := state

	number ^= number << 13
	number ^= number >> 17
	number ^= number << 5

	state = number

	return number
}

// getRandom64BitNumber generates a random 64 bit pseudo legal number. (FFFF == 65535 == 16 bits)
func getRandom64BitNumber() uint64 {
	n1 := uint64(getRandom32BitNumber()) & 0xFFFF
	n2 := uint64(getRandom32BitNumber()) & 0xFFFF
	n3 := uint64(getRandom32BitNumber()) & 0xFFFF
	n4 := uint64(getRandom32BitNumber()) & 0xFFFF

	return n1 | (n2 << 16) | (n3 << 32) | (n4 << 48)
}

// generateMagicNumber generates a magic number candidate.
func generateMagicNumber() uint64 {
	return getRandom64BitNumber() & getRandom64BitNumber() & getRandom64BitNumber()
}

// findMagicNumber checks whether a magic number candidate is viable.
func findMagicNumber(square, relevantBits int, bishop bool) uint64 {
	var (
		occupancies [4096]uint64
		attacks     [4096]uint64
		attackMask  uint64
	)
	if bishop {
		attackMask = maskBishopMoves(square)
	} else {
		attackMask = maskRookMoves(square)
	}
	occupancyIndizes := 1 << relevantBits

	for i := 0; i < occupancyIndizes; i++ {
		occupancies[i] = setOccupancy(relevantBits, i, attackMask)

		if bishop {
			attacks[i] = generateBishopMovesOnTheFly(square, occupancies[i])
		} else {
			attacks[i] = generateRookMovesOnTheFly(square, occupancies[i])
		}
	}

	for randomCount := 0; randomCount < 100000000; randomCount++ {
		magicNumber := generateMagicNumber()
		var fail int
		var usedAttacks [4096]uint64

		if popCount((attackMask*magicNumber)&0xFF00000000000000) < 6 {
			continue
		}

		for i := 0; fail == 0 && i < occupancyIndizes; i++ {
			magicIndex := uint64((occupancies[i] * magicNumber) >> (64 - relevantBits))

			if usedAttacks[magicIndex] == 0 {
				usedAttacks[magicIndex] = attacks[i]
			} else if usedAttacks[magicIndex] != attacks[i] {
				fail = 1
			}
		}

		if fail == 0 {
			return magicNumber
		}
	}
	return 0
}

var magicNumberRook [64]uint64
var magicNumberBishop [64]uint64

// initMagicNumbers initalizes the magicnumbers.
func initMagicNumbers() {

	for i := 0; i < 64; i++ {
		magicNumberRook[i] = findMagicNumber(i, relevantBitsRook[i], false)
	}
	for i := 0; i < 64; i++ {
		magicNumberBishop[i] = findMagicNumber(i, relevantBitsBishop[i], true)
	}
}
