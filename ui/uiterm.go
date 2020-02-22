package ui

import (
	"fmt"

	"github.com/wehard/hive-master/game"
)

type UITerm struct{}

func (ui UITerm) Draw(level *game.Level) {
	for y, row := range level.Map {
		for x, _ := range row {
			fmt.Printf("%c", level.Map[y][x])
		}
		fmt.Print("\n")
	}
}
