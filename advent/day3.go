package advent

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type Day3 struct {
}

func (d *Day3) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

var p = message.NewPrinter(language.English)

func (d *Day3) readInput(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	input, err := io.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("error reading file contents %w", err)
	}

	return input, nil

}

func (d *Day3) part1(filename string) error {

	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	result, err := d.evaluateInput(input)
	if err != nil {
		return err
	}

	fmt.Printf("day3a: %s: %d", filename, result)

	return nil
}

func (d *Day3) part2(filename string) error {

	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	window := input[:]
	result := 0
	offsetIndex := 0
	for {
		// find the don't() and match up to it
		donotIndex := strings.Index(string(window), "don't()")
		fmt.Printf("evaluating: %d-%d\n", offsetIndex, offsetIndex+donotIndex)
		subResult, err := d.evaluateInput(window[:donotIndex])
		if err != nil {
			return err
		}
		result += subResult
		p.Printf("total: %v\n", number.Decimal(result))

		// find the next do
		doIndex := strings.Index(string(window[donotIndex:]), "do()")
		if doIndex == -1 {
			fmt.Printf("no more do()s after %d\n", offsetIndex+donotIndex)
			break
		}
		// reset the window to start at the do()
		fmt.Printf("skipping %d-%d\n%s\n", offsetIndex+donotIndex, offsetIndex+donotIndex+doIndex, window[donotIndex:donotIndex+doIndex+4])
		window = window[donotIndex+doIndex:]
		offsetIndex += doIndex
	}
	p.Printf("total: %v\n", number.Decimal(result))

	fmt.Printf("day3b: %s: %d\n", filename, result)

	return nil
}

func (d *Day3) evaluateInput(input []byte) (int, error) {
	pattern := `mul\((\d{1,3}),(\d{1,3})\)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return 0, fmt.Errorf("bad regex %s %w", pattern, err)
	}

	// Find all matches
	matches := re.FindAllSubmatch(input, -1)

	result := 0
	for _, mul := range matches {
		x, err := strconv.Atoi(string(mul[1]))
		if err != nil {
			return 0, fmt.Errorf("match not a number %s %s %w", mul, mul[1], err)
		}

		y, err := strconv.Atoi(string(mul[2]))
		if err != nil {
			return 0, fmt.Errorf("match not a number %s %s %w", mul, mul[2], err)
		}

		result += x * y
		p.Printf("evaluated %s, subtotal: %v\n", mul, number.Decimal(result))
	}

	return result, nil
}
