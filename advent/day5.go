package advent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Day5 struct {
}

func (d *Day5) Run(part int, filename string, opts ...Option) error {
	switch part {
	case 1:
		return d.part1(filename)
	case 2:
		return d.part2(filename)
	default:
		return fmt.Errorf("part %d not valid", part)
	}
}

type inputDay5 struct {
	orderingRules [][2]int
	pageUpdates   [][]int
}

type day5OrderKey struct {
	before int
	after  int
}

// read day4 input as a series of lines
func (d *Day5) readInput(filename string) (inputDay5, error) {
	file, err := os.Open(filename)
	if err != nil {
		return inputDay5{}, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var input inputDay5
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		if strings.Contains(line, "|") {
			split := strings.Split(line, "|")
			rule := [2]int{}
			if len(split) != 2 {
				return inputDay5{}, fmt.Errorf("line %s doesn't have the right values", line)
			}

			rule[0], err = strconv.Atoi(split[0])
			if err != nil {
				return inputDay5{}, fmt.Errorf("line %s invalid number %w", line, err)
			}
			rule[1], err = strconv.Atoi(split[1])
			if err != nil {
				return inputDay5{}, fmt.Errorf("line %s invalid number %w", line, err)
			}
			input.orderingRules = append(input.orderingRules, rule)
		} else if strings.Contains(line, ",") {
			split := strings.Split(line, ",")
			pages := []int{}
			for _, s := range split {
				page, err := strconv.Atoi(s)
				if err != nil {
					return inputDay5{}, fmt.Errorf("line %s invalid number %w", line, err)
				}
				pages = append(pages, page)
			}
			input.pageUpdates = append(input.pageUpdates, pages)
		}
	}

	if err := scanner.Err(); err != nil {
		return inputDay5{}, fmt.Errorf("error reading file: %w", err)
	}

	return input, nil
}

func (d *Day5) part1(filename string) error {

	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	// make a map of before and after page rules
	rules := make(map[day5OrderKey]bool)
	for _, rule := range input.orderingRules {
		rules[day5OrderKey{rule[0], rule[1]}] = true
	}

	var valids [][]int
	for _, update := range input.pageUpdates {
		valid, badRule := d.eval(update, rules)

		if valid {
			valids = append(valids, update)
			fmt.Printf("%s %v\n", correctResultStyle.Render("valid"), update)
		} else {
			fmt.Printf("%s %v - %v\n", incorrectResultStyle.Render("invalid"), update, badRule)
		}
	}

	totalMids := 0
	numValid := len(valids)
	for _, update := range valids {
		mid := update[len(update)/2]
		totalMids += mid
	}

	fmt.Printf("numValid: %s, total mids: %s\n", numberStyle.Render(strconv.Itoa(numValid)), solutionStyle.Render(strconv.Itoa(totalMids)))

	return nil
}

func (d *Day5) part2(filename string) error {

	input, err := d.readInput(filename)
	if err != nil {
		return err
	}

	// make a map of before and after page rules
	rules := make(map[day5OrderKey]bool)
	for _, rule := range input.orderingRules {
		rules[day5OrderKey{rule[0], rule[1]}] = true
	}

	var invalids [][]int
	for _, update := range input.pageUpdates {
		valid, _ := d.eval(update, rules)

		if !valid {
			invalids = append(invalids, update)
		}
	}

	fmt.Printf("Fixing bad pages...\n\n")

	// sort the invalids
	numInvalid := len(invalids)
	totalMids := 0
	for _, update := range invalids {
		fmt.Printf("%s %v => ", incorrectResultStyle.Render("invalid"), update)
		for {
			valid, _ := d.eval(update, rules)
			if valid {
				break
			}
			// not good yet, repair it
			d.repair(update, rules)
		}

		fmt.Printf("%v %s\n", update, correctResultStyle.Render("valid"))

		mid := update[len(update)/2]
		totalMids += mid
	}

	fmt.Printf("Num Invalid: %s, total mids: %s\n", numberStyle.Render(strconv.Itoa(numInvalid)), solutionStyle.Render(strconv.Itoa(totalMids)))

	return nil
}

func (d *Day5) eval(update []int, rules map[day5OrderKey]bool) (bool, *day5OrderKey) {
	for i, page := range update {
		if i < len(update)-1 {
			// for 75,47,61,53,29
			// 47 must not come before 75
			rule := day5OrderKey{update[i+1], page}
			if rules[rule] {
				return false, &rule
			}
		}
		if i > 0 {
			// for 75,47,61,53,29 if on page 47
			// 47 must no come before 75
			rule := day5OrderKey{page, update[i-1]}
			if rules[rule] {
				return false, &rule
			}
		}
	}
	return true, nil
}

func (d *Day5) repair(update []int, rules map[day5OrderKey]bool) {
	for i, page := range update {
		if i < len(update)-1 {
			// for
			// 97|75
			// 75,97,47,61,53
			// 97 must come before 75, so swap them
			rule := day5OrderKey{update[i+1], page}
			if rules[rule] {
				update[i], update[i+1] = update[i+1], update[i]
			}
		}
		if i > 0 {
			// for
			// 29|13
			// 61,13,29
			// 29 comes before 13, so swap them
			rule := day5OrderKey{page, update[i-1]}
			if rules[rule] {
				update[i], update[i-1] = update[i-1], update[i]
			}
		}
	}
}
