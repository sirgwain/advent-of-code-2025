package advent

func ValidPosition(p Position, width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

// GetBoardValue returns a rune/int/bool at x,y in the input or the empty value if out of bounds
func GetBoardValue[T int | uint | rune | bool](x, y int, board [][]T) T {
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
func FindValue[T int | uint | rune | bool](board [][]T, c T) (x, y int) {
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] == c {
				return x, y
			}
		}
	}
	return 0, 0
}
