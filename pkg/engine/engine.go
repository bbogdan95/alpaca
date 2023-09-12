package engine

// SQ64 is an array that maps 120-square board indices to 64-square board indices.
// It allows for efficient conversion between the 120-square and 64-square representations of the board.
var SQ64 [BRD_SQ_NUM]int

// SQ120 is an array that maps 64-square board indices to 120-square board indices.
// It allows for efficient conversion between the 64-square and 120-square representations of the board.
var SQ120 [64]int

var FilesBrd [BRD_SQ_NUM]int
var RanksBrd [BRD_SQ_NUM]int

var PceChar = ".PNBRQKpnbrqk"
var SideChar = "wb-"
var RankChar = "12345678"
var FileChar = "abcdefgh"

// PieceBig is an array that helps identify if a piece is considered "big" (rook or queen).
// It maps piece types (represented by integers) to a boolean value (TRUE or FALSE).
var PieceBig = [13]int{FALSE, FALSE, TRUE, TRUE, TRUE, TRUE, TRUE, FALSE, TRUE, TRUE, TRUE, TRUE, TRUE}

// PieceMaj is an array that helps identify if a piece is considered "major" (rook, queen, or king).
// It maps piece types (represented by integers) to a boolean value (TRUE or FALSE).
var PieceMaj = [13]int{FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE}

// PieceMin is an array that helps identify if a piece is considered "minor" (bishop or knight).
// It maps piece types (represented by integers) to a boolean value (TRUE or FALSE).
var PieceMin = [13]int{FALSE, FALSE, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE, FALSE, FALSE}

// PieceVal is an array that assigns values to different piece types.
// It maps piece types (represented by integers) to their respective values.
var PieceVal = [13]int{0, 100, 325, 325, 550, 1000, 50000, 100, 325, 325, 550, 1000, 50000}

// PieceCol is an array that represents the color of a piece (WHITE or BLACK).
// It maps piece types (represented by integers) to their respective colors.
var PieceCol = [13]int{BOTH, WHITE, WHITE, WHITE, WHITE, WHITE, WHITE, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK}

// PiecePawn is a boolean array that indicates whether a piece type is a pawn (1) or not (0).
var PiecePawn = [13]int{FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE}

// PieceKnight is a boolean array that indicates whether a piece type is a knight (1) or not (0).
var PieceKnight = [13]int{FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE}

// PieceKing is a boolean array that indicates whether a piece type is a king (1) or not (0).
var PieceKing = [13]int{FALSE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE}

// PieceRookQueen is a boolean array that indicates whether a piece type is a rook or queen (1) or not (0).
var PieceRookQueen = [13]int{FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE, FALSE, FALSE, FALSE, TRUE, TRUE, FALSE}

// PieceBishopQueen is a boolean array that indicates whether a piece type is a bishop or queen (1) or not (0).
var PieceBishopQueen = [13]int{FALSE, FALSE, FALSE, TRUE, FALSE, TRUE, FALSE, FALSE, FALSE, TRUE, FALSE, TRUE, FALSE}

// PieceSlides is a boolean array that indicates whether a piece type can slide across the board (1) or not (0).
var PieceSlides = [13]int{FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, FALSE, FALSE, FALSE, TRUE, TRUE, TRUE, FALSE}

// used for movegen
var LoopSlidePieces = [8]int{WB, WR, WQ, 0, BB, BR, BQ, 0}
var LoopSlideIndex = [2]int{0, 4}
var LoopNonSlidePieces = [6]int{WN, WK, 0, BN, BK, 0}
var LoopNonSlideIndex = [2]int{0, 3}

// PieceDir is a two-dimensional array that stores the possible directions each piece type can move on a chessboard.
// It is used for move generation, helping to determine the valid move directions for each piece.
//
// Explanation of Values:
// - The array has 13 rows, each corresponding to a piece type (including empty squares), and 8 columns representing different directions.
// - Values in the array represent the relative positions a piece can move in a given direction.
// - For example, PieceDir[WN][0] represents the first direction a white knight can move, which is two squares up and one square left.
//
// Usage:
//   - This array is essential for generating legal moves for pieces, allowing the engine to identify valid move directions
//     and calculate potential destination squares.
//
// Note: The values in this array are relative positions on the board, and their interpretation may vary depending on the piece type.
var PieceDir = [13][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{0, 0, 0, 0, 0, 0, 0, 0},
	{-8, -19, -21, -12, 8, 19, 21, 12},
	{-9, -11, 11, 9, 0, 0, 0, 0},
	{-1, -10, 1, 10, 0, 0, 0, 0},
	{-1, -10, 1, 10, -9, -11, 11, 9},
	{-1, -10, 1, 10, -9, -11, 11, 9},
}

// NumDir is an array that represents the number of possible directions a piece can move for each piece type.
// It is utilized in the move generation function to determine the valid directions for generating moves for each piece.
//
// Explanation of Values:
// - The values in the array correspond to piece types, with indices ranging from 0 (unused) to 12 (for all piece types, including empty squares).
// - For example, NumDir[WN] indicates the number of valid directions a white knight can move.
// - The array helps control the generation of legal moves by specifying the number of directions to consider for each piece.
var NumDir = [13]int{0, 0, 8, 4, 4, 8, 8, 0, 8, 4, 4, 8, 8}

// Every time we move a piece, we will do castle_permissions &= CastlePerm[from]
// and castle_permissions &= CastlePerm[from]. The result of these operations is 1111 == 15
// except for A1, E1, H1 & A8, E8, H8
// When the rooks of the queen moves on either side, it takes out the castle permissions for that side
// eq. Black queen moves from E8 to E7. castle_permissions &= 3 -> gives 0011 -> which means BLACK side lost castling permissions
// on both queen and king side.
var CastlePerm = [120]int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

var FEN0 = "8/3q4/8/8/4Q3/8/8/8 w - - 0 2"
var FEN1 = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var FEN2 = "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
var FEN3 = "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
var FEN4 = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

var PAWNMOVES_W = "rnbqkb1r/pp1p1pPp/8/2p1pP2/1P1P4/3P3P/P1P1P3/RNBQKBNR w KQkq e6 0 1"
var PAWNMOVES_B = "rnbqkbnr/p1p1p3/3p3p/1p1p4/2P1Pp2/8/PP1P1PpP/RNBQKB1R b KQkq e3 0 1"
var KNIGHTSKINGSFEN = "5k2/1n6/4n3/6N1/8/3N4/8/5K2 w - - 0 1"
var ROOKSFEN = "6k1/8/5r2/8/1nR5/5N2/8/6K1 b - - 0 1"
var QUEENSFEN = "6k1/8/4nq2/8/1nQ5/5N2/1N6/6K1 b - - 0 1"
var BISHOPSFEN = "6k1/1b6/4n3/8/1n4B1/1B3N2/1N6/2b3K1 b - - 0 1"
var CASTLE1FEN = "r3k1r1/8/8/8/8/8/8/R3K2R w KQq - 0 1"
var CASTLE2FEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"
var PERFTFEN = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1"

var START_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func InitAll() {
	InitSq120To64()
	InitBitMasks()
	InitHashKeys()
	InitFilesRanksBrd()
}
