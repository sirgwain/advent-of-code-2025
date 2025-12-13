package advent

func ValidPosition(p Point, width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

func MakeBoard[T any](width, height int) [][]T {
	board := make([][]T, height)
	for y := range height {
		board[y] = make([]T, width)
	}
	return board
}

// GetBoardValue returns a rune/int/bool at x,y in the input or the empty value if out of bounds
func GetBoardValue[T int | uint | byte | rune | bool](x, y int, board [][]T) T {
	var zero T
	if y < 0 || y >= len(board) {
		return zero
	}
	if x < 0 || x >= len(board[y]) {
		return zero
	}
	return board[y][x]
}

// find
func FindValue[T int | uint | byte | rune | bool](board [][]T, c T) (x, y int) {
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] == c {
				return x, y
			}
		}
	}
	return 0, 0
}
