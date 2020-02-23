package game

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Enemy struct {
	Aggressive bool
	Character
	path []Position
}

func NewEnemy(name string, level float64, pos Position, texture *sdl.Texture, textureRect *sdl.Rect) *Enemy {
	var newEnemy Enemy
	newEnemy.Name = name
	newEnemy.Level = level
	newEnemy.Pos = pos
	newEnemy.texture = texture
	newEnemy.TextureRect = *textureRect
	return &newEnemy
}

func (enemy *Enemy) Draw(renderer *sdl.Renderer, tileSize, offsetX, offsetY int32) {
	enemyDestRect := sdl.Rect{
		X: int32(enemy.Pos.X)*tileSize + offsetX,
		Y: int32(enemy.Pos.Y)*tileSize + offsetY,
		W: int32(tileSize),
		H: int32(tileSize),
	}
	renderer.Copy(enemy.texture, &enemy.TextureRect, &enemyDestRect)
}

func (enemy *Enemy) distanceToCharacter(character *Character) int {
	dx := float64(enemy.Pos.X - character.Pos.X)
	dy := float64(enemy.Pos.Y - character.Pos.Y)
	distance := math.Sqrt(dx*dx + dy*dy)
	return int(distance)
}

func (enemy *Enemy) Update(level *Level) {
	if enemy.distanceToCharacter(&level.Player.Character) < 5 {
		enemy.Aggressive = true
	}
	if !enemy.Aggressive && enemy.path == nil {
		enemy.path = astar(level, enemy.Pos, getRandomPositionInsideCircle(10, enemy.Pos))
	} else if enemy.Aggressive {
		enemy.path = astar(level, enemy.Pos, level.Player.Pos)
		if enemy.path == nil {
			enemy.Aggressive = false
		}
	}
	if enemy.path != nil {
		for _, p := range enemy.path {
			level.Debug[p] = true
		}
		if enemy.Pos == enemy.path[len(enemy.path)-1] {
			fmt.Println("enemy reached end of path")
			enemy.path = nil
			return
		}
		enemy.Move(enemy.path[0], level)
		if len(enemy.path) > 1 {
			enemy.path = enemy.path[1:]
		}
	}
}
