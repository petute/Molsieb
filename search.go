package main

import "fmt"

var depth int
var stopSearch bool

//search searches for the best move
func search() move {
	moveList := getLegalMoves(position.color)
	fmt.Println(moveList)
	score := make([]int, len(moveList))
	max := -100
	var index int
	for i, move := range moveList {
		p := makeMove(move, position)
		score[i] = negamax(p, depth)
		if score[i] > max {
			max = score[i]
			index = i
		}
	}
	return moveList[index]
}

// negamax is my implementation of the negamax algorithm.
func negamax(position pos, ply int) int {
	moveList := getLegalMoves(position.color)
	if ply == 0 || moveList == nil || stopSearch {
		return evaluate(position)
	}
	max := -100
	for _, move := range moveList {
		p := makeMove(move, position)
		value := -negamax(p, ply-1)

		if value > max {
			max = value
		}
	}
	return max
}
