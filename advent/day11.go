package advent

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Day11 struct {
	*Options

	// keep track of how many stones a stone blinked a
	// certain number of times creates
	stoneCache map[day11CacheEntry]int
}

type day11CacheEntry struct {
	stone int
	times int
}

// Run is the main entry point for a day. It reads the input file and runs the part
func (d *Day11) Run(part int, filename string, opts ...Option) error {
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

func (d *Day11) readInput(filename string) ([]int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	numStrs := strings.Split(string(content), " ")
	input := make([]int, len(numStrs))

	for i, str := range numStrs {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("failed to read number: %s %w", str, err)
		}
		input[i] = num
	}

	return input, nil
}

func (d *Day11) part1(input []int) error {
	// numStones := d.blinkStones(input, 25)
	numStones := d.blinkStonesShantz(input, 25)

	fmt.Printf("\nTotal Stones: %s\n", solutionStyle.Render(strconv.Itoa(numStones)))
	return nil
}

func (d *Day11) part2(input []int) error {
	numStones := d.blinkStones(input, 75)

	fmt.Printf("\nTotal Stones: %s\n", solutionStyle.Render(strconv.Itoa(numStones)))
	return nil
}

// blinkStones blinks each stone n times and sums up their counts
func (d *Day11) blinkStones(stones []int, times int) int {
	d.stoneCache = make(map[day11CacheEntry]int)
	numStones := 0

	for _, stone := range stones {
		numStones += d.blink(stone, times)
	}
	return numStones
}

func (d *Day11) blinkStonesShantz(input []int, times int) int {
	// keep a map of stones and their counts
	// blink all stones of the same id as a group
	// from shantz1. I wanted to try this idea to see how it compared to my original
	stones := make(map[int]int)
	for _, stone := range input {
		stones[stone]++
	}

	for range times {
		// make a new map for this blink round
		newStones := make(map[int]int, len(stones))
		// iterate over every stone in the previous group and run the rules
		for stone, count := range stones {

			stone1, stone2 := d.runRules(stone)
			// the first stone changed number, so put it's count in the new stones map
			newStones[stone1] += count

			if stone2 != -1 {
				// stone2 is a dupe of an old stone batch with a new number
				newStones[stone2] += count
			}
		}
		// reset our stones map for the next loop
		stones = newStones
	}

	// count all the stones after all the blinks
	numStones := 0
	for _, count := range stones {
		numStones += count
	}
	return numStones
}

// blink blinks this stone a number a times and returns the total number of stones at the end
func (d *Day11) blink(stone int, times int) int {
	if times == 0 {
		return 1
	}

	// if we already counted this stone blink n times, return the cache
	if count, ok := d.stoneCache[day11CacheEntry{stone: stone, times: times}]; ok {
		return count
	}

	// run the rules for this new stone
	stone1, stone2 := d.runRules(stone)

	count := 0
	if stone2 != -1 {
		// if we made a second stone, recursively blink it
		// and cache the results
		numStones := d.blink(stone2, times-1)
		count += numStones
		d.stoneCache[day11CacheEntry{stone: stone2, times: times - 1}] = numStones
	}

	// recursively blink this stone one less time
	// and cache the results
	numStones := d.blink(stone1, times-1)
	d.stoneCache[day11CacheEntry{stone: stone1, times: times - 1}] = numStones
	count += numStones
	return count
}

func (d *Day11) runRules(stone int) (stone1, stone2 int) {
	if new, ok := d.rule1(stone); ok {
		return new, -1
	} else if new1, new2, ok := d.rule2(stone); ok {
		return new1, new2
	} else {
		new1, _ := d.rule3(stone)
		return new1, -1
	}
}

// rule1 - 0 becomes 1
func (d *Day11) rule1(stone int) (int, bool) {
	if stone == 0 {
		return 1, true
	}
	return 0, false
}

// rule2 - even numbers are split into two stones
func (d *Day11) rule2(stone int) (int, int, bool) {
	digits := int(math.Log10(float64(stone))) + 1
	if digits%2 == 0 {
		// 1234 split is 1234 / 100 = 12
		//               1234 % 100 = 34
		k := int(math.Pow10(digits / 2))
		stone1 := stone / k
		stone2 := stone % k

		return stone1, stone2, true
	}
	return 0, 0, false
}

// rule2 - even numbers are split into two stones
func (d *Day11) rule3(stone int) (int, bool) {
	return stone * 2024, true
}
