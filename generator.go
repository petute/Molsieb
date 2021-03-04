package main

var (
	notAFile  uint64 = 9187201950435737471
	notHFile  uint64 = 18374403900871474942
)

// GeneratePawnAttacks generates all possible attacks for pawns and returns them in an [2][64]int array.
func GeneratePawnAttacks() (pawnAttacks [2][64]uint64) {
	for i := uint64(0); i < 64; i++ {
		pawnAttacks[0][i] = setBit(0, i)
		pawnAttacks[0][i] = (pawnAttacks[0][i]&notAFile)>>7 ^ (pawnAttacks[0][i]&notHFile)>>9

		pawnAttacks[1][i] = setBit(0, i)
		pawnAttacks[1][i] = (pawnAttacks[1][i]&notHFile)<<7 ^ (pawnAttacks[1][i]&notAFile)<<9
	}
	return pawnAttacks
}
