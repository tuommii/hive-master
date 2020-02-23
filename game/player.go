package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	Name  string
	Level float64
	Entity
}

func NewPlayer(name string, level float64, pos Position, textureRect *sdl.Rect) *Player {
	var player Player
	player.Name = name
	player.Level = level
	player.TextureRect = *textureRect
	player.Pos = pos
	return &player
}
