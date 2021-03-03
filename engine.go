package main

import "fmt"

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

//getBit is a function to get a bit on a certain square.
func getBit(bitboard, square uint64) (rgw uint64) {
	if bitboard&(1<<square) != 0 {
		rgw = 1
	}
	return rgw
}

//setBit sets a bit on a certain square to 1.
func setBit(bitboard, square uint64) uint64 {
	return bitboard | (1 << square)
}

//popCount counts the bits in a bitboard.
func popCount(bitboard uint64) (count int) {
	if bitboard != 0 && (bitboard&bitboard-1 == 0) {
		count = 1
	} else if bitboard != 0 {
		for bitboard != 0 {
			count++
			bitboard &= bitboard - 1
		}
	}
	return count
}

//printBitboard is a debug function to print every bit in a bitboard in a 8x8 square.
func printBitboard(bitboard uint64) {
	fmt.Printf(" Bitboard %d:\n", bitboard)
	for rank := uint64(0); rank < 8; rank++ {
		for file := uint64(0); file < 8; file++ {
			if file == 0 {
				fmt.Printf("\033[31m %d\033[0m", 8-rank)
			}
			square := rank*8 + file
			fmt.Printf(" %d", getBit(bitboard, square))
		}
		fmt.Printf("\n")
	}
	fmt.Println("\033[31m", "  A B C D E F G H")
}

//initPosition sets the starting position for all 8 bitboards
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

//countMaterial counts the material for one side
func countMaterial(color uint64) (mat int) {
	mat += popCount(color & position.pawns)
	mat += popCount(color&position.bishops) * 3
	mat += popCount(color&position.knights) * 3
	mat += popCount(color&position.rooks) * 5
	mat += popCount(color&position.queens) * 9

	return mat
}

//evaluate evaluates the position.
func evaluate() int {
	return countMaterial(position.white) - countMaterial(position.black)
}

func main() {
	initPosition()
	fmt.Println(evaluate())
}
