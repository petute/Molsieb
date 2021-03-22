package main

import (
	"fmt"
	"strconv"
	"strings"
)

// getMoveFromString converts a string to a move.
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

// parseFENString sets a position from a FENstring
func parseFenString(groups []string) {
	positionStr := strings.Split(groups[0], "/")
	square := 0
	for _, row := range positionStr {
		runes := []rune(row)
		for _, ascii := range runes {
			if ascii > 97 {
				position.black = setBit(position.black, square)
				ascii -= 32
			} else if ascii > 65 {
				position.white = setBit(position.white, square)
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

// uciIn processes the UCI commands sent by the GUI.
func uci(in string) {
	slice := strings.Split(in, " ")
	switch slice[0] {
	case "uci":
		fmt.Println("id name Molsieb")
		fmt.Println("id author Petute")
		fmt.Println("uciok")
	case "position":
		handlePosition(slice[1:])
	case "go":
		stopSearch = false
		handleGo(slice[1:])
	case "stop":
		stopSearch = true
	default:
		fmt.Println("This command is not implemented.")
	}
}

// handlePosition handles the UCI position command.
func handlePosition(slice []string) {
	if slice[0] == "startpos" {
		position.pawns = 71776119061282560
		position.knights = 4755801206503243842
		position.bishops = 2594073385365405732
		position.rooks = 9295429630892703873
		position.kings = 1152921504606846992
		position.queens = 576460752303423496
		position.black = 65535
		position.white = 18446462598732840960
		position.castle = 15
		position.moveNumber = 0
		position.moveRule = 0
		position.color = true
	} else if slice[0] == "fen" {
		parseFenString(slice[1:7])
	}
	if len(slice) > 8 && slice[7] == "moves" {
		for _, m := range slice[8:] {
			mov := getMoveFromString(m)
			position = makeMove(mov, position)
		}
	} else if len(slice) > 2 && slice[1] == "moves" {
		for _, m := range slice[2:] {
			mov := getMoveFromString(m)
			position = makeMove(mov, position)
		}
	}
}

// handleGo handles the UCI go command.
func handleGo(slice []string) {
	for len(slice) >= 1 {
		c := 1
		switch slice[0] {
		case "wtime":
			if wtime, err := strconv.ParseFloat(slice[1], 64); err == nil {
				game.wtime = wtime
			}
			c++
		case "btime":
			if btime, err := strconv.ParseFloat(slice[1], 64); err == nil {
				game.btime = btime
			}
			c++
		case "winc":
			if winc, err := strconv.ParseFloat(slice[1], 64); err == nil {
				game.winc = winc
			}
			c++
		case "binc":
			if binc, err := strconv.ParseFloat(slice[1], 64); err == nil {
				game.binc = binc
			}
			c++
		case "depth":
			if d, err := strconv.Atoi(slice[1]); err == nil {
				depth = d
			}
			c++
		case "infinite":
			depth = 100
		default:
			fmt.Println("This command is not implemented.")
		}
		slice = slice[c:]
	}
	fmt.Println(search())
}
