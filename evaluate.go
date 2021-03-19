package main

// countMaterial counts the material for one side.
func countMaterial(color uint64) (mat int) {
	mat += popCount(color & position.pawns)
	mat += popCount(color&position.bishops) * 3
	mat += popCount(color&position.knights) * 3
	mat += popCount(color&position.rooks) * 5
	mat += popCount(color&position.queens) * 9

	return mat
}

// evaluate evaluates the position.
func evaluate(position pos) int {
	return countMaterial(position.white) - countMaterial(position.black)
}
