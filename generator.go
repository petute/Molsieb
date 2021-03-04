package main

var (
	notAFile  uint64 = 9187201950435737471
	notHFile  uint64 = 18374403900871474942
	notABFile uint64 = 4557430888798830399
	notGHFile uint64 = 18229723555195321596
)

// <<-------------------------------- Attacks -------------------------------->>

// GeneratePawnAttacks generates all possible attacks for pawns.
func GeneratePawnAttacks() (pawnAttacks [2][64]uint64) {
	for i := uint64(0); i < 64; i++ {
		pawnAttacks[0][i] = setBit(0, i)
		pawnAttacks[0][i] = (pawnAttacks[0][i]&notAFile)>>7 ^ (pawnAttacks[0][i]&notHFile)>>9

		pawnAttacks[1][i] = setBit(0, i)
		pawnAttacks[1][i] = (pawnAttacks[1][i]&notHFile)<<7 ^ (pawnAttacks[1][i]&notAFile)<<9
	}
	return pawnAttacks
}

// GenerateKnightAttacks generates all possible attacks for knights.
func GenerateKnightAttacks() (knightAttacks [64]uint64) {
	for i := uint64(0); i < 64; i++ {
		knightAttacks[i] = setBit(0, i)
		knightAttacks[i] = ((knightAttacks[i])>>6)&notGHFile ^
			((knightAttacks[i] >> 10) & notABFile) ^
			((knightAttacks[i] << 6) & notABFile) ^
			((knightAttacks[i] << 10) & notGHFile) ^
			((knightAttacks[i] >> 17) & notAFile) ^
			((knightAttacks[i] >> 15) & notHFile) ^
			((knightAttacks[i] << 17) & notHFile) ^
			((knightAttacks[i] << 15) & notAFile)
	}
	return knightAttacks
}
