package advent

// duplicate2DSlice creates a deep copy of a 2D slice.
func duplicate2DSlice[T any](original [][]T) [][]T {
	if original == nil {
		return nil
	}

	// Create a new 2D slice with the same dimensions as the original
	duplicate := make([][]T, len(original))
	for i := range original {
		duplicate[i] = make([]T, len(original[i]))
		copy(duplicate[i], original[i])
	}
	return duplicate
}

// thanks chatgpt!
func removeIndex(original []int, index int) []int {
	if index < 0 || index >= len(original) {
		// Return a copy of the original slice if the index is out of range
		return append([]int(nil), original...)
	}

	// Create a new slice excluding the specified index
	newSlice := make([]int, 0, len(original)-1)
	newSlice = append(newSlice, original[:index]...)
	newSlice = append(newSlice, original[index+1:]...)
	return newSlice
}

// getChar returns a char at x,y in the input or "" if no char is present
func getChar(x, y int, input []string) string {
	if y < 0 || y >= len(input) {
		return ""
	}
	if x < 0 || x >= len(input[y]) {
		return ""
	}
	return string(input[y][x])
}
