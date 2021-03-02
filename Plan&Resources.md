# Chess Engine

This is a collection of resources and at the same time a (**ambitious**) plan for my engine Molsieb.

## Plan & Resources

Most of these links lead to the [ChessProgrammingWiki](https://www.chessprogramming.org/Main_Page)

Priority of implementation:

**M** Must

**S** Should

### Board representation

- [Bitboards](https://www.chessprogramming.org/Bitboards)

### GUI communication

- [UCI](https://www.chessprogramming.org/UCI)

### [Move generation](https://www.chessprogramming.org/Move_Generation#Bitboards)

- [magic bitboards](https://www.chessprogramming.org/Magic_Bitboards)

### [Move evaluation](https://www.chessprogramming.org/Evaluation)
- **M**
	- [material](https://www.chessprogramming.org/Material)
- **S**
	- position
	- pawns
	- mobility
	- king

### Move search
- **M**
	- [negamax](https://www.chessprogramming.org/Negamax)
    - iterative deepening
- **S**
	- quiescence search

### Pruning
- **M**
	- [alpha/beta](https://www.chessprogramming.org/Alpha-Beta)
- **S**
	- transposition tables
		- zobrist hashing
	- Forward Pruning
		- Null Move Pruning
		- late move reductions
		- futility pruning

### [Ordering](https://www.chessprogramming.org/Move_Ordering)
- **S**
	- pv
	- killer
	- history

### Database
- **M**
	- [opening](https://www.chessprogramming.org/Opening_Book)
	- [endgame](https://www.chessprogramming.org/Endgame_Tablebases)

### Videos:
- [Tord Romstad (Stockfish) - How modern chess programs work](https://vimeo.com/216463393)
- [Maksim Korzh (Chess Programming YT) - Bitboard chess engine in C (95 vids)](https://www.youtube.com/playlist?list=PLmN0neTso3Jxh8ZIylk74JpwfiWNI76Cs)