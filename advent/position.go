package advent

type position struct {
	x int
	y int
}

type positionDirection struct {
	position
	direction direction
}

func (p1 position) addDirection(dir direction) position {
	x, y := dir.offsetMultiplier()
	return position{p1.x + x, p1.y + y}
}
