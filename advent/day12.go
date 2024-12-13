package advent

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code-2024/advent/tui"
)

type Day12 struct {
	*Options
}

type day12Input [][]rune

type day12Board struct {
	board  [][]rune
	width  int
	height int

	foundEdges [][]uint
	regions    []day12Region
	solution   int
	onStep     func()
}

type day12Region struct {
	plotType rune
	area     int
	sides    int
}

// Run is the main entry point for a day. It reads the input file and runs the part
func (d *Day12) Run(part int, filename string, opts ...Option) error {
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
func (d *Day12) RunVisual(part int, filename string, opts ...Option) error {
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

func (d *Day12) readInput(filename string) (day12Input, error) {
	return readInputAsRunes(filename)
}

func (d *Day12) part1(input day12Input) error {

	solution := 0

	board := duplicate2DSlice(input)
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			plotType := board[y][x]
			if unicode.IsLower(plotType) {
				// we visted this rune already
				continue
			}
			// find all plots this type
			area, perimeter := d.findRegion(plotType, board, x, y)

			fmt.Printf("Found area for %s, area: %d, perimeter %d\n", string(plotType), area, perimeter)
			// add the cost to our solution
			solution += area * perimeter
		}
	}

	fmt.Printf("\nSolution: %s\n", solutionStyle.Render(strconv.Itoa(solution)))
	return nil
}

// regions look like this
// +-+-+-+-+
// |A A A A|
// +-+-+-+-+
// but can also go up/down and form funny shapes
func (d *Day12) findRegion(plotType rune, board [][]rune, x, y int) (area, perimeter int) {
	// mark this path as visited
	visited := unicode.ToLower(plotType)
	board[y][x] = visited

	area = 1
	for _, dir := range cardinalDirections {
		offsetX, offsetY := dir.offsetMultiplier()
		test := getBoardValue(x+offsetX, y+offsetY, board)

		if test == visited {
			continue
		}
		if test != plotType {
			// we found an edge
			perimeter++
			continue
		}

		// a neighbor is the same plotType, add its area and perimiter to ours
		a, p := d.findRegion(plotType, board, x+offsetX, y+offsetY)
		area += a
		perimeter += p
	}

	// clear out this board entry
	return area, perimeter
}

func (d *Day12) part2(input day12Input) error {
	board := day12Board{
		board:  input,
		height: len(input),
		width:  len(input[0]),
	}
	board.findPlots()
	fmt.Printf("%s\n%s\n%s", board.view(), board.viewRegions(), board.viewSolution())
	return nil
}

func (d *Day12) part2Visual(input [][]rune) error {
	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 12 - Part 2").WithMinWidth(100))

	board := day12Board{
		board:  input,
		height: len(input),
		width:  len(input[0]),
	}
	// run the solver in a gouroutine and Send a message to the bubbletea program to update the viewport
	// on each step
	width := len(board.board[0])
	go func() {
		// update the UI
		board.onStep = func() {
			content := fmt.Sprintf("%s\n%s\n%s", board.view(), board.viewRegions(), board.viewSolution())
			p.Send(tui.UpdateViewport(content, width))
			if d.Options.Delay > 0 {
				time.Sleep(time.Millisecond * time.Duration(d.Options.Delay))
			}
		}
		board.findPlots()
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

func (b *day12Board) findPlots() {
	board := b.board
	b.foundEdges = make2DSlice(len(board[0]), len(board))

	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			plotType := board[y][x]
			if unicode.IsLower(plotType) {
				// we visted this rune already
				continue
			}
			// find all plots this type
			b.regions = append(b.regions, day12Region{plotType: plotType})
			region := &b.regions[len(b.regions)-1]
			area, sides := b.findSidedRegion(region, x, y)

			// fmt.Printf("Found area qfor %s, area: %d, sides %d = %d\n", string(plotType), area, sides, area*sides)
			// add the cost to our solution
			b.solution += area * sides
			if b.onStep != nil {
				b.onStep()
			}
		}
	}

}

func (b *day12Board) view() string {
	var sb strings.Builder
	board := b.board
	for y := 0; y < len(board); y++ {
		for x := 0; x < len(board[y]); x++ {
			if unicode.IsLower(board[y][x]) {
				// a visited square, start with color "200"
				color := 200 + int(board[y][x]) - int('a')
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(color))).Render(string(board[y][x])))
			} else {
				// a non visted square, start with color "30"
				color := 30 + int(board[y][x]) - int('a')
				sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(color))).Render(string(board[y][x])))
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *day12Board) viewRegions() string {
	var sb strings.Builder

	for _, region := range b.regions {
		color := 200 + int(region.plotType) - int('A')
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(strconv.Itoa(color))).Render(string(region.plotType)))
		sb.WriteString(fmt.Sprintf(" Area: %s, Sides: %s",
			numberStyle.Render(strconv.Itoa(region.area)),
			numberStyle.Render(strconv.Itoa(region.sides)),
		))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *day12Board) viewSolution() string {
	return fmt.Sprintf("Cost: %s\n", solutionStyle.Render(strconv.Itoa(b.solution)))
}

