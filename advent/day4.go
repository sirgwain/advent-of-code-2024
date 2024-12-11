package advent

import (
	"bufio"
	"fmt"
	"os"
)

type Day4 struct {
}

func (d *Day4) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

// read day4 input as a series of lines
func (d *Day4) readInput(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var input []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())

	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return input, nil
}

func (d *Day4) part1(filename string) error {

	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	var matches []day4Match
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			matches = append(matches, findXmas("XMAS", x, y, input)...)
		}
	}

	output := make([][]rune, len(input))
	for y := 0; y < len(input); y++ {
		line := []rune(input[y])
		for x := 0; x < len(input[y]); x++ {
			line[x] = rune('.')
		}
		output[y] = line
	}

	for _, match := range matches {
		offsetX, offsetY := match.direction.offsetMultiplier()
		for i, r := range match.match {
			output[match.y+i*offsetY][match.x+i*offsetX] = r
		}
	}

	for _, line := range output {
		fmt.Println(string(line))
	}

	fmt.Printf("day4a: %s: %d\n\n", filename, len(matches))

	return nil
}

// Day4b finds X-MASes
// up
// ====
// M.M
// .A.
// S.S
//
// right
// ====
// S.M
// .A.
// S.M
//
// down
// ====
// S.S
// .A.
// M.M
//
// left
// ====
// M.S
// .A.
// M.S
func (d *Day4) part2(filename string) error {
	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	var board [][]rune = make([][]rune, 0, len(input))
	for _, line := range input {
		board = append(board, []rune(line))
	}

	numMatches := 0
	for y := 1; y < len(input)-1; y++ {
		for x := 1; x < len(input[y])-1; x++ {
			if board[y][x] == 'A' {
				// check top M, bottom S
				if board[y-1][x-1] == 'M' && board[y-1][x+1] == 'M' &&
					board[y+1][x-1] == 'S' && board[y+1][x+1] == 'S' {
					numMatches++
				}
				// check right M, left S
				if board[y-1][x-1] == 'S' && board[y-1][x+1] == 'M' &&
					board[y+1][x-1] == 'S' && board[y+1][x+1] == 'M' {
					numMatches++
				}
				// check bottom M, top S
				if board[y-1][x-1] == 'S' && board[y-1][x+1] == 'S' &&
					board[y+1][x-1] == 'M' && board[y+1][x+1] == 'M' {
					numMatches++
				}
				// check left M, right S
				if board[y-1][x-1] == 'M' && board[y-1][x+1] == 'S' &&
					board[y+1][x-1] == 'M' && board[y+1][x+1] == 'S' {
					numMatches++
				}
			}
		}
	}

	fmt.Printf("day4b: %s: %d\n\n", filename, numMatches)
	return nil
}

type day4Match struct {
	x         int
	y         int
	match     string
	direction direction
}

// horizontal, vertical, diagonal - forwards and backwards
// S . . S . . S
// . A . A . A .
// S A M X M A S
// . A . A . A .
// S . . S . . S
const possibleXMasMatches int = 8

func findXmas(matchString string, x, y int, input []string) []day4Match {
	matchLength := len(matchString)

	possibilities := make([]string, possibleXMasMatches)
	for i := 0; i < matchLength; i++ {
		// check each direction
		for dir := 0; dir < possibleXMasMatches; dir++ {
			offsetX, offsetY := direction(dir).offsetMultiplier()
			possibilities[dir] += getChar(x+i*offsetX, y+i*offsetY, input)
		}
	}

	// count numMatches
	numMatches := 0
	var matches []day4Match
	for dir, possibility := range possibilities {
		if possibility == matchString {
			numMatches++
			matches = append(matches, day4Match{
				x:         x,
				y:         y,
				match:     matchString,
				direction: direction(dir),
			})
		}
	}

	return matches
}
