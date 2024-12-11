package advent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sirgwain/advent-of-code-2024/advent/tui"
)

type Day6 struct {
	Delay         int
	UpdateUIMoves int
}

type day6Board struct {
	position
	board         [][]rune
	direction     direction
	vistedSquares int
	obstaclesHit  map[positionDirection]int
	complete      bool
	cycle         bool
	onMove        func()
}

func (d *Day6) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day6) part1(filename string) error {

	input, err := readInputAsRunes(filename)
	if err != nil {
		return err
	}

	x, y := findValue(input, '^')
	board := day6Board{board: input, position: position{x: x, y: y}, direction: directionUp, vistedSquares: 1, obstaclesHit: make(map[positionDirection]int)}

	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 6 - Part 1"))

	width := len(board.board[0])
	go func() {
		count := 0
		board.onMove = func() {
			// update the UI every 10th call
			count++
			if count > 10 {
				p.Send(tui.UpdateViewport(board.boardView(), width))
				count = 0
			}
		}
		board.runBoard()

		// all done
		p.Send(tea.QuitMsg{})

	}()
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	if board.complete {
		fmt.Printf("%s\n\n", board.boardView())
		fmt.Printf("Total Steps: %s\n", solutionStyle.Render(strconv.Itoa(board.vistedSquares)))
	}

	return nil
}

func (d *Day6) part2(filename string) error {
	input, err := readInputAsRunes(filename)
	if err != nil {
		return err
	}

	startX, startY := findValue(input, '^')
	board := day6Board{board: input, position: position{x: startX, y: startY}, direction: directionUp, vistedSquares: 1, obstaclesHit: make(map[positionDirection]int)}
	initialRun := board.duplicate()
	initialRun.runBoard()

	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 6 - Part 2"))

	// add an obstacle in every visited square except the start
	obstacles := make([]position, 0, initialRun.vistedSquares-1)

	for y := 0; y < len(initialRun.board); y++ {
		for x := 0; x < len(initialRun.board[y]); x++ {
			if y == startY && x == startX {
				continue
			}
			if initialRun.board[y][x] != '.' && initialRun.board[y][x] != '#' {
				// path, add an obstacle here
				obstacles = append(obstacles, position{x, y})
			}
		}
	}

	width := len(board.board[0])
	cycleBoards := make([]day6Board, 0)
	go func() {
		for _, obstacle := range obstacles {
			testBoard := board.duplicate()
			testBoard.board[obstacle.y][obstacle.x] = '#'

			done := make(chan struct{})
			// Start the goroutine
			go func() {
				count := 0
				testBoard.onMove = func() {
					// update the UI every 10th call
					if d.UpdateUIMoves != 0 {
						count++
						if count > d.UpdateUIMoves {
							p.Send(tui.UpdateViewport(testBoard.boardView(), width))
							count = 0
						}
					}

					if d.Delay != 0 {
						time.Sleep(time.Duration(d.Delay * int(time.Millisecond)))
					}
				}
				testBoard.runBoard()
				close(done)
			}()
			// Wait for run to finish
			<-done

			if testBoard.cycle {
				testBoard.board[obstacle.y][obstacle.x] = 'O'
				cycleBoards = append(cycleBoards, testBoard)
			}
		}

		// signal we are done
		p.Send(tea.Quit())
	}()

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	if d.UpdateUIMoves != 0 {
		for _, b := range cycleBoards {
			fmt.Printf("%s\n", b.boardView())
		}
	}

	fmt.Printf("Total Cycles: %s\n", solutionStyle.Render(strconv.Itoa(len(cycleBoards))))

	return nil
}

func (b *day6Board) duplicate() day6Board {
	dup := *b
	dup.obstaclesHit = make(map[positionDirection]int)
	dup.board = duplicate2DSlice(b.board)
	return dup
}

func (b *day6Board) runBoard() {
	for {
		b.moveGuard()
		if b.onMove != nil {
			b.onMove()
		}
		if b.cycle {
			break
		}
		if b.complete {
			break
		}
	}
}

func (b *day6Board) moveGuard() {
	var x, y int

	switch b.direction {
	case directionUp:
		y = -1
	case directionRight:
		x = 1
	case directionDown:
		y = 1
	case directionLeft:
		x = -1
	}

	nextPosition := position{
		x: b.x + x,
		y: b.y + y,
	}

	if nextPosition.y < 0 || nextPosition.y == len(b.board) {
		// at edge, no more moves
		b.complete = true
		return
	}

	if nextPosition.x < 0 || nextPosition.x == len(b.board[b.y]) {
		// at edge, no more moves
		b.complete = true
		return
	}

	nextSquare := getBoardValue(nextPosition.x, nextPosition.y, b.board)

	if nextSquare == '#' {
		// obstacle, turn
		key := positionDirection{
			position:  nextPosition,
			direction: b.direction,
		}
		if hits := b.obstaclesHit[key]; hits > 0 {
			// we've hit this obstacle from this direction before, we're in a loop
			b.cycle = true
			return
		}
		b.obstaclesHit[key]++
		b.direction = b.direction.turnRight()
		return
	}

	// check if we have been here before
	if nextSquare == '.' {
		b.vistedSquares = b.vistedSquares + 1
	}

	// move the guardq
	b.board[b.y][b.x] = 'X'
	b.y = b.y + y
	b.x = b.x + x
	b.board[b.y][b.x] = b.direction.getChar()
}

// prerender some styled characters
var (
	renderedPath     = pathStyle.Render("X")
	renderedObstacle = obstacleStyle.Render("#")
)

func (b *day6Board) boardView() string {
	guard := b.direction.getChar()
	renderedGuard := guardStyle.Render(string(guard))

	var sb strings.Builder
	for _, line := range b.board {
		for _, r := range line {
			switch r {
			case 'X':
				sb.WriteString(renderedPath)
			case '#':
				sb.WriteString(renderedObstacle)
			case guard:
				sb.WriteString(renderedGuard)
			default:
				sb.WriteRune(r)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
