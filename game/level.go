package game

import (
	"bufio"
	"os"
)

type Level struct {
	Map    [][]Tile
	Player *Player
	Debug  map[Position]bool
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
				t.Rune = Blank
			case '#':
				t.Rune = Wall
			case '.':
				t.Rune = Floor
			case '|':
				t.Rune = ClosedDoor
			case '/':
				t.Rune = OpenDoor
			case '*':
				t.Rune = ClosedChest
			case 'o':
				t.Rune = OpenChest
			default:
				t.Rune = Blank
			}
			level.Map[y][x] = t
		}
	}
	return level
}

func canMove(pos Position, level *Level) bool {
	if level.Map[pos.Y][pos.X].Rune == Wall || level.Map[pos.Y][pos.X].Rune == ClosedDoor || level.Map[pos.Y][pos.X].Rune == ClosedChest {
		return false
	}
	return true
}

func checkDoor(pos Position, level *Level) {
	if level.Map[pos.Y][pos.X].Rune == ClosedDoor {
		level.Map[pos.Y][pos.X].Rune = OpenDoor
	} else if level.Map[pos.Y][pos.X].Rune == ClosedChest {
		level.Map[pos.Y][pos.X].Rune = OpenChest
	}
}

func getNeighbors(level *Level, pos Position) []Position {
	neighbors := make([]Position, 0, 4)
	left := Position{pos.X - 1, pos.Y}
	right := Position{pos.X + 1, pos.Y}
	up := Position{pos.X, pos.Y - 1}
	down := Position{pos.X, pos.Y + 1}

	if canMove(left, level) {
		neighbors = append(neighbors, left)
	}
	if canMove(right, level) {
		neighbors = append(neighbors, right)
	}
	if canMove(up, level) {
		neighbors = append(neighbors, up)
	}
	if canMove(down, level) {
		neighbors = append(neighbors, down)
	}

	return neighbors
}
