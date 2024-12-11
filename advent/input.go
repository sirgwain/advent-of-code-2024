package advent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// read input as a series of rune lines
func readInputAsRunes(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var input [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		input = append(input, []rune(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return input, nil
}

// read input as a series of rune lines
func readInputAsIntBoard(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var input [][]int
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, make([]int, len(line)))
		for i, r := range []rune(line) {
			input[lineNum][i], err = strconv.Atoi(string(r))
			if err != nil {
				return nil, fmt.Errorf("%v is not a number %w", r, err)
			}
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return input, nil
}
