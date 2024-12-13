package advent

func validPosition(p position, width, height int) bool {
	return p.x >= 0 && p.x < width && p.y >= 0 && p.y < height
}

// getBoardValue returns a rune/int/bool at x,y in the input or the empty value if out of bounds
func getBoardValue[T int | uint | rune | bool](x, y int, board [][]T) T {
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
func findValue[T int | uint | rune | bool](board [][]T, c T) (x, y int) {
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if board[y][x] == c {
				return x, y
			}
		}
	}
	return 0, 0
}
