package game

import (
	"bufio"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
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
				t.TileType = Blank
			case '#':
				t.TileType = Wall
			case '.':
				t.TileType = Floor
			case '|':
				t.TileType = ClosedDoorV
			case '/':
				t.TileType = OpenDoorV
			case '-':
				t.TileType = ClosedDoorH
			case '\\':
				t.TileType = ClosedDoorH
			case '*':
				t.TileType = ClosedChest
			case 'o':
				t.TileType = OpenChest
			default:
				t.TileType = Blank
			}
			level.Map[y][x] = t
		}
	}
	checkLevelWallCorners(level)
	checkLevelDoors(level)
	return level
}

func LoadLevelFromCSVFile(filename string) *Level {
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
		l := strings.Split(line, ",")
		if len(l) > cols {
			cols = len(l)
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
		line := strings.Split(levelLines[y], ",")
		for x := 0; x < len(line); x++ {
			c, _ := strconv.Atoi(line[x])
			var t Tile
			switch c {
			case -2:
				t.TileType = Blank
			case 64, 65, 66, 73, 74, 75, 96, 98, 128, 129, 130:
				t.TileType = Wall
			case -1:
				t.TileType = Floor
			case 102, 71:
				t.TileType = ClosedDoorV
			case 104, 135:
				t.TileType = OpenDoorV
			case 224:
				t.TileType = ClosedChest
			case 226:
				t.TileType = OpenChest
			default:
				t.TileType = Blank
			}
			level.Map[y][x] = t
		}
	}
	checkLevelWallCorners(level)
	checkLevelDoors(level)
	return level
}

func checkLevelDoors(level *Level) {
	for y, rows := range level.Map {
		for x, _ := range rows {
			if isDoor(level, Position{x, y}) {
				_, flags := getNeighbors(level, Position{x, y})
				if flags == 12 {
					level.Map[y][x].TileType = ClosedDoorH
				} else if flags == 11 {
					level.Map[y][x].TileType = ClosedDoorH
				}
			}
		}
	}
}

func checkLevelWallCorners(level *Level) {
	for y, rows := range level.Map {
		for x, _ := range rows {
			if isWall(level, Position{x, y}) {
				flags := getWallNeighbors(level, Position{x, y})
				if flags == 10 {
					level.Map[y][x].TileType = WallSW
				} else if flags == 12 {
					level.Map[y][x].TileType = WallNS
				} else if flags == 6 {
					level.Map[y][x].TileType = WallNW
				} else if flags == 5 {
					level.Map[y][x].TileType = WallNE
				} else if flags == 9 {
					level.Map[y][x].TileType = WallSE
				} else if flags == 4 {
					level.Map[y][x].TileType = WallN
				} else if flags == 8 {
					level.Map[y][x].TileType = WallS
				} else if flags == 11 {
					level.Map[y][x].TileType = WallSWE
				} else if flags == 15 {
					level.Map[y][x].TileType = WallSWE
				} else if flags == 14 {
					level.Map[y][x].TileType = WallNSW
				} else if flags == 13 {
					level.Map[y][x].TileType = WallNSE
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
	if pos.X < 0 || pos.X > level.Width-1 {
		//fmt.Println("X position out of level bounds!", "pos", pos, "width", level.Width)
		return true
	}
	if pos.Y < 0 || pos.Y > level.Height-1 {
		//fmt.Println("Y position out of level bounds!", "pos", pos, "height", level.Height)
		return true
	}
	if level.Map[pos.Y][pos.X].TileType == Wall ||
		level.Map[pos.Y][pos.X].TileType == WallSW ||
		level.Map[pos.Y][pos.X].TileType == WallNW ||
		level.Map[pos.Y][pos.X].TileType == WallNS ||
		level.Map[pos.Y][pos.X].TileType == WallNE ||
		level.Map[pos.Y][pos.X].TileType == WallSE ||
		level.Map[pos.Y][pos.X].TileType == WallN ||
		level.Map[pos.Y][pos.X].TileType == WallS ||
		level.Map[pos.Y][pos.X].TileType == WallE ||
		level.Map[pos.Y][pos.X].TileType == WallW ||
		level.Map[pos.Y][pos.X].TileType == WallSWE ||
		level.Map[pos.Y][pos.X].TileType == WallNSE ||
		level.Map[pos.Y][pos.X].TileType == WallNSW {
		return true
	}
	return false
}

func isDoor(level *Level, pos Position) bool {
	if level.Map[pos.Y][pos.X].TileType == ClosedDoorV ||
		level.Map[pos.Y][pos.X].TileType == OpenDoorV ||
		level.Map[pos.Y][pos.X].TileType == ClosedDoorH ||
		level.Map[pos.Y][pos.X].TileType == OpenDoorH {
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
	if isWall(level, pos) ||
		level.Map[pos.Y][pos.X].TileType == ClosedDoorV ||
		level.Map[pos.Y][pos.X].TileType == ClosedDoorH ||
		level.Map[pos.Y][pos.X].TileType == ClosedChest {
		return false
	}
	exists, _ := hasEnemy(pos, level)
	if exists {
		return false
	}
	return true
}

func checkDoor(pos Position, level *Level) {
	if level.Map[pos.Y][pos.X].TileType == ClosedDoorV {
		level.Map[pos.Y][pos.X].TileType = OpenDoorV
	} else if level.Map[pos.Y][pos.X].TileType == OpenDoorV {
		level.Map[pos.Y][pos.X].TileType = ClosedDoorV
	}
	if level.Map[pos.Y][pos.X].TileType == ClosedDoorH {
		level.Map[pos.Y][pos.X].TileType = OpenDoorH
	} else if level.Map[pos.Y][pos.X].TileType == OpenDoorH {
		level.Map[pos.Y][pos.X].TileType = ClosedDoorH
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

func isBlank(level *Level, pos Position) bool {
	if level.Map[pos.Y][pos.X].TileType == Blank {
		return true
	}
	return false
}

func (level *Level) getRandomPosition() Position {
	pos := Position{-1, -1}
	for pos.X < 0 || pos.Y < 0 || isWall(level, pos) || isBlank(level, pos) {
		e, _ := hasEnemy(pos, level)
		if e {
			pos = Position{-1, -1}
			continue
		}
		pos.X = rand.Intn(level.Width - 1)
		pos.Y = rand.Intn(level.Height - 1)
	}
	return pos
}
