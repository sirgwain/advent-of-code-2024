package advent

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Day9 struct {
	*Options
}

func (d *Day9) Run(part int, filename string, opts ...Option) error {
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

func (d *Day9) readInput(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	input := make([]int, len(content))
	for i, c := range content {
		num, err := strconv.Atoi(string(c))
		if err != nil {
			return nil, fmt.Errorf("character is not an int: %v %w", c, err)
		}
		input[i] = num
	}

	return input, nil
}

// part1 uses two indices to move forward and in reverse through the input
func (d *Day9) part1(input []int) error {

	// initialize the end of array index
	// and the endID (for 20 numbers, we have an end ID of 9 because the last block has no empty space)
	// we will be adding these end blocks to empty slots, so record how many endBlocks we have to move
	j := len(input) - 1
	endID := (len(input) - 1) / 2
	endBlocksToMove := input[j]

	// disk := make([]int, 0, len(input))
	// compute the checksum as we go
	checksum := 0
	diskIndex := 0
	for i, id := 0, 0; i < len(input); i, id = i+2, id+1 {
		if i >= j {
			// add the remaining end blocks to the end
			// compute the final checksum, and we're done
			for range endBlocksToMove {
				// disk = append(disk, endID)
				checksum = checksum + diskIndex*endID
				// fmt.Printf("%d", endID)
			}

			break
		}
		// each part is two numbers, the number of blocks and the number of empty spaces
		numBlocks := input[i]
		numEmpty := input[i+1]

		// compute the checksum for each of these blocks
		for range numBlocks {
			// disk = append(disk, id)

			// update the checksum and the current index on the new disk
			checksum = checksum + diskIndex*id
			diskIndex++
			// fmt.Printf("%d", id)
		}

		// for each empty block, grab a block from the end and fill it in
		for range numEmpty {
			if endBlocksToMove == 0 {
				// this end block is fully moved, move to the next end block
				j -= 2 // skip the empty

				// reset endBlocksToMove to our new end block
				endBlocksToMove = input[j]
				// decrement the endID
				endID--
			}
			// disk = append(disk, endID)

			// move an endID block into the empty spot
			// update the checksum
			checksum = checksum + diskIndex*endID
			diskIndex++

			// we used up one block
			endBlocksToMove--
			// fmt.Printf("%d", endID)
		}
	}

	// fmt.Println()
	// for _, block := range disk {
	// 	fmt.Printf("%d", block)
	// }
	fmt.Printf("\nChecksum: %s\n", solutionStyle.Render(strconv.Itoa(checksum)))

	return nil
}

func (d *Day9) part2(input []int) error {

	size := 0
	for _, n := range input {
		size += n
	}

	diskImage := make([]int, size)
	diskIndex := 0
	for i, id := 0, 1; i < len(input); i, id = i+2, id+1 {
		// each part is two numbers, the number of blocks and the number of empty spaces
		numBlocks := input[i]

		var numEmpty int
		if i < len(input)-2 {
			numEmpty = input[i+1]
		}

		for range numBlocks {
			diskImage[diskIndex] = id
			diskIndex++
		}
		for range numEmpty {
			diskImage[diskIndex] = 0
			diskIndex++
		}
	}
	// output disk image
	// don't do this on the big test, it's slow and kind of useless
	// fmt.Printf("S: %s\n", d.view(diskImage))

	// go backwards and see if each end block can move
	id := diskImage[len(diskImage)-1]
	numToMove := 0
	for j := len(diskImage) - 1; j > 0; j-- {
		if diskImage[j] == id {
			numToMove++
			continue
		}

		if id == 0 {
			// skip empties
			id = diskImage[j]
			numToMove = 1
			continue
		}

		// try and place this id somewhere in the beginning
		emptyIndex := d.findEmptyIndex(diskImage[:j+1], numToMove)
		if emptyIndex != -1 {
			// move here
			for i := 0; i < numToMove; i++ {
				diskImage[emptyIndex+i] = id
				// clear out where it was
				diskImage[j+i+1] = 0
			}
		}
		// don't print on big test, it's SLOW
		// fmt.Printf("%d: %s\n", id-1, d.view(diskImage))

		// start the next one
		id = diskImage[j]
		numToMove = 0
		if id != 0 {
			numToMove = 1
		}

	}

	fmt.Printf("F: %s\n", d.view(diskImage))
	fmt.Printf("\nChecksum: %s\n", solutionStyle.Render(strconv.Itoa(d.checksum(diskImage))))

	return nil
}

func (d *Day9) findEmptyIndex(diskImage []int, size int) int {

	foundEmpty := false
	emptySize := 0
	emptyIndex := -1
	for i := 0; i < len(diskImage); i++ {
		if diskImage[i] == 0 && !foundEmpty {
			emptyIndex = i
			emptySize = 0
			foundEmpty = true
		}

		// record this empty space
		if diskImage[i] == 0 {
			emptySize++
		}

		if emptySize >= size {
			return emptyIndex
		}

		// ran out of empties without finding the one we want, skip it
		if diskImage[i] != 0 {
			emptyIndex = -1
			emptySize = 0
			foundEmpty = false
		}
	}

	return -1
}

func (d *Day9) checksum(diskImage []int) int {
	checksum := 0
	for i, n := range diskImage {
		if n == 0 {
			continue
		}
		// ids are stored as id+1 to make empties easier
		checksum += i * (n - 1)
	}
	return checksum
}

func (d *Day9) view(diskImage []int) string {
	var sb strings.Builder
	for _, n := range diskImage {
		if n == 0 {
			sb.WriteRune('.')
			continue
		}
		sb.WriteString(numberStyle.Render(strconv.Itoa(n - 1)))
	}

	return sb.String()
}
