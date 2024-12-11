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

type Day7 struct {
}

type day7Equation struct {
	num      int
	result   int
	values   []int
	solution []operator
}

type operator int

const (
	operatorAdd operator = iota
	operatorMul
	operatorCat
)

func (d *Day7) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day7) readInput(filename string) ([]day7Equation, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var input []day7Equation

	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		resultSplit := strings.Split(line, ":")
		if len(resultSplit) != 2 {
			return nil, fmt.Errorf("line doesn't contain two numbers: %s", line)
		}

		var err error
		result, err := strconv.Atoi(resultSplit[0])
		if err != nil {
			return nil, fmt.Errorf("result is not a number: %s %v", resultSplit[0], err)
		}

		valueSplit := strings.Split(strings.TrimSpace(resultSplit[1]), " ")
		values := make([]int, len(valueSplit))
		for i, v := range valueSplit {
			values[i], err = strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("value is not a number line: %s value: %s %v", line, v, err)
			}
		}

		input = append(input, day7Equation{num: index, result: result, values: values})
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return input, nil
}

func (d *Day7) part1(filename string) error {

	equations, err := d.readInput(filename)
	if err != nil {
		return err
	}

	for i := range equations {
		eq := &equations[i]
		numValues := len(eq.values)
		combos := generateCombinations([]operator{operatorAdd, operatorMul}, numValues-1)

		for _, operators := range combos {
			result := d.eval(eq.values[0], eq.values[1], operators[0])
			for i := 2; i < numValues; i++ {
				result = d.eval(result, eq.values[i], operators[i-1])
			}

			fmt.Printf("%s\n", eq.view(operators, result))

			if result == eq.result {
				eq.solution = operators
				break
			}
		}
	}

	sum := 0
	count := 0
	for _, eq := range equations {
		if eq.solution != nil {
			sum += eq.result
			count++
		}
	}

	// output the result
	fmt.Printf("Valid Tests: %s, Sum of Test Values: %s\n", correctResultStyle.Render(strconv.Itoa(count)), solutionStyle.Render(strconv.Itoa(sum)))

	// Output the dist
	return nil
}

func (d *Day7) part2(filename string) error {
	equations, err := d.readInput(filename)
	if err != nil {
		return err
	}

	numWorkers := 19

	// create a bubbletea program
	p := tui.NewViewportProgram(tui.NewModel("Day 7 - Part 2").WithViewport(make([]string, numWorkers+1)))

	jobs := make(chan *day7Equation, len(equations))    // Channel to queue jobs
	results := make(chan *day7Equation, len(equations)) // Channel to collect results

	// Worker function
	worker := func(id int, equations <-chan *day7Equation, results chan<- *day7Equation) {
		for eq := range equations {
			time.Sleep(50 * time.Millisecond)
			numValues := len(eq.values)
			combos := generateCombinations([]operator{operatorAdd, operatorMul, operatorCat}, numValues-1)

			var solution []operator
			var solutionResult int
			for _, operators := range combos {
				result := d.eval(eq.values[0], eq.values[1], operators[0])
				for i := 2; i < numValues; i++ {
					result = d.eval(result, eq.values[i], operators[i-1])
				}

				solution = operators
				solutionResult = result

				if result == eq.result {
					eq.solution = operators
					break
				}
			}
			p.Send(tui.UpdateViewportLine(id, fmt.Sprintf("%d: %s", eq.num, eq.view(solution, solutionResult))))
			results <- eq
		}
	}

	sum := 0
	count := 0

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		go worker(i, jobs, results)
	}

	// send all the jobs to the workers
	for i := 0; i < len(equations); i++ {
		jobs <- &equations[i]
	}
	close(jobs)

	// update the ui as jobs come in
	go func() {
		for a := 0; a < len(equations); a++ {
			result := <-results
			if result.solution != nil {
				count++
				sum += result.result
				// p.Send(updateViewportLine{lineNum: numWorkers, line: fmt.Sprintf("Valid Tests: %s, Sum of Test Values: %s", correctResultStyle.Render(strconv.Itoa(count)), solutionStyle.Render(strconv.Itoa(sum)))})
				p.Send(tui.UpdateViewportLine(numWorkers, fmt.Sprintf("Valid Tests: %s, Sum of Test Values: %s", correctResultStyle.Render(strconv.Itoa(count)), solutionStyle.Render("<redacted>"))))
			}
		}
	}()

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("could not start program: %v", err)
	}

	// output the result
	fmt.Printf("Valid Tests: %s, Sum of Test Values: %s\n", correctResultStyle.Render(strconv.Itoa(count)), solutionStyle.Render(strconv.Itoa(sum)))

	// Output the dist
	return nil
}

func (eq *day7Equation) view(solution []operator, result int) string {
	var sb strings.Builder

	sb.WriteString(numberStyle.Render(strconv.Itoa(eq.values[0])))
	for i := 1; i < len(eq.values); i++ {
		if solution == nil {
			sb.WriteString(operatorStyle.Render(" ? "))
		} else {
			switch solution[i-1] {
			case operatorAdd:
				sb.WriteString(operatorStyle.Render(" + "))
			case operatorMul:
				sb.WriteString(operatorStyle.Render(" * "))
			case operatorCat:
				sb.WriteString(operatorStyle.Render(" || "))
			}
		}

		sb.WriteString(numberStyle.Render(strconv.Itoa(eq.values[i])))
	}

	if result == eq.result {
		sb.WriteString(" = " + correctResultStyle.Render(strconv.Itoa(result)))
	} else {
		sb.WriteString(" = " + incorrectResultStyle.Render(strconv.Itoa(result)))
	}

	return sb.String()
}

// generateCombinations generates all combinations of operators for a given size
func generateCombinations(ops []operator, size int) [][]operator {
	if size == 0 {
		return [][]operator{{}}
	}

	smallerCombos := generateCombinations(ops, size-1) // Recursive step
	var result [][]operator

	for _, combo := range smallerCombos {
		for _, op := range ops {
			newCombo := append([]operator{}, combo...) // Copy current combination
			newCombo = append(newCombo, op)            // Add the new operator
			result = append(result, newCombo)
		}
	}

	return result
}

func (d *Day7) eval(num1, num2 int, op operator) int {
	switch op {
	case operatorAdd:
		return num1 + num2
	case operatorMul:
		return num1 * num2
	case operatorCat:
		cat, err := strconv.Atoi(strconv.Itoa(num1) + strconv.Itoa(num2))
		if err != nil {
			panic(fmt.Sprintf("can't cat %d %d", num1, num2))
		}
		return cat
	}

	panic("unknown operator")
}
