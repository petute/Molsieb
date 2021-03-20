package main

import (
	"fmt"
)

// <<--------------------------------- Masks --------------------------------->>

var (
	notAFile  uint64 = 9187201950435737471
	notHFile  uint64 = 18374403900871474942
	notABFile uint64 = 4557430888798830399
	notGHFile uint64 = 18229723555195321596
)
var (
	pawnAttacks   [2][64]uint64
	pawnMoves     [2][64]uint64
	knightAttacks [64]uint64
	kingAttacks   [64]uint64
	rookMasks     [64]uint64
	bishopMasks   [64]uint64
	bishopAttacks [64][512]uint64
	rookAttacks   [64][4096]uint64
)

// maskPawnAttacks generates all possible attacks for pawns. pawnAttacks[0][x] == white
func maskPawnAttacks() (pawnAttacks [2][64]uint64) {
	for i := 0; i < 64; i++ {
		pawnAttacks[0][i] = setBit(0, i)
		pawnAttacks[0][i] = (pawnAttacks[0][i]&notAFile)>>7 ^ (pawnAttacks[0][i]&notHFile)>>9

		pawnAttacks[1][i] = setBit(0, i)
		pawnAttacks[1][i] = (pawnAttacks[1][i]&notHFile)<<7 ^ (pawnAttacks[1][i]&notAFile)<<9
	}
	return pawnAttacks
}

// generatePawnMoves masks possible pawn moves. pawnMoves[0][x] == white
func maskPawnMoves() (pawnMoves [2][64]uint64) {
	for i := 0; i < 64; i++ {
		pawnMoves[0][i] = setBit(0, i)
		pawnMoves[0][i] >>= 8
		if i/8 == 6 {
			pawnMoves[0][i] |= pawnMoves[0][i] >> 8
		}

		pawnMoves[1][i] = setBit(0, i)
		pawnMoves[1][i] <<= 8
		if i/8 == 1 {
			pawnMoves[1][i] |= pawnMoves[1][i] << 8
		}
	}
	return pawnMoves
}

// maskKnightMoves generates all possible moves for knights.
func maskKnightMoves(square int) (knightMove uint64) {
	knightMove = setBit(0, square)
	knightMove = (((knightMove >> 6) & notGHFile) ^
		((knightMove >> 10) & notABFile) ^
		((knightMove << 6) & notABFile) ^
		((knightMove << 10) & notGHFile) ^
		((knightMove >> 17) & notAFile) ^
		((knightMove >> 15) & notHFile) ^
		((knightMove << 17) & notHFile) ^
		((knightMove << 15) & notAFile))

	return knightMove
}

// maskKingMoves generates all possible moves for kings.
func maskKingMoves(square int) (kingMove uint64) {
	kingMove = setBit(0, square)
	kingMove = (((kingMove >> 1) & notAFile) ^
		((kingMove >> 9) & notAFile) ^
		((kingMove << 7) & notAFile) ^
		(kingMove >> 8) ^
		(kingMove << 8) ^
		((kingMove >> 7) & notHFile) ^
		((kingMove << 1) & notHFile) ^
		((kingMove << 9) & notHFile))

	return kingMove
}

// maskRookMoves generates all relevant occupancy bits of rooks for magic bitboards.
func maskRookMoves(square int) (rookMove uint64) {
	rank := square / 8
	file := square % 8

	for r := rank + 1; r <= 6; r++ {
		rookMove |= (1 << (r*8 + file))
	}
	for r := rank - 1; r >= 1; r-- {
		rookMove |= (1 << (r*8 + file))
	}
	for f := file + 1; f <= 6; f++ {
		rookMove |= (1 << (rank*8 + f))
	}
	for f := file - 1; f >= 1; f-- {
		rookMove |= (1 << (rank*8 + f))
	}

	return rookMove
}

