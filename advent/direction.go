package advent

import (
	"fmt"
	"slices"
)

// clockwise directions
type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionUpRight
	DirectionDownRight
	DirectionDownLeft
	DirectionUpLeft
)

var CardinalDirections = []Direction{
	DirectionUp,
	DirectionRight,
	DirectionDown,
	DirectionLeft,
}

var AdjacentDirections = []Direction{
	DirectionUp,
	DirectionUpRight,
	DirectionRight,
	DirectionDownRight,
	DirectionDown,
	DirectionDownLeft,
	DirectionLeft,
	DirectionUpLeft,
}

const (
	SideUp    = 0x01
	SideRight = 0x02
	SideDown  = 0x04
	SideLeft  = 0x08
)

func removeSide(sides uint, dir Direction) uint {
	switch dir {
	case DirectionUp:
		return (sides ^ SideUp)
	case DirectionRight:
		return (sides ^ SideRight)
	case DirectionDown:
		return (sides ^ SideDown)
	case DirectionLeft:
		return (sides ^ SideLeft)
	}
	panic(fmt.Sprintf("can't remove side for direction %v", dir))
}

// turn 90 degrees right
func (d Direction) turnRight() Direction {
	switch d {
	case DirectionUp:
		return DirectionRight
	case DirectionRight:
		return DirectionDown
	case DirectionDown:
		return DirectionLeft
	case DirectionLeft:
		return DirectionUp
	}

	return d
}

func directionFromChar(c rune) Direction {
	switch c {
	case '^':
		return DirectionUp
	case '>':
		return DirectionRight
	case 'v':
		return DirectionDown
	case '<':
		return DirectionLeft
	case '↗':
		return DirectionUpRight
	case '↘':
		return DirectionDownRight
	case '↙':
		return DirectionDownLeft
	case '↖':
		return DirectionUpLeft
	}

	return DirectionUp
}

var _ = directionFromChar // might need this, don't want to retype it, don't like the warning

func (d Direction) OffsetMultiplier() (x, y int) {
	switch d {
	case DirectionUp:
		return 0, -1
	case DirectionRight:
		return 1, 0
	case DirectionDown:
		return 0, 1
	case DirectionLeft:
		return -1, 0
	case DirectionUpRight:
		return 1, -1
	case DirectionDownRight:
		return 1, 1
	case DirectionDownLeft:
		return -1, 1
	case DirectionUpLeft:
		return -1, -1
	}
	return 0, 0
}

func DirectionFromOffset(offset Position) Direction {
	switch offset {
	case Position{0, -1}:
		return DirectionUp
	case Position{1, 0}:
		return DirectionRight
	case Position{0, 1}:
		return DirectionDown
	case Position{-1, 0}:
		return DirectionLeft
	case Position{1, -1}:
		return DirectionUpRight
	case Position{1, 1}:
		return DirectionDownRight
	case Position{-1, 1}:
		return DirectionDownLeft
	case Position{-1, -1}:
		return DirectionUpLeft
	}
	panic(fmt.Sprintf("invalid offset %v", offset))
}

func (d Direction) getChar() rune {
	switch d {
	case DirectionUp:
		return '^'
	case DirectionRight:
		return '>'
	case DirectionDown:
		return 'v'
	case DirectionLeft:
		return '<'
	case DirectionUpRight:
		return '↗'
	case DirectionDownRight:
		return '↘'
	case DirectionDownLeft:
		return '↙'
	case DirectionUpLeft:
		return '↖'
	}

	return ' '
}

func (d Direction) MinTurns(other Direction) int {
	dIndex := slices.Index(CardinalDirections, d)
	otherIndex := slices.Index(CardinalDirections, other)
	if otherIndex > dIndex {
		return min(2, otherIndex-dIndex)
	}
	return min(2, dIndex-otherIndex)
}

func (d Direction) String() string {
	c := d.getChar()
	if c == ' ' {
		return ""
	}
	return string(c)
}
