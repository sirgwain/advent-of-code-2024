package advent

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Day1 struct {
}

func (d *Day1) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day1) readInput(filename string) ([]int, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var slice1, slice2 []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("line doesn't contain two numbers: %s", line)
		}

		num1, err1 := strconv.Atoi(parts[0])
		num2, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return nil, nil, fmt.Errorf("error parsing numbers on line: %s", line)
		}

		slice1 = append(slice1, num1)
		slice2 = append(slice2, num2)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading file: %w", err)
	}

	return slice1, slice2, nil
}

func (d *Day1) part1(filename string) error {

	slice1, slice2, err := d.readInput(filename)
	if err != nil {
		return err
	}

	slices.Sort(slice1)
	slices.Sort(slice2)

	dist := 0
	for i := range slice1 {
		dist += int(math.Abs(float64(slice1[i] - slice2[i])))
	}

	// Output the dist
	fmt.Printf("day1: %s dist = %d\n", filename, dist)
	return nil
}

func (d *Day1) part2(filename string) error {
	slice1, slice2, err := d.readInput(filename)
	if err != nil {
		return err
	}

	slice2Occurances := make(map[int]int, len(slice2))
	for _, val := range slice2 {
		slice2Occurances[val]++
	}

	similarity := 0
	for _, val := range slice1 {
		similarity += val * slice2Occurances[val]
	}
	// Output the dist
	fmt.Printf("day2: %s similarity = %d\n", filename, similarity)

	return nil
}