// maskBishopMoves generates all relevant occupancy bits of bishops for magic bitboards.
func maskBishopMoves(square int) (bishopMove uint64) {
	rank := square / 8
	file := square % 8

	for r, f := rank+1, file+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		bishopMove |= (1 << (r*8 + f))
	}
	for r, f := rank-1, file+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		bishopMove |= (1 << (r*8 + f))
	}
	for r, f := rank+1, file-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		bishopMove |= (1 << (r*8 + f))
	}
	for r, f := rank-1, file-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		bishopMove |= (1 << (r*8 + f))
	}

	return bishopMove
}

// generateRookMovesOnTheFly generates the rook moves for a certain blockboard (position).
func generateRookMovesOnTheFly(square int, blockboard uint64) (rookMoves uint64) {
	var rank, file int

	rank = square / 8
	file = square % 8

	for r := rank + 1; r <= 7; r++ {
		rookMoves |= (1 << (r*8 + file))
		if (1<<(r*8+file))&blockboard != 0 {
			break
		}
	}
	for r := rank - 1; r >= 0; r-- {
		rookMoves |= (1 << (r*8 + file))
		if (1<<(r*8+file))&blockboard != 0 {
			break
		}

	}
	for f := file + 1; f <= 7; f++ {
		rookMoves |= (1 << (rank*8 + f))
		if (1<<(rank*8+f))&blockboard != 0 {
			break
		}
	}
	for f := file - 1; f >= 0; f-- {
		rookMoves |= (1 << (rank*8 + f))
		if (1<<(rank*8+f))&blockboard != 0 {
			break
		}
	}

	return rookMoves
}

// generateBishopMovesOnTheFly generates the bishop moves for a certain blockboard (position).
func generateBishopMovesOnTheFly(square int, blockboard uint64) (bishopMoves uint64) {
	var rank, file int

	rank = square / 8
	file = square % 8

	for r, f := rank+1, file+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	for r, f := rank-1, file+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	for r, f := rank+1, file-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	for r, f := rank-1, file-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		bishopMoves |= (1 << (r*8 + f))
		if (1<<(r*8+f))&blockboard != 0 {
			break
		}
	}
	return bishopMoves
}

// <<--------------------------------- Magic --------------------------------->>

// Amount oft relevant occupancy bits per square for the bishop.
var relevantBitsBishop = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// Amount oft relevant occupancy bits per square for the rook.
var relevantBitsRook = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

// state is used to generate the random numbers.
var state uint32 = 1804289383

// magicNumberRook and magicNumberBishop are magic numbers for every square for the rook and the bishop.
var magicNumbersBishop = [64]uint64{
	18018832060792964,
	9011737055478280,
	4531088509108738,
	74316026439016464,
	396616115700105744,
	2382975967281807376,
	1189093273034424848,
	270357282336932352,
	1131414716417028,
	2267763835016,
	2652629010991292674,
	283717117543424,
	4411067728898,
	1127068172552192,
	288591295206670341,
	576743344005317120,
	18016669532684544,
	289358613125825024,
	580966009790284034,
	1126071732805635,
	37440604846162944,
	9295714164029260800,
	4098996805584896,
	9223937205167456514,
	153157607757513217,
	2310364244010471938,
	95143507244753921,
	9015995381846288,
	4611967562677239808,
	9223442680644702210,
	64176571732267010,
	7881574242656384,
	9224533161443066400,
	9521190163130089986,
	2305913523989908488,
	9675423050623352960,
	9223945990515460104,
	2310346920227311616,
	7075155703941370880,
	4755955152091910658,
	146675410564812800,
	4612821438196357120,
	4789475436135424,
	1747403296580175872,
	40541197101432897,
	144397831292092673,
	1883076424731259008,
	9228440811230794258,
	360435373754810368,
	108227545293391872,
	4611688277597225028,
	3458764677302190090,
	577063357723574274,
	9165942875553793,
	6522483364660839184,
	1127033795058692,
	2815853729948160,
	317861208064,
	5765171576804257832,
	9241386607448426752,
	11258999336993284,
	432345702206341696,
	9878791228517523968,
	4616190786973859872,
}
var magicNumbersRook = [64]uint64{
	9979994641325359136,
	90072129987412032,
	180170925814149121,
	72066458867205152,
	144117387368072224,
	216203568472981512,
	9547631759814820096,
	2341881152152807680,
	140740040605696,
	2316046545841029184,
	72198468973629440,
	81205565149155328,
	146508277415412736,
	703833479054336,
	2450098939073003648,
	576742228899270912,
	36033470048378880,
	72198881818984448,
	1301692025185255936,
	90217678106527746,
	324684134750365696,
	9265030608319430912,
	4616194016369772546,
	2199165886724,
	72127964931719168,
	2323857549994496000,
	9323886521876609,
	9024793588793472,
	562992905192464,
	2201179128832,
	36038160048718082,
	36029097666947201,
	4629700967774814240,
	306244980821723137,
	1161084564161792,
	110340390163316992,
	5770254227613696,
	2341876206435041792,
	82199497949581313,
	144120019947619460,
	324329544062894112,
	1152994210081882112,
	13545987550281792,
	17592739758089,
	2306414759556218884,
	144678687852232706,
	9009398345171200,
	2326183975409811457,
	72339215047754240,
	18155273440989312,
	4613959945983951104,
	145812974690501120,
	281543763820800,
	147495088967385216,
	2969386217113789440,
	19215066297569792,
	180144054896435457,
	2377928092116066437,
	9277424307650174977,
	4621827982418248737,
	563158798583922,
	5066618438763522,
	144221860300195844,
	281752018887682,
}

