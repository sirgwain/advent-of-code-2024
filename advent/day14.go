package advent

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2024/advent/color"
	"github.com/sirgwain/advent-of-code-2024/advent/tui"
)

type Day14 struct {
	*Options
}

type day14Input = day14Board

type day14robot struct {
	position
	velocity position
}

type day14Board struct {
	width  int
	height int
	robots []day14robot

	seconds       int
	confidence    int
	midUpperLeft  position
	midLowerRight position
}

const robotChar = '☹'

var robotRedRender = robotGreenStyle.Render(string(robotChar))
var robotGreenRender = robotRedStyle.Render(string(robotChar))
var midRender = lipgloss.NewStyle().Background(color.DarkBlue).Render(".")

// Run is the main entry point for a day. It reads the input file and runs the part
func (d *Day14) Run(part int, filename string, opts ...Option) error {
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

func (d *Day14) readInput(filename string) (day14Input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return day14Input{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	boardPattern := `w=(\d+),h=(\d+)`
	reBoard, err := regexp.Compile(boardPattern)
	if err != nil {
		return day14Input{}, fmt.Errorf("bad regex %s %w", boardPattern, err)
	}

	robotPattern := `p=(\d+),(\d+) v=([-\d]+),([-\d]+)`
	reRobot, err := regexp.Compile(robotPattern)
	if err != nil {
		return day14Input{}, fmt.Errorf("bad regex %s %w", robotPattern, err)
	}

	var input day14Input

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if i == 0 {
			match := reBoard.FindSubmatch(line)
			input.width = mustAtoi(match[1])
			input.height = mustAtoi(match[2])
			i++
			continue
		}
		match := reRobot.FindSubmatch(line)
		r := day14robot{
			position: position{mustAtoi(match[1]), mustAtoi(match[2])},
			velocity: position{mustAtoi(match[3]), mustAtoi(match[4])},
		}
		input.robots = append(input.robots, r)
	}
	return input, nil
}

// RunVisual is like Run, but it starts a bubbletea program and runs the solver in a goroutine
func (d *Day14) RunVisual(part int, filename string, opts ...Option) error {
	d.Options = newRun(opts...)
	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	switch part {
	case 1:
		fallthrough
	case 2:
		return d.part2Visual(input)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day14) part1(input day14Input) error {
	solution := 0
	board := input
	board.midLowerRight = position{board.width - 1, board.height - 1}
	for range 100 {
		// fmt.Printf("%s\n\n", input.view())
		board.move(1)
	}

	fmt.Printf("\nFinal:\n%s\n\n", board.view())
	solution = board.safetyFactor()
	fmt.Printf("\nSolution: %s\n", solutionStyle.Render(strconv.Itoa(solution)))
	return nil
}

func (d *Day14) part2(input day14Input) error {
	board := input
	board.midUpperLeft, board.midLowerRight = board.treeArea()
	fmt.Printf("%s\n", board.view())

	for {
		board.move(1)
		board.seconds++
		if board.tree() > .5 {
			board.confidence = int(board.tree() * 100)
			break
		}
	}
	fmt.Println(board.view())
	fmt.Println(board.viewSolution())

	return nil
}

func (d *Day14) part2Visual(input day14Input) error {
	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 14"))

	board := input
	board.midUpperLeft, board.midLowerRight = board.treeArea()

	// run the solver in a gouroutine and Send a message to the bubbletea program to update the viewport
	// on each step
	go func() {
		view := board.view()
		seconds := 0
		for {
			board.move(1)
			seconds++
			treeConfidence := board.tree()
			if treeConfidence > .5 {
				view = board.view()
				board.seconds = seconds
				board.confidence = int(treeConfidence * 100)
			}
			content := fmt.Sprintf("%s\n%s - current: %d (%d)", view, board.viewSolution(), seconds, int(treeConfidence*100))
			p.Send(tui.UpdateViewport(content, board.width))
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

func (b day14Board) view() string {
	robots := make(map[position]int, len(b.robots))
	for _, r := range b.robots {
		robots[r.position]++
	}
	var sb strings.Builder

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if count, ok := robots[position{x, y}]; ok {
				if count > 1 {
					sb.WriteString(robotRedRender)
				} else {
					sb.WriteString(robotGreenRender)
				}
			} else {
				// draw the boundary we think the tree might be in
				if ((x == b.midUpperLeft.x || x == b.midLowerRight.x) && y > b.midUpperLeft.y && y < b.midLowerRight.y) ||
					((y == b.midUpperLeft.y || y == b.midLowerRight.y) && x > b.midUpperLeft.x && x < b.midLowerRight.x) {
					sb.WriteString(midRender)
				} else {
					sb.WriteRune('.')
				}
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
func (b day14Board) viewSolution() string {
	return fmt.Sprintf("\nseconds: %s - %s%%", solutionStyle.Render(strconv.Itoa(b.seconds)), solutionStyle.Render(strconv.Itoa(b.confidence)))
}

// guess that the tree is in the middle 1/3rd
func (b day14Board) treeArea() (midUpperLeft, midLowerRight position) {
	midLeft := b.width / 3
	midRight := b.width - midLeft
	midTop := b.height / 3
	midBottom := b.height - midTop

	return position{midLeft, midTop}, position{midRight, midBottom}
}

// return what percent of robots are in that square
// it's about a 1200sq area and 500 total robots, so if we have a large clump there, it's probably a tree?
func (b day14Board) tree() float64 {

	midRobots := 0
	for _, r := range b.robots {
		if r.x >= b.midUpperLeft.x && r.x < b.midLowerRight.x && r.y >= b.midUpperLeft.y && r.y < b.midLowerRight.y {
			midRobots++
		}
	}

	return float64(midRobots) / float64(len(b.robots))
}
func (b *day14Board) move(moves int) {
	for i := range b.robots {
		r := &b.robots[i]
		x, y := r.x+(r.velocity.x*moves), r.y+(r.velocity.y*moves)
		// account for loop arounds
		x = x % b.width
		y = y % b.height
		if x < 0 {
			x = b.width + x
		}
		if y < 0 {
			y = b.height + y
		}
		// move the robot
		r.x, r.y = x, y
	}
}

func (b *day14Board) safetyFactor() int {
	// clockwise upper left to lower left quadrant
	q1, q2, q3, q4 := b.quadrants()

	return q1 * q2 * q3 * q4
}

// clockwise upper left to lower left quadrant
func (b *day14Board) quadrants() (q1, q2, q3, q4 int) {
	mid := position{b.width / 2, b.height / 2}
	for i, r := range b.robots {
		_ = i
		if r.x == mid.x || r.y == mid.y {
			// doesn't count
			continue
		}
		if r.x < mid.x {
			if r.y < mid.y {
				// slog.Debug(fmt.Sprintf("r%d (%v) in upper left", i, r.position))
				q1++
			} else {
				// slog.Debug(fmt.Sprintf("r%d (%v) in lower left", i, r.position))
				q4++
			}
		} else {
			if r.y < mid.y {
				// slog.Debug(fmt.Sprintf("r%d (%v) in upper right", i, r.position))
				q2++
			} else {
				// slog.Debug(fmt.Sprintf("r%d (%v) in lower right", i, r.position))
				q3++
			}
		}
	}

	slog.Debug(fmt.Sprintf("safety: q1: %d, q2: %d, q3: %d, q4: %d", q1, q2, q3, q4))
	return q1, q2, q3, q4
}
