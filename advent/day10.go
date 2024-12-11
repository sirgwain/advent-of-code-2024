package advent

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2024/advent/tui"
)

type Day10 struct {
	*Options
}

type day10Board struct {
	board     [][]int
	width     int
	height    int
	trailEnd  int
	visited   [][]bool
	solutions map[position]map[position]int
	solution1 int
	solution2 int
	onStep    func(p position)
}

// base height styles before rendering the board. These start at dark blue and get lighter
var heightRenders = []string{
	lipgloss.NewStyle().Foreground(lipgloss.Color("30")).Render("0"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("31")).Render("1"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("32")).Render("2"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("33")).Render("3"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render("4"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("35")).Render("5"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Render("6"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("37")).Render("7"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("38")).Render("8"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("9"),
}

// styles for visited locations, these are pink and orange
var heightVisitedRenders = []string{
	lipgloss.NewStyle().Foreground(lipgloss.Color("200")).Render("0"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("201")).Render("1"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Render("2"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Render("3"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Render("4"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("5"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("206")).Render("6"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("207")).Render("7"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Render("8"),
	lipgloss.NewStyle().Foreground(lipgloss.Color("209")).Render("9"),
}

// Run is the main entry point for a day. It reads the input file and runs the part
func (d *Day10) Run(part int, filename string, opts ...Option) error {
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

// RunVisual is like Run, but it starts a bubbletea program and runs the solver in a goroutine
func (d *Day10) RunVisual(part int, filename string, opts ...Option) error {
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

func (d *Day10) readInput(filename string) ([][]int, error) {
	return readInputAsIntBoard(filename)
}

func (d *Day10) part1(input [][]int) error {
	board := day10Board{
		board:    input,
		trailEnd: 9,
		height:   len(input),
		width:    len(input[0]),
	}
	fmt.Println(board.view())
	board.findTrails()
	fmt.Println(board.viewSolution(false))

	return nil
}

func (d *Day10) part2(input [][]int) error {
	// part1 and part2 are the same for this day
	d.part1(input)
	return nil
}

func (d *Day10) part2Visual(input [][]int) error {
	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 10"))

	// find the solution so we can hide it from the output
	silentBoard := day10Board{board: duplicate2DSlice(input),
		trailEnd: 9,
		height:   len(input),
		width:    len(input[0]),
	}
	silentBoard.findTrails()

	// create a board with the solution numbers there already so we can redact them as we get close
	board := day10Board{board: input,
		trailEnd:  9,
		height:    len(input),
		width:     len(input[0]),
		solution1: silentBoard.trails(),
		solution2: silentBoard.distinctTrails(),
	}

	// run the solver in a gouroutine and Send a message to the bubbletea program to update the viewport
	// on each step
	width := len(board.board[0])
	go func() {
		// update the UI
		board.onStep = func(pos position) {
			content := fmt.Sprintf("%s\n%s", board.view(), board.viewSolution(d.RedactSolution))
			p.Send(tui.UpdateViewport(content, width))
		}
		board.findTrails()
	}()

	// execute the bubbletea program. This will block until the user pressed q or esc
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	// output the board and final result before exiting the program
	fmt.Println(board.view())
	fmt.Println(board.viewSolution(false))
	return nil
}

// find all trails starting at 0 and ending at 9
func (b *day10Board) findTrails() {
	// the solutions map contains starting posistions and a map of ending positions with
	// a count for each time we arrive
	b.solutions = map[position]map[position]int{}

	// keep track of what positions we've visited just for UI updates to look pretty
	b.visited = make([][]bool, b.height)
	for y := 0; y < b.height; y++ {
		b.visited[y] = make([]bool, b.width)
	}

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			// for every board value of 0, recursively find all trails that end at 9
			if getBoardValue(x, y, b.board) == 0 {
				start := position{x, y}
				trailEnds := make(map[position]int)

				// recursively find the trail
				b.findTrail(start, 0, trailEnds)

				// if our findTrail came with solutions, add them to our solutions map
				// for this starting position
				if len(trailEnds) > 0 {
					b.solutions[start] = trailEnds
				}
			}
		}
	}
}

// findTrail will search for the next height
func (b day10Board) findTrail(pos position, currentHeight int, trailHeads map[position]int) {
	b.visited[pos.y][pos.x] = true
	if b.onStep != nil {
		b.onStep(pos)
	}
	if currentHeight == b.trailEnd {
		trailHeads[pos]++
		return
	}

	// try each four positions
	for _, dir := range cardinalDirections {
		testPos := pos.addDirection(dir)
		if validPosition(testPos, b.width, b.height) && getBoardValue(testPos.x, testPos.y, b.board) == currentHeight+1 {
			b.findTrail(testPos, currentHeight+1, trailHeads)
		}
	}
}

// view will render the board as a string.
// This is called by the UI for every step to update
// it will also be called at the end of the program to display the final board
func (b day10Board) view() string {
	var sb strings.Builder
	board := b.board
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			// we have 10 height styles of increasing color, one slice for visited positions and one for unvisited
			// render each number as a style
			if b.visited != nil && b.visited[y][x] {
				sb.WriteString(heightVisitedRenders[board[y][x]])
			} else {
				sb.WriteString(heightRenders[board[y][x]])
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// count the total number of trails, for part 1
func (b day10Board) trails() int {
	total := 0
	for _, solution := range b.solutions {
		total += len(solution)
	}

	return total
}

// count the number of distinct trails, for part 2
func (b day10Board) distinctTrails() int {
	distinct := 0
	for _, solution := range b.solutions {
		for _, paths := range solution {
			distinct += paths
		}
	}

	return distinct
}

// viewSolution renders the solution for part 1 and 2 as a string
func (b day10Board) viewSolution(redactSolution bool) string {
	// hide the solution if it's within 100
	totalTrails := b.trails()
	distinctTrails := b.distinctTrails()
	trails := solutionStyle.Render(strconv.Itoa(totalTrails))
	if redactSolution && float64(float64(totalTrails)/float64(b.solution1)) > .75 {
		trails = solutionStyle.Render("<redacted>")
	}
	distinctPaths := solutionStyle.Render(strconv.Itoa(distinctTrails))
	if redactSolution && float64(float64(distinctTrails)/float64(b.solution2)) > .75 {
		distinctPaths = solutionStyle.Render("<redacted>")
	}

	return fmt.Sprintf("Trails: %s, Distinct Paths: %s\n", trails, distinctPaths)
}