// setOccupancy generates the relevant occupancy bitboard for a given rook or bishop moves bitboard.
func setOccupancy(bitsInMask, index int, moveMask uint64) (occupancy uint64) {
	for i := 0; i < bitsInMask; i++ {
		square := getLS1BIndex(moveMask)
		moveMask = popBit(moveMask, square)

		if index&(1<<i) != 0 {
			occupancy |= (1 << square)
		}
	}
	return occupancy
}

// getRandom32BitNumber generates pseudoRandom numbers (XORSHIFT32).
func getRandom32BitNumber() uint32 {
	number := state

	number ^= number << 13
	number ^= number >> 17
	number ^= number << 5

	state = number

	return number
}

// getRandom64BitNumber generates a random 64 bit pseudo legal number. (FFFF == 65535 == 16 bits)
func getRandom64BitNumber() uint64 {
	n1 := uint64(getRandom32BitNumber()) & 0xFFFF
	n2 := uint64(getRandom32BitNumber()) & 0xFFFF
	n3 := uint64(getRandom32BitNumber()) & 0xFFFF
	n4 := uint64(getRandom32BitNumber()) & 0xFFFF

	return n1 | (n2 << 16) | (n3 << 32) | (n4 << 48)
}

// generateMagicNumber generates a magic number candidate.
func generateMagicNumber() uint64 {
	return getRandom64BitNumber() & getRandom64BitNumber() & getRandom64BitNumber()
}

// findMagicNumber checks whether a magic number candidate is viable.
func findMagicNumber(square, relevantBits int, bishop bool) uint64 {
	var (
		occupancies [4096]uint64
		attacks     [4096]uint64
		attackMask  uint64
	)
	if bishop {
		attackMask = maskBishopMoves(square)
	} else {
		attackMask = maskRookMoves(square)
	}
	occupancyIndizes := 1 << relevantBits

	for i := 0; i < occupancyIndizes; i++ {
		occupancies[i] = setOccupancy(relevantBits, i, attackMask)

		if bishop {
			attacks[i] = generateBishopMovesOnTheFly(square, occupancies[i])
		} else {
			attacks[i] = generateRookMovesOnTheFly(square, occupancies[i])
		}
	}

	for randomCount := 0; randomCount < 100000000; randomCount++ {
		magicNumber := generateMagicNumber()
		var fail int
		var usedAttacks [4096]uint64

		if popCount((attackMask*magicNumber)&0xFF00000000000000) < 6 {
			continue
		}

		for i := 0; fail == 0 && i < occupancyIndizes; i++ {
			magicIndex := uint64((occupancies[i] * magicNumber) >> (64 - relevantBits))

			if usedAttacks[magicIndex] == 0 {
				usedAttacks[magicIndex] = attacks[i]
			} else if usedAttacks[magicIndex] != attacks[i] {
				fail = 1
			}
		}

		if fail == 0 {
			return magicNumber
		}
	}
	return 0
}

