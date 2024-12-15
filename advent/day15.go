package advent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirgwain/advent-of-code-2024/advent/tui"
)

type Day15 struct {
	*Options
}

type day15Input = day15Board

type day15Type = rune

const (
	empty day15Type = '.'
	wall  day15Type = '#'
	box   day15Type = 'O'
	robot day15Type = '@'
)

type day15Board struct {
	width  int
	height int
	board  [][]rune
	robot  position
	moves  []direction
	onStep func()

	move     int
	solution int
}

var robotRender = robotStyle.Render("@")
var wallRender = wallStyle.Render("#")
var boxRender = boxStyle.Render("O")

// Run is the main entry point for a day. It reads the input file and runs the part
func (d *Day15) Run(part int, filename string, opts ...Option) error {
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

func (d *Day15) readInput(filename string) (day15Input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return day15Input{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var input day15Input

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		// skip empties
		if strings.TrimSpace(string(line)) == "" {
			continue
		}
		if line[0] == '#' { // wall means board, read as runes
			input.board = append(input.board, []rune(string(line)))
		}
		// everything else is directions
		for _, r := range line {
			switch rune(r) {
			case '>':
				input.moves = append(input.moves, directionRight)
			case '<':
				input.moves = append(input.moves, directionLeft)
			case '^':
				input.moves = append(input.moves, directionUp)
			case 'v':
				input.moves = append(input.moves, directionDown)
			}
		}
	}

	input.height = len(input.board)
	input.width = len(input.board[0])
	input.robot = input.robotPosition()
	return input, nil
}

// RunVisual is like Run, but it starts a bubbletea program and runs the solver in a goroutine
func (d *Day15) RunVisual(part int, filename string, opts ...Option) error {
	d.Options = newRun(opts...)
	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	switch part {
	case 1:
		fallthrough
	case 2:
		return d.visual(input)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day15) part1(input day15Input) error {
	board := input

	board.solve()

	fmt.Printf("\n%s\n\n", board.view())
	fmt.Printf("%s\n", board.viewSolution())
	return nil
}

func (d *Day15) part2(input day15Input) error {
	board := input
	fmt.Println(board.view())
	fmt.Println(board.viewSolution())

	return nil
}

func (d *Day15) visual(input day15Input) error {
	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 15"))

	board := input

	// run the solver in a gouroutine and Send a message to the bubbletea program to update the viewport
	// on each step
	go func() {
		for {

			// update the UI
			board.onStep = func() {
				content := fmt.Sprintf("%s\n%s", board.view(), board.viewSolution())
				p.Send(tui.UpdateViewport(content, board.width))

				if d.Options.Delay > 0 {
					time.Sleep(time.Millisecond * time.Duration(d.Options.Delay))
				}

			}

			board.solve()
		}
	}()

	// execute the bubbletea program. This will block until the user pressed q or esc
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	// output the board and final result before exiting the program
	fmt.Println(board.view())
	fmt.Println(board.viewSolution())
	return nil
}

func (b day15Board) robotPosition() position {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			switch b.board[y][x] {
			case '@':
				return position{x, y}
			}
		}
	}
	return position{}
}

func (b day15Board) gps() int {
	solution := 0
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			switch b.board[y][x] {
			case 'O':
				solution += 100*y + x
			}
		}
	}
	return solution
}

func (b *day15Board) solve() {
	for _, dir := range b.moves {
		b.moveRobot(dir)
		b.solution = b.gps()
		if b.onStep != nil {
			b.onStep()
		}
	}

}

func (b *day15Board) moveRobot(dir direction) {
	x, y := b.robot.x, b.robot.y
	dx, dy := dir.offsetMultiplier()
	tx, ty := x+dx, y+dy
	target := getBoardValue(tx, ty, b.board)
	switch target {
	case wall: // wall, do nothing
		return
	case empty: // empty space, move
		b.board[ty][tx] = robot
		b.board[y][x] = empty
		b.robot.x, b.robot.y = tx, ty
	case box: // box, try and move it
		if b.moveBox(tx, ty, dir) {
			b.board[ty][tx] = robot
			b.board[y][x] = empty
			b.robot.x, b.robot.y = tx, ty
		}
	}
	b.move++
}

func (b *day15Board) moveBox(x, y int, dir direction) bool {
	dx, dy := dir.offsetMultiplier()
	tx, ty := x+dx, y+dy
	target := getBoardValue(tx, ty, b.board)
	switch target {
	case wall: // wall, do nothing
		return false
	case empty: // empty space, move
		b.board[ty][tx] = box
		b.board[y][x] = empty
		return true
	case box: // box, try and move it
		if b.moveBox(tx, ty, dir) {
			b.board[ty][tx] = box
			b.board[y][x] = empty
			return true
		}
	}
	return false
}

func (b day15Board) view() string {
	var sb strings.Builder

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			switch b.board[y][x] {
			case wall:
				sb.WriteString(wallRender)
			case box:
				sb.WriteString(boxRender)
			case robot:
				sb.WriteString(robotRender)
			default:
				sb.WriteRune(empty)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
func (b day15Board) viewSolution() string {
	return fmt.Sprintf("Move %d, Solution: %s", b.move, solutionStyle.Render(strconv.Itoa(b.solution)))
}
