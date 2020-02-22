package game

import (
	"bufio"
	"fmt"
	"os"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
)

type Input struct {
	Type InputType
}

type Tile rune

const (
	Blank Tile = ' '
	Wall  Tile = '#'
	Floor Tile = '.'
)

type Level struct {
	Map    [][]Tile
	Player Player
}

func LoadLevelFromFile(filename string) *Level {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	levelLines := make([]string, 0)
	cols := 0
	rows := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > cols {
			cols = len(line)
		}
		levelLines = append(levelLines, line)
		rows++
	}
	level := &Level{}
	level.Map = make([][]Tile, rows)
	for i := range level.Map {
		level.Map[i] = make([]Tile, cols)
	}
	for y := 0; y < rows; y++ {
		line := levelLines[y]
		for x := 0; x < len(line); x++ {
			c := line[x]
			var t Tile
			switch c {
			case ' ', '\t', '\r':
				t = Blank
			case '#':
				t = Wall
			case '.':
				t = Floor
			default:
				t = Blank
			}
			level.Map[y][x] = t
		}
	}
	return level
}

func handleInput(level *Level, input *Input) {
	switch input.Type {
	case Up:
		level.Player.Y--
	case Down:
		level.Player.Y++
	case Left:
		level.Player.X--
	case Right:
		level.Player.X++
	}
	fmt.Println(input.Type)
}

func Run(ui GameUI) {
	level := LoadLevelFromFile("game/maps/level1.map")
	level.Player.X = 5
	level.Player.Y = 5
	for {
		ui.Draw(level)
		input := ui.GetInput()
		if input.Type == Quit {
			return
		}
		handleInput(level, input)
	}
}