// initMagicNumbers initializes the magicnumbers.
func initMagicNumbers() {
	fmt.Println("rook")
	for i := 0; i < 64; i++ {
		fmt.Printf("%d,\n", findMagicNumber(i, relevantBitsRook[i], false))
	}
	fmt.Println("bishop")
	for i := 0; i < 64; i++ {
		fmt.Printf("%d,\n", findMagicNumber(i, relevantBitsBishop[i], true))
	}
}

// <<------------------------------- initMoves ------------------------------->>

// initLeaperAttacks initializes the attack tables for the leaper pieces
func initLeaperAttacks() {
	pawnAttacks = maskPawnAttacks()
	pawnMoves = maskPawnMoves()
	for i := 0; i < 64; i++ {
		knightAttacks[i] = maskKnightMoves(i)
		kingAttacks[i] = maskKingMoves(i)
	}
}

// initSliderAttacks initializes the attack tables for the slider pieces.
func initSliderAttacks(bishop bool) {
	for i := 0; i < 64; i++ {
		bishopMasks[i] = maskBishopMoves(i)
		rookMasks[i] = maskRookMoves(i)
		var attackMask uint64

		if bishop {
			attackMask = bishopMasks[i]
		} else {
			attackMask = rookMasks[i]
		}
		relevantBitsCount := popCount(attackMask)
		occupancyIndizes := 1 << relevantBitsCount

		for j := 0; j < occupancyIndizes; j++ {
			if bishop {
				occupancy := setOccupancy(relevantBitsCount, j, attackMask)
				magicIndex := int((occupancy * magicNumbersBishop[i]) >> (64 - relevantBitsBishop[i]))

				bishopAttacks[i][magicIndex] = generateBishopMovesOnTheFly(i, occupancy)
			} else {
				occupancy := setOccupancy(relevantBitsCount, j, attackMask)
				magicIndex := int((occupancy * magicNumbersRook[i]) >> (64 - relevantBitsRook[i]))

				rookAttacks[i][magicIndex] = generateRookMovesOnTheFly(i, occupancy)
			}
		}
	}
}

// <<------------------------------- Generator ------------------------------->>

type move struct {
	fromSquare int
	toSquare   int
	pieceType  string
}

// getBishopAttacks returns the attack for a square and occupancy.
func getBishopAttacks(square int, occupancy uint64) uint64 {
	occupancy &= bishopMasks[square]
	occupancy *= magicNumbersBishop[square]
	occupancy >>= 64 - relevantBitsBishop[square]

	return bishopAttacks[square][occupancy]
}

// getRookAttacks returns the attack for a square and occupancy.
func getRookAttacks(square int, occupancy uint64) uint64 {
	occupancy &= rookMasks[square]
	occupancy *= magicNumbersRook[square]
	occupancy >>= 64 - relevantBitsRook[square]

	return rookAttacks[square][occupancy]
}

// convertMoves converts the moves from bitboards to move type.
func convertMoves(square int, bitboard uint64, pieceType string) (moveList []move) {
	moveList = make([]move, 0, 10)
	for bitboard > 0 {
		toSquare := getLS1BIndex(bitboard)
		moveList = append(moveList, move{fromSquare: square, toSquare: toSquare, pieceType: pieceType})
		bitboard = popBit(bitboard, toSquare)
	}
	return moveList
}

