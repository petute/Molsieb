package main

import "fmt"

//search searches for the best move
func search(white bool) move {
	moveList := getLegalMoves(white)
	score := make([]int, len(moveList))
	max := -100
	var index int
	for i, move := range moveList {
		p := makeMove(move, white, position)
		score[i] = negamax(p, 2, white)
		if score[i] > max {
			max = score[i]
			index = i
		}
	}
	fmt.Println(moveList)
	fmt.Println(score)
	return moveList[index]
}

// negamax is my implementation of the negamax algorithm.
func negamax(position pos, depth int, white bool) int {
	moveList := getLegalMoves(white)
	if depth == 0 || moveList == nil {
		return evaluate(position)
	}
	max := -100
	for _, move := range moveList {
		p := makeMove(move, white, position)
		value := -negamax(p, depth-1, !white)

		if value > max {
			max = value
		}
	}
	return max
}
