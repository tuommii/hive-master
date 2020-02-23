package game

type InputType int

type Input struct {
	Type InputType
}

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Search
	Quit
)

func handleInput(level *Level, input *Input) {
	toPos := level.Player.Pos
	switch input.Type {
	case Up:
		toPos.Y--
	case Down:
		toPos.Y++
	case Left:
		toPos.X--
	case Right:
		toPos.X++
	}
	if canMove(toPos, level) {
		level.Player.Move(toPos)
	} else {
		checkDoor(toPos, level)
	}
}
