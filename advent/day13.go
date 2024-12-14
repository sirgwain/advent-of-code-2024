package advent

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Day13 struct {
	*Options
}

type day13Input []day13Machine

type button uint

const (
	buttonA button = iota
	buttonB
)

type day13Machine struct {
	buttonA position
	buttonB position
	prize   position

	tried map[day13Solution]bool
}

type day13Solution struct {
	a int
	b int
}

func (s day13Solution) tokens() int {
	return s.a*3 + s.b
}

func (s day13Solution) empty() bool {
	return s.a == 0 && s.b == 0
}

func (m day13Machine) press(button button, pos position) position {

	if button == buttonA {
		return position{pos.x + m.buttonA.x, pos.y + m.buttonA.y}
	}
	return position{pos.x + m.buttonB.x, pos.y + m.buttonB.y}
}

// Run is the main entry point for a day. It reads the input file and runs the part
func (d *Day13) Run(part int, filename string, opts ...Option) error {
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

func (d *Day13) readInput(filename string) (day13Input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	buttonPattern := `Button [AB]: X\+(\d+), Y\+(\d+)`
	reButton, err := regexp.Compile(buttonPattern)
	if err != nil {
		return nil, fmt.Errorf("bad regex %s %w", buttonPattern, err)
	}

	prizePattern := `Prize: X=(\d+), Y=(\d+)`
	rePrize, err := regexp.Compile(prizePattern)
	if err != nil {
		return nil, fmt.Errorf("bad regex %s %w", prizePattern, err)
	}

	var input []day13Machine

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var machine day13Machine

		// Button A: X+94, Y+34
		line := scanner.Bytes()
		match := reButton.FindSubmatch(line)
		machine.buttonA = position{mustAtoi(match[1]), mustAtoi(match[2])}

		// Button B: X+22, Y+67
		scanner.Scan()
		line = scanner.Bytes()
		match = reButton.FindSubmatch(line)
		machine.buttonB = position{mustAtoi(match[1]), mustAtoi(match[2])}

		scanner.Scan()
		line = scanner.Bytes()
		match = rePrize.FindSubmatch(line)
		machine.prize = position{mustAtoi(match[1]), mustAtoi(match[2])}
		scanner.Scan()

		input = append(input, machine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return input, nil
}

func (d *Day13) part1(input day13Input) error {

	solution := 0

	for _, machine := range input {
		// presses := machine.findBestSolution()
		presses := machine.findBestSolutionWithAlgrebra(0, 100)
		if presses.empty() {
			// no solution
			fmt.Printf("Prize: X=%d, Y=%d\n", machine.prize.x, machine.prize.y)
			fmt.Printf("No solution\n")
			continue
		}
		fmt.Printf("Button A: X+%d, Y+%d\n", machine.buttonA.x, machine.buttonA.y)
		fmt.Printf("Button B: X+%d, Y+%d\n", machine.buttonB.x, machine.buttonB.y)
		fmt.Printf("Prize: X=%d, Y=%d\n", machine.prize.x, machine.prize.y)
		fmt.Printf("\nBest Presses: A: %s, B: %s => %s tokens\n\n",
			solutionStyle.Render(strconv.Itoa(presses.a)),
			solutionStyle.Render(strconv.Itoa(presses.b)),
			solutionStyle.Render(strconv.Itoa(presses.tokens())),
		)
		solution += presses.tokens()
	}

	fmt.Printf("\nSolution: %s\n", solutionStyle.Render(strconv.Itoa(solution)))
	return nil
}

// Button A: X+24, Y+90
// Button B: X+85, Y+62
// Prize: X=6844, Y=6152
//
// oh yeah, I guess this is just two equations two unknowns...
//
// 24a + 85b = 6844 = ax*a + bx*b = px
// 90a + 62b = 6152 = ay*a + by*b = py
// or in matrix land AX=B
func (m *day13Machine) findBestSolutionWithAlgrebra(prizeOffset int, buttonLimit int) day13Solution {

	// Coefficients of the equations
	ax, bx, px := m.buttonA.x, m.buttonB.x, m.prize.x+prizeOffset
	ay, by, py := m.buttonA.y, m.buttonB.y, m.prize.y+prizeOffset

	// Determinant of A
	det := ax*by - ay*bx
	if det == 0 {
		// no solution
		return day13Solution{}
	}

	// Inverse of A multiplied by B
	aPresses := (px*by - py*bx) / det
	bPresses := (ax*py - ay*px) / det

	// check the answer
	arrivedX := m.buttonA.x*aPresses + m.buttonB.x*bPresses
	arrivedY := m.buttonA.y*aPresses + m.buttonB.y*bPresses

	if arrivedX != m.prize.x+prizeOffset || arrivedY != m.prize.y+prizeOffset {
		return day13Solution{}
	}

	fmt.Printf("Pressing button A %d times moves to (%d,%d)\n", aPresses, m.buttonA.x*aPresses, m.buttonA.y*aPresses)
	fmt.Printf("Pressing button B %d times moves to (%d,%d)\n", bPresses, m.buttonB.x*bPresses, m.buttonB.y*bPresses)
	fmt.Printf("Arriving at (%d,%d) prize is at (%d,%d)\n", arrivedX, arrivedY, m.prize.x+prizeOffset, m.prize.y+prizeOffset)
	return day13Solution{a: aPresses, b: bPresses}
}

func (m *day13Machine) findBestSolution() day13Solution {
	m.tried = make(map[day13Solution]bool, 100*100)
	return m.findSolutionFrom(position{}, day13Solution{}, day13Solution{})
}

func (m *day13Machine) findSolutionFrom(pos position, presses, best day13Solution) day13Solution {

	// see if this solution is a winner
	if m.tried[presses] {
		return best
	}

	m.tried[presses] = true
	if pos == m.prize {
		fmt.Printf("Solution: A: %d, B: %d => %d,%d\n", presses.a, presses.b, pos.x, pos.y)
		// found solution
		if presses.tokens() < best.tokens() || best.empty() {
			best = presses
		}
		return best
	}

	// if we've exceeded the total number of button presses
	// if presses.a > 100 || presses.b > 100 {
	// 	return best
	// }

	// try button a
	bestA := m.findSolutionFrom(m.press(buttonA, pos), day13Solution{a: presses.a + 1, b: presses.b}, best)
	bestB := m.findSolutionFrom(m.press(buttonB, pos), day13Solution{a: presses.a, b: presses.b + 1}, best)

	if bestA.empty() && !bestB.empty() {
		return bestB
	} else if bestB.empty() && !bestA.empty() {
		return bestA
	}
	if bestA.tokens() < bestB.tokens() {
		return bestA
	}

	return bestB
}

func (d *Day13) part2(input []day13Machine) error {
	solution := 0

	prizeOffset := 10000000000000
	for i, machine := range input {
		fmt.Printf("Machine %d\n", i+1)
		presses := machine.findBestSolutionWithAlgrebra(prizeOffset, 0)
		if presses.empty() {
			// no solution
			fmt.Printf("Prize: X=%d, Y=%d\n", machine.prize.x+prizeOffset, machine.prize.y+prizeOffset)
			fmt.Printf("No solution\n")
			fmt.Println()
			continue
		}

		fmt.Printf("Button A: X+%d, Y+%d\n", machine.buttonA.x, machine.buttonA.y)
		fmt.Printf("Button B: X+%d, Y+%d\n", machine.buttonB.x, machine.buttonB.y)
		fmt.Printf("Prize: X=%d, Y=%d\n", machine.prize.x+prizeOffset, machine.prize.y+prizeOffset)
		fmt.Printf("\nBest Presses: A: %s, B: %s => %s tokens\n\n",
			solutionStyle.Render(strconv.Itoa(presses.a)),
			solutionStyle.Render(strconv.Itoa(presses.b)),
			solutionStyle.Render(strconv.Itoa(presses.tokens())),
		)
		solution += presses.tokens()
		fmt.Println()
	}

	fmt.Printf("\nSolution: %s\n", solutionStyle.Render(strconv.Itoa(solution)))
	return nil
}
