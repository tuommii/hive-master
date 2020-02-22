package game

import "github.com/veandco/go-sdl2/sdl"

type Entity struct {
	Pos         Position
	TextureRect sdl.Rect
}

type Movable interface {
	Move(pos Position)
}

func (e *Entity) Move(pos Position) {
	e.Pos = pos
}