// getAttackMap returns a bitboard with all attacks for one side.
func getAttackMap(white bool, occupancy uint64) uint64 {
	var color uint64
	var attacks uint64

	if white {
		color = position.white
	} else {
		color = position.black
	}

	pieceCount := popCount(color)
	rooks := color & position.rooks
	pawns := color & position.pawns
	knights := color & position.knights
	bishops := color & position.bishops
	queens := color & position.queens
	kings := color & position.kings

	for pieceCount > 0 {
		count := 0
		if rooks > 0 {
			square := getLS1BIndex(rooks)
			rooks = popBit(rooks, square)

			attacks |= getRookAttacks(square, occupancy)
			count++
		}
		if queens > 0 {
			square := getLS1BIndex(queens)
			queens = popBit(queens, square)

			attacks |= getRookAttacks(square, occupancy) | getBishopAttacks(square, occupancy)
			count++
		}
		if bishops > 0 {
			square := getLS1BIndex(bishops)
			bishops = popBit(bishops, square)

			attacks |= getBishopAttacks(square, occupancy)
			count++
		}
		if knights > 0 {
			square := getLS1BIndex(knights)
			knights = popBit(knights, square)

			attacks |= knightAttacks[square]
			count++
		}
		if pawns > 0 {
			square := getLS1BIndex(pawns)
			pawns = popBit(pawns, square)

			if white {
				attacks |= pawnAttacks[0][square]
			} else {
				attacks |= pawnAttacks[1][square]
			}
			count++
		}
		if kings > 0 {
			square := getLS1BIndex(kings)
			kings = popBit(kings, square)

			attacks |= kingAttacks[square]
			count++
		}
		pieceCount -= count
	}

	return attacks
}

// getLegalMoves returns a []move with all possible attacks.
func getLegalMoves(white bool) (moveList []move) {
	moveList = make([]move, 0, 35)
	var (
		color uint64
	)
	occupancy := position.white | position.black
	if white {
		color = position.white
	} else {
		color = position.black
	}

	pieceCount := popCount(color)
	rooks := color & position.rooks
	pawns := color & position.pawns
	knights := color & position.knights
	bishops := color & position.bishops
	queens := color & position.queens
	kings := color & position.kings
	attacks := getAttackMap(!white, occupancy)

	for pieceCount > 0 {
		count := 0
		if pawns > 0 {
			var move uint64
			square := getLS1BIndex(pawns)
			pawns = popBit(pawns, square)

			if white {
				move = pawnMoves[0][square] &^ occupancy
				if square/8 == 6 && move != pawnMoves[0][square] {
					move &^= 1 << (square - 16)
				}
				move |= pawnAttacks[0][square] & (position.black | (position.enPassant & 0xFF0000))
			} else {
				move = pawnMoves[1][square] &^ occupancy
				if square/8 == 1 && move != pawnMoves[1][square] {
					move &^= occupancy & (1 >> (square + 16))
				}
				move |= pawnAttacks[1][square] & (position.white | (position.enPassant & 0xFF0000000000))
			}

			moveList = append(moveList, convertMoves(square, move, "pawn")...)
			count++
		}
		if rooks > 0 {
			square := getLS1BIndex(rooks)
			rooks = popBit(rooks, square)

			moveList = append(moveList, convertMoves(square, getRookAttacks(square, occupancy)&^color, "rook")...)
			count++
		}
		if bishops > 0 {
			square := getLS1BIndex(bishops)
			bishops = popBit(bishops, square)

			moveList = append(moveList, convertMoves(square, getBishopAttacks(square, occupancy)&^color, "bishop")...)
			count++
		}
		if knights > 0 {
			square := getLS1BIndex(knights)
			knights = popBit(knights, square)

			moveList = append(moveList, convertMoves(square, knightAttacks[square]&^color, "knight")...)
			count++
		}
		if queens > 0 {
			square := getLS1BIndex(queens)
			queens = popBit(queens, square)

			moveList = append(moveList, convertMoves(square, (getRookAttacks(square, occupancy)&getBishopAttacks(square, occupancy))&^color, "queen")...)
			count++
		}
		if kings > 0 {
			square := getLS1BIndex(kings)
			kings = popBit(kings, square)

			moveList = append(moveList, convertMoves(square, (kingAttacks[square]&^color)&^attacks, "king")...)
			count++
		}
		pieceCount -= count
	}

	if getBit(attacks, getLS1BIndex(color&position.kings)) == 1 {
		var checkMoves []move
		for _, move := range moveList {
			if move.pieceType == "king" {
				checkMoves = append(checkMoves, move)
			} else if getBit(attacks, move.toSquare) == 1 {
				test := setBit(occupancy, move.toSquare)
				test = popBit(test, move.fromSquare)
				tAttack := getAttackMap(!white, test)
				if getBit(tAttack, getLS1BIndex(color&position.kings)) != 1 {
					checkMoves = append(checkMoves, move)
				}
			}
		}
		moveList = checkMoves
	}
	return moveList
}

