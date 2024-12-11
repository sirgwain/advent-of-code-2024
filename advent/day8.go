package advent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirgwain/advent-of-code-2024/advent/tui"
)

type Day8 struct {
	*Options
}

type day8Board struct {
	board           [][]rune
	antennas        map[rune][]position
	antinodes       map[position]bool
	onAntinodeFound func(p position)
	solution        int
}

func (d *Day8) Run(part int, filename string, opts ...Option) error {
	d.Options = newRun(opts...)
	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	switch part {
	case 1:
		return d.part1(input)
	case 2:
		return d.part2(input)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day8) RunVisual(part int, filename string, opts ...Option) error {
	d.Options = newRun(opts...)
	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	switch part {
	case 2:
		return d.part2Visual(input)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day8) readInput(filename string) ([][]rune, error) {
	return readInputAsRunes(filename)
}

func (d *Day8) part1(input [][]rune) error {

	board := day8Board{board: input}
	board.findAntinodes()

	fmt.Printf("%s\n", board.view())
	fmt.Printf("Antenna Types: %s, valid Antinodes: %s\n",
		correctResultStyle.Render(strconv.Itoa(len(board.antennas))),
		solutionStyle.Render(strconv.Itoa(len(board.antinodes))))

	return nil
}

func (d *Day8) part2(input [][]rune) error {
	board := day8Board{board: input}
	board.findAntinodesWithResonance()

	// fmt.Printf("%s\n", board.view())
	fmt.Printf("Antenna Types: %s, valid Antinodes: %s\n",
		correctResultStyle.Render(strconv.Itoa(len(board.antennas))),
		solutionStyle.Render(strconv.Itoa(len(board.antinodes))))

	return nil
}

func (d *Day8) part2Visual(input [][]rune) error {
	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 8 - Part 2"))

	// find the solution so we can hide it from the output
	silentBoard := day8Board{board: duplicate2DSlice(input)}
	silentBoard.findAntinodesWithResonance()

	board := day8Board{board: input, solution: len(silentBoard.antinodes)}
	width := len(board.board[0])
	go func() {
		// update the UI
		board.onAntinodeFound = func(pos position) {
			content := fmt.Sprintf("%s\n%s", board.view(), board.viewSolution(d.RedactSolution))
			p.Send(tui.UpdateViewport(content, width))
			if d.Delay != 0 {
				time.Sleep(time.Millisecond * time.Duration(d.Delay))
			}
		}
		board.findAntinodesWithResonance()
	}()

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	fmt.Printf("%s\n", board.view())
	fmt.Printf("Antenna Types: %s, valid Antinodes: %s\n",
		correctResultStyle.Render(strconv.Itoa(len(board.antennas))),
		solutionStyle.Render(strconv.Itoa(len(board.antinodes))))

	return nil
}

// find the antinode for two positions
// ..........
// ...#......  3,1 a1
// ..........
// ....a.....  4,3
// ..........
// .....a....  5,5
// ..........
// ......#...  6,7 a2
//
// in the above example, dist between p1, p2
// is -1,-2 so the antinode is
// a1 = p2+(-1,-2)*2
// a2 = p1-(-1,-2)*2
func findAntinodePositions(p1 position, p2 position) (a1 position, a2 position) {
	xdist := p1.x - p2.x
	ydist := p1.y - p2.y

	a1 = position{p2.x + xdist*2, p2.y + ydist*2}
	a2 = position{p1.x - xdist*2, p1.y - ydist*2}

	return a1, a2
}

func findAntinodeLinePoints(p1 position, p2 position, width, height int) []position {
	xdist := p1.x - p2.x
	ydist := p1.y - p2.y

	linePoints := []position{p1, p2}
	// do the y=mx+b forward
	pos := p1
	for {
		pos = position{pos.x + xdist, pos.y + ydist}
		if !validPosition(pos, width, height) {
			break
		}
		linePoints = append(linePoints, pos)
	}

	// now backward
	pos = p2
	for {
		pos = position{pos.x - xdist, pos.y - ydist}
		if !validPosition(pos, width, height) {
			break
		}
		linePoints = append(linePoints, pos)
	}

	return linePoints
}

func (b *day8Board) findAntinodes() {
	height, width := len(b.board), len(b.board[0])

	antennas := make(map[rune][]position)
	antinodes := make(map[position]bool)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			antenna := b.board[y][x]
			if antenna == '.' {
				continue
			}
			// found an antenna
			p1 := position{x, y}
			for _, p2 := range antennas[antenna] {
				a1, a2 := findAntinodePositions(p1, p2)
				if validPosition(a1, width, height) {
					antinodes[a1] = true
					if b.onAntinodeFound != nil {
						b.onAntinodeFound(a1)
					}

				}
				if validPosition(a2, width, height) {
					antinodes[a2] = true
					if b.onAntinodeFound != nil {
						b.onAntinodeFound(a2)
					}
				}
			}
			antennas[antenna] = append(antennas[antenna], p1)
		}
	}

	b.antennas = antennas
	b.antinodes = antinodes
}

func (b *day8Board) findAntinodesWithResonance() {
	height, width := len(b.board), len(b.board[0])

	b.antennas = make(map[rune][]position)
	b.antinodes = make(map[position]bool)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			antenna := b.board[y][x]
			if antenna == '.' {
				continue
			}
			// found an antenna
			p1 := position{x, y}
			for _, p2 := range b.antennas[antenna] {
				points := findAntinodeLinePoints(p1, p2, width, height)
				for _, point := range points {
					b.antinodes[point] = true
					if b.onAntinodeFound != nil {
						b.onAntinodeFound(point)
					}
				}
			}
			b.antennas[antenna] = append(b.antennas[antenna], p1)
		}
	}
}

// prerender some styled characters
var (
	renderedAntinode = antinodeStyle.Render("#")
)

func (b *day8Board) view() string {
	var sb strings.Builder
	for y, line := range b.board {
		for x, r := range line {
			switch r {
			case '.':
				// if there is an antinode here, render it
				if b.antinodes[position{x, y}] {
					sb.WriteString(renderedAntinode)
				} else {
					sb.WriteRune(r)
				}
			default:
				if b.antinodes[position{x, y}] {
					sb.WriteString(antennaWithAntinodeStyle.Render(string(r)))
				} else {
					sb.WriteString(antennaStyle.Render(string(r)))
				}
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *day8Board) viewSolution(redactSolution bool) string {
	// hide the solution if it's within 100
	solution := solutionStyle.Render(strconv.Itoa(len(b.antinodes)))
	if redactSolution && b.solution-len(b.antinodes) < 100 {
		solution = solutionStyle.Render("<redacted>")
	}

	return fmt.Sprintf("Antenna Types: %s, valid Antinodes: %s\n",
		correctResultStyle.Render(strconv.Itoa(len(b.antennas))),
		solution)

}
