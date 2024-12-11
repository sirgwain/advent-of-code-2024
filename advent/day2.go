package advent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Day2 struct {
}

func (d *Day2) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

func (d *Day2) readInput(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var reports [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var levels []int

		parts := strings.Fields(line)
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("failed to read number from line: %s %w", line, err)
			}
			levels = append(levels, num)
		}
		reports = append(reports, levels)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return reports, nil
}

func (d *Day2) part1(filename string) error {

	reports, err := d.readInput(filename)
	if err != nil {
		return err
	}

	safeReports := 0
	for _, report := range reports {
		// make sure they are all increasing/decreasing at the same level
		safe := d.checkLevels(report)

		if safe {
			safeReports++
			fmt.Printf("%v: Safe\n", report)
		}
	}

	fmt.Printf("day2: %s: %d reports are safe", filename, safeReports)

	return nil
}

func (d *Day2) part2(filename string) error {

	reports, err := d.readInput(filename)
	if err != nil {
		return err
	}

	safeReports := 0
	for i, report := range reports {
		fmt.Printf("report %d: ", i)
		// make sure they are all increasing/decreasing at the same level
		safe := d.checkLevels(report)

		if !safe {
			// try by removing a level each time
			for li := 0; li < len(report); li++ {
				fmt.Printf("report %d (without item %d): ", i, li)

				safe = d.checkLevels(removeIndex(report, li))
				if safe {
					break
				}
			}
		}

		if safe {
			safeReports++
			fmt.Printf("%v: Safe\n", report)
		}
	}

	fmt.Printf("day2: %s: %d reports are safe", filename, safeReports)

	return nil
}

func (d *Day2) checkLevels(report []int) bool {
	var lastDiff *int
	for j := 1; j < len(report); j++ {
		l0 := report[j-1]
		l1 := report[j]
		diff := int(l0 - l1)
		if diff > 3 || diff < -3 {
			fmt.Printf("%v: Unsafe because level %d -> %d too large %d\n", report, l0, l1, diff)
			return false
		}
		if diff == 0 {
			fmt.Printf("%v: Unsafe because level %d -> %d not decreasing or increasing\n", report, l0, l1)
			return false
		}

		if lastDiff == nil {
			lastDiff = &diff
		} else {
			if (*lastDiff > 0 && diff < 0) || (*lastDiff < 0 && diff > 0) {
				fmt.Printf("%v: Unsafe because level %d -> %d is %d, last was %d\n", report, l0, l1, diff, *lastDiff)
				return false
			}
		}
	}
	return true
}
