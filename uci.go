package main

import (
	"strconv"
	"strings"
)

// parseFENString sets a position from a FENstring
func parseFenString(groups []string) {
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
	for _, ascii := range castle {
		switch ascii {
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

	if groups[3] != "-" {
		square := 0
		ascii := []rune(groups[3])

		square = (int(ascii[0]) - 97) + (-int(ascii[1]-56))*8
		position.enPassant = setBit(position.enPassant, square)
	}

	if value, err := strconv.Atoi(groups[4]); err == nil {
		position.moveRule = value
	}

	if value, err := strconv.Atoi(groups[5]); err == nil {
		position.moveNumber = value
	}
}

func getMoveFromString(m string) move {
	var mov move
	ascii := []rune(m)

	mov.fromSquare = (int(ascii[0]) - 97) + (-int(ascii[1]-56))*8
	mov.toSquare = (int(ascii[2]) - 97) + (-int(ascii[3]-56))*8
	if getBit(position.pawns, mov.fromSquare) == 1 {
		mov.pieceType = "pawn"
	} else if getBit(position.bishops, mov.fromSquare) == 1 {
		mov.pieceType = "bishop"
	} else if getBit(position.knights, mov.fromSquare) == 1 {
		mov.pieceType = "knight"
	} else if getBit(position.rooks, mov.fromSquare) == 1 {
		mov.pieceType = "rook"
	} else if getBit(position.queens, mov.fromSquare) == 1 {
		mov.pieceType = "queen"
	} else if getBit(position.kings, mov.fromSquare) == 1 {
		mov.pieceType = "king"
	}
	return mov
}
