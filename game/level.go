package game

import (
	"bufio"
	"math"
	"math/rand"
	"os"
)

type Level struct {
	Map     [][]Tile
	Player  *Player
	Enemies []*Enemy
	Width   int
	Height  int
	Debug   map[Position]bool
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
	level.Width = cols
	level.Height = rows
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
	checkLevelWallCorners(level)
	return level
}

func checkLevelWallCorners(level *Level) {
	for y, rows := range level.Map {
		for x, _ := range rows {
			if isWall(level, Position{x, y}) {
				flags := getWallNeighbors(level, Position{x, y})
				if flags == 10 {
					level.Map[y][x].Rune = WallSW
				} else if flags == 12 {
					level.Map[y][x].Rune = WallNS
				} else if flags == 6 {
					level.Map[y][x].Rune = WallNW
				} else if flags == 5 {
					level.Map[y][x].Rune = WallNE
				} else if flags == 9 {
					level.Map[y][x].Rune = WallSE
				} else if flags == 4 {
					level.Map[y][x].Rune = WallN
				} else if flags == 8 {
					level.Map[y][x].Rune = WallS
				} else if flags == 11 {
					level.Map[y][x].Rune = WallSWE
				} else if flags == 15 {
					level.Map[y][x].Rune = WallSWE
				}

			}
		}
	}
}

func getWallNeighbors(level *Level, pos Position) uint8 {
	var flags uint8
	left := Position{pos.X - 1, pos.Y}
	right := Position{pos.X + 1, pos.Y}
	up := Position{pos.X, pos.Y - 1}
	down := Position{pos.X, pos.Y + 1}

	if left.X >= 0 && left.X < level.Width && isWall(level, left) {
		flags |= 1 << 0
	}
	if right.X >= 0 && right.X < level.Width && isWall(level, right) {
		flags |= 1 << 1
	}
	if up.Y >= 0 && up.Y < level.Height && isWall(level, up) {
		flags |= 1 << 2
	}
	if down.Y >= 0 && down.Y < level.Height && isWall(level, down) {
		flags |= 1 << 3
	}

	return flags
}

func isWall(level *Level, pos Position) bool {
	if level.Map[pos.Y][pos.X].Rune == Wall ||
		level.Map[pos.Y][pos.X].Rune == WallSW ||
		level.Map[pos.Y][pos.X].Rune == WallNW ||
		level.Map[pos.Y][pos.X].Rune == WallNS ||
		level.Map[pos.Y][pos.X].Rune == WallNE ||
		level.Map[pos.Y][pos.X].Rune == WallSE ||
		level.Map[pos.Y][pos.X].Rune == WallN ||
		level.Map[pos.Y][pos.X].Rune == WallS ||
		level.Map[pos.Y][pos.X].Rune == WallE ||
		level.Map[pos.Y][pos.X].Rune == WallW ||
		level.Map[pos.Y][pos.X].Rune == WallSWE {
		return true
	}
	return false
}

func hasEnemy(pos Position, level *Level) (bool, *Enemy) {
	for _, e := range level.Enemies {
		if pos == e.Pos {
			return true, e
		}
	}
	return false, nil
}

func canMove(pos Position, level *Level) bool {
	if isWall(level, pos) || level.Map[pos.Y][pos.X].Rune == ClosedDoor || level.Map[pos.Y][pos.X].Rune == ClosedChest {
		return false
	}
	exists, _ := hasEnemy(pos, level)
	if exists {
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

func getNeighbors(level *Level, pos Position) ([]Position, uint8) {
	var flags uint8
	neighbors := make([]Position, 0, 4)
	left := Position{pos.X - 1, pos.Y}
	right := Position{pos.X + 1, pos.Y}
	up := Position{pos.X, pos.Y - 1}
	down := Position{pos.X, pos.Y + 1}

	if canMove(left, level) {
		neighbors = append(neighbors, left)
		flags |= 1 << 0
	}
	if canMove(right, level) {
		neighbors = append(neighbors, right)
		flags |= 1 << 1
	}
	if canMove(up, level) {
		neighbors = append(neighbors, up)
		flags |= 1 << 2
	}
	if canMove(down, level) {
		neighbors = append(neighbors, down)
		flags |= 1 << 3
	}

	return neighbors, flags
}

func getRandomPositionInsideCircle(radius int, pos Position) Position {
	var p Position
	angle := 2.0 * math.Pi * rand.Float64()
	r := float64(radius) * math.Sqrt(rand.Float64())
	p.X = int(r*math.Cos(angle) + float64(pos.X))
	p.Y = int(r*math.Sin(angle) + float64(pos.Y))
	return p
}
