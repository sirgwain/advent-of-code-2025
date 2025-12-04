package advent

type Position struct {
	X int
	Y int
}

type PositionDirection struct {
	Position
	Direction Direction
}

func (p1 Position) addDirection(dir Direction) Position {
	x, y := dir.OffsetMultiplier()
	return Position{p1.X + x, p1.Y + y}
}
