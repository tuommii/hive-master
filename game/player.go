package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Player struct {
	Character
}

func NewPlayer(name string, level float64, pos Position, texture *sdl.Texture, textureRect *sdl.Rect) *Player {
	var player Player
	player.Name = name
	player.Level = level
	player.texture = texture
	player.TextureRect = *textureRect
	player.Pos = pos
	return &player
}

func (player *Player) Draw(renderer *sdl.Renderer, tileSize, offsetX, offsetY int32) {
	playerDestRect := sdl.Rect{
		X: int32(player.Pos.X)*tileSize + offsetX,
		Y: int32(player.Pos.Y)*tileSize + offsetY,
		W: int32(tileSize),
		H: int32(tileSize),
	}
	renderer.Copy(player.texture, &player.TextureRect, &playerDestRect)
	//label.Draw(Position{level.Player.Pos.X*tileSize + int(offsetX) - tileSize/2, level.Player.Pos.Y*tileSize + int(offsetY) - tileSize/2})
}
