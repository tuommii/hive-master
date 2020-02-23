package game

import "github.com/veandco/go-sdl2/sdl"

type Drawable interface {
	Draw(renderer *sdl.Renderer)
}

type Movable interface {
	Move(pos Position)
}

type Character struct {
	Pos         Position
	Name        string
	Level       float64
	Health      int
	IsDead      bool
	texture     *sdl.Texture
	TextureRect sdl.Rect
}

func (e *Character) Move(pos Position, level *Level) {
	if canMove(pos, level) {
		e.Pos = pos
	}
}
