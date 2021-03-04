package main

import "fmt"

// getBit is a function to get a bit on a certain square.
func getBit(bitboard, square uint64) (rgw uint64) {
	if bitboard&(1<<square) != 0 {
		rgw = 1
	}
	return rgw
}

// setBit sets a bit on a certain square 1.
func setBit(bitboard, square uint64) uint64 {
	return bitboard | (1 << square)
}

// popBit sets a bit on a certain square from 1 to 0.
func popBit(bitboard, square uint64) uint64 {
	if getBit(bitboard, square) != 0 {
		bitboard ^= (1 << square)
	}
	return bitboard
}

// popCount counts the bits in a bitboard.
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

// printBitboard is a debug function to print every bit in a bitboard in a 8x8 square.
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