// makeMove makes a move and en-passant. TODO: Check for checks.
func makeMove(move move, white bool, position pos) pos {
	var capture bool
	if white {
		position.white = popBit(position.white, move.fromSquare)
		if getBit(position.black, move.toSquare) != 0 {
			capture = true
			position.black = popBit(position.black, move.toSquare)
		} else if getBit(position.enPassant, move.toSquare) != 0 {
			capture = true
			position.black = popBit(position.black, move.toSquare+8)
			position.pawns = popBit(position.pawns, move.toSquare+8)

		}
		position.white = setBit(position.white, move.toSquare)
		position.enPassant = position.enPassant &^ 0xFF0000
		position.moveNumber++
		position.moveRule++
	} else {
		position.black = popBit(position.black, move.fromSquare)
		if getBit(position.white, move.toSquare) != 0 {
			capture = true
			position.white = popBit(position.white, move.toSquare)
		} else if getBit(position.enPassant, move.toSquare) != 0 {
			capture = true
			position.white = popBit(position.white, move.toSquare-8)
			position.pawns = popBit(position.pawns, move.toSquare-8)
		}
		position.black = setBit(position.black, move.toSquare)
		position.enPassant = position.enPassant &^ 0xFF0000000000
	}

	if move.pieceType == "pawn" && move.toSquare-move.fromSquare == 16 {
		position.enPassant = setBit(position.enPassant, move.fromSquare+8)
	} else if move.pieceType == "pawn" && move.fromSquare-move.toSquare == 16 {
		position.enPassant = setBit(position.enPassant, move.fromSquare-8)
	}

	if capture {
		position.moveRule = 0
		if getBit(position.pawns, move.toSquare) == 1 {
			position.pawns = popBit(position.pawns, move.toSquare)
		} else if getBit(position.rooks, move.toSquare) == 1 {
			position.rooks = popBit(position.rooks, move.toSquare)
		} else if getBit(position.knights, move.toSquare) == 1 {
			position.knights = popBit(position.knights, move.toSquare)
		} else if getBit(position.bishops, move.toSquare) == 1 {
			position.bishops = popBit(position.bishops, move.toSquare)
		} else if getBit(position.queens, move.toSquare) == 1 {
			position.queens = popBit(position.queens, move.toSquare)
		}
	}

	switch move.pieceType {
	case "pawn":
		position.pawns = popBit(position.pawns, move.fromSquare)
		position.pawns = setBit(position.pawns, move.toSquare)
	case "bishop":
		position.bishops = popBit(position.bishops, move.fromSquare)
		position.bishops = setBit(position.bishops, move.toSquare)
	case "knight":
		position.knights = popBit(position.knights, move.fromSquare)
		position.knights = setBit(position.knights, move.toSquare)
	case "rook":
		position.rooks = popBit(position.rooks, move.fromSquare)
		position.rooks = setBit(position.rooks, move.toSquare)
	case "queen":
		position.queens = popBit(position.queens, move.fromSquare)
		position.queens = setBit(position.queens, move.toSquare)
	case "king":
		position.kings = popBit(position.kings, move.fromSquare)
		position.kings = setBit(position.kings, move.toSquare)
	}

	return position
}
