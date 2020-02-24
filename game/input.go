package game

import "fmt"

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
	Action
	Search
	ZoomIn
	ZoomOut
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
	case Action:
		ns, _ := getNeighbors(level, level.Player.Pos)
		for _, n := range ns {
			checkDoor(n, level)
		}
	}
	if canMove(toPos, level) {
		level.Player.Move(toPos, level)
	} else {
		exists, e := hasEnemy(toPos, level)
		if exists {
			damageAmount := level.Player.Level * 5
			fmt.Println(level.Player.Name, "attacked", e.Name, "for", damageAmount, "damage!")
			e.Health -= int(damageAmount)
			if e.Health <= 0 {
				fmt.Println(e.Name, "is dead.")
				e.IsDead = true
			}
		}
		checkDoor(toPos, level)
	}

}
