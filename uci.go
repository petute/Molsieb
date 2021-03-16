package main

import (
	"strconv"
	"strings"
)

// parseFENString sets a position from a FENstring
func parseFenString(fen string) {
	groups := strings.Split(fen, " ")
	positionStr := strings.Split(groups[0], "/")
	square := 0
	for _, row := range positionStr {
		runes := []rune(row)
		for _, ascii := range runes {
			if ascii > 65 && ascii < 98 {
				position.white = setBit(position.white, square)
			} else if ascii > 97 {
				ascii -= 32
				position.black = setBit(position.black, square)
			}
			switch ascii {
			case 66:
				position.bishops = setBit(position.bishops, square)
			case 75:
				position.kings = setBit(position.kings, square)
			case 78:
				position.knights = setBit(position.knights, square)
			case 80:
				position.pawns = setBit(position.pawns, square)
			case 81:
				position.queens = setBit(position.queens, square)
			case 82:
				position.rooks = setBit(position.rooks, square)
			default:
				square += int(ascii) - 49
			}

			square++
		}
	}
	if groups[1] == "w" {
		position.color = true
	} else {
		position.color = false
	}
	castle := []rune(groups[2])
	for _, char := range castle {
		switch char {
		case 75:
			position.castle++
		case 81:
			position.castle += 2
		case 107:
			position.castle += 4
		case 113:
			position.castle += 8
		}
	}
	/* TODO: implement en-passant mechanism and then parse groups[3] here.*/

	if value, err := strconv.Atoi(groups[4]); err == nil {
		position.moveRule = value
	}

	if value, err := strconv.Atoi(groups[5]); err == nil {
		position.moveNumber = value
	}
}
