package advent

import "fmt"

// clockwise directions
type direction int

const (
	directionUp direction = iota
	directionRight
	directionDown
	directionLeft
	directionUpRight
	directionDownRight
	directionDownLeft
	directionUpLeft
)

var cardinalDirections = []direction{
	directionUp,
	directionRight,
	directionDown,
	directionLeft,
}

const (
	sideUp    = 0x01
	sideRight = 0x02
	sideDown  = 0x04
	sideLeft  = 0x08
)

func removeSide(sides uint, dir direction) uint {
	switch dir {
	case directionUp:
		return (sides ^ sideUp)
	case directionRight:
		return (sides ^ sideRight)
	case directionDown:
		return (sides ^ sideDown)
	case directionLeft:
		return (sides ^ sideLeft)
	}
	panic(fmt.Sprintf("can't remove side for direction %v", dir))
}

// turn 90 degrees right
func (d direction) turnRight() direction {
	switch d {
	case directionUp:
		return directionRight
	case directionRight:
		return directionDown
	case directionDown:
		return directionLeft
	case directionLeft:
		return directionUp
	}

	return d
}

func directionFromChar(c rune) direction {
	switch c {
	case '^':
		return directionUp
	case '>':
		return directionRight
	case 'v':
		return directionDown
	case '<':
		return directionLeft
	case '↗':
		return directionUpRight
	case '↘':
		return directionDownRight
	case '↙':
		return directionDownLeft
	case '↖':
		return directionUpLeft
	}

	return directionUp
}

var _ = directionFromChar // might need this, don't want to retype it, don't like the warning

func (d direction) offsetMultiplier() (x, y int) {
	switch d {
	case directionUp:
		return 0, -1
	case directionRight:
		return 1, 0
	case directionDown:
		return 0, 1
	case directionLeft:
		return -1, 0
	case directionUpRight:
		return 1, -1
	case directionDownRight:
		return 1, 1
	case directionDownLeft:
		return -1, 1
	case directionUpLeft:
		return -1, -1
	}
	return 0, 0
}

func (d direction) getChar() rune {
	switch d {
	case directionUp:
		return '^'
	case directionRight:
		return '>'
	case directionDown:
		return 'v'
	case directionLeft:
		return '<'
	case directionUpRight:
		return '↗'
	case directionDownRight:
		return '↘'
	case directionDownLeft:
		return '↙'
	case directionUpLeft:
		return '↖'
	}

	return ' '
}