// sided regions look like this
// +-------+
// |A A A A|
// +-------+
// but can also go up/down and form funny shapes
func (b *day12Board) findSidedRegion(region *day12Region, x, y int) (area, sides int) {
	plotType := region.plotType
	// mark this path as visited
	b.visit(x, y)
	area += 1
	region.area += 1

	if b.onStep != nil {
		b.onStep()
	}

	// find the sides of this square
	currentEdges := b.findSides(plotType, x, y)
	b.foundEdges[y][x] = currentEdges

	if (currentEdges&sideUp) > 0 && b.placeNewEdge(plotType, sideUp, x, y) {
		sides++
		region.sides++
		slog.Debug(fmt.Sprintf("%s %d,%d added side: up", string(plotType), x, y))
	}
	if (currentEdges&sideRight) > 0 && b.placeNewEdge(plotType, sideRight, x, y) {
		sides++
		region.sides++
		slog.Debug(fmt.Sprintf("%s %d,%d added side: right", string(plotType), x, y))
	}
	if (currentEdges&sideDown) > 0 && b.placeNewEdge(plotType, sideDown, x, y) {
		sides++
		region.sides++
		slog.Debug(fmt.Sprintf("%s %d,%d added side: down", string(plotType), x, y))
	}
	if (currentEdges&sideLeft) > 0 && b.placeNewEdge(plotType, sideLeft, x, y) {
		sides++
		region.sides++
		slog.Debug(fmt.Sprintf("%s %d,%d added side: left", string(plotType), x, y))
	}

	if b.onStep != nil {
		b.onStep()
	}
	slog.Debug(fmt.Sprintf("%s area %d, sides: %d", string(plotType), region.area, region.sides))

	// now move to any like squares, sending our sides wih it
	for _, dir := range cardinalDirections {
		offsetX, offsetY := dir.offsetMultiplier()

		if b.visited(x+offsetX, y+offsetY) || !b.samePlotType(plotType, x+offsetX, y+offsetY) {
			continue
		}

		// a neighbor is the same plotType, add its area and perimiter to ours
		a, s := b.findSidedRegion(region, x+offsetX, y+offsetY)
		area += a
		sides += s
	}

	// clear out this board entry
	return area, sides
}

func (b *day12Board) placeNewEdge(plotType rune, side uint, x, y int) bool {
	currentEdges := getBoardValue(x, y, b.foundEdges)
	offsetX, offsetY := 0, 0
	switch side {
	// up/down side checks need to look left and right for already found edges
	case sideUp:
		fallthrough
	case sideDown:
		offsetX = 1
	case sideLeft:
		fallthrough
	case sideRight:
		offsetY = 1
	}

	if (currentEdges & side) == 0 {
		return false
	}

	// if any neighbor in this line has placed an edge, don't place one
	count := 1
	for {
		pos := position{x + (offsetX * count), y + (offsetY * count)}
		// stop checking if we hit a board position that is different from us
		if !validPosition(pos, b.width, b.height) || !b.samePlotType(plotType, pos.x, pos.y) {
			break
		}
		if (getBoardValue(pos.x, pos.y, b.foundEdges) & side) > 0 {
			return false
		}
		// don't continue searching neighbors if we've visited this one
		// we only skip over unvisited spaces
		if b.visited(pos.x, pos.y) {
			break
		}
		if b.findSides(plotType, pos.x, pos.y)&side == 0 {
			break
		}
		count++
	}

	// check neighbors in the other direction
	count = 1
	for {
		pos := position{x - (offsetX * count), y - (offsetY * count)}
		// stop checking if we hit a board position that is different from us
		if !validPosition(pos, b.width, b.height) || !b.samePlotType(plotType, pos.x, pos.y) {
			break
		}
		if (getBoardValue(pos.x, pos.y, b.foundEdges) & side) > 0 {
			return false
		}
		// don't continue searching neighbors if we've visited this one
		// we only skip over unvisited spaces
		if b.visited(pos.x, pos.y) {
			break
		}
		// vc
		// CC
		// Qc ^ checking up, we have a left side and another left side up a ways on the board, but the side is broken by unvisited squares
		// if the line is broken, don't continue
		if b.findSides(plotType, pos.x, pos.y)&side == 0 {
			break
		}
		count++
	}
	return true
}

func (b *day12Board) samePlotType(plotType rune, x, y int) bool {
	return unicode.ToUpper(getBoardValue(x, y, b.board)) == plotType
}

func (b *day12Board) visit(x, y int) {
	b.board[y][x] = unicode.ToLower(b.board[y][x])
}

func (b *day12Board) visited(x, y int) bool {
	return unicode.IsLower(getBoardValue(x, y, b.board))
}

func (b *day12Board) findSides(plotType rune, x, y int) uint {
	var sidesFound uint
	for _, dir := range cardinalDirections {
		offsetX, offsetY := dir.offsetMultiplier()
		test := getBoardValue(x+offsetX, y+offsetY, b.board)

		if test == unicode.ToLower(plotType) {
			continue
		}
		if test != plotType {
			// only add sides if we haven't found it yet
			if dir == directionUp {
				sidesFound |= sideUp
			}
			if dir == directionRight {
				sidesFound |= sideRight
			}
			if dir == directionDown {
				sidesFound |= sideDown
			}
			if dir == directionLeft {
				sidesFound |= sideLeft
			}
		}
	}
	return sidesFound
}
