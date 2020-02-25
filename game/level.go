package game

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Level struct {
	Map     [][]Tile
	Visible [][]bool
	Player  *Player
	Enemies []*Enemy
	Width   int
	Height  int
	Debug   map[Position]bool
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
	level.Visible = make([][]bool, rows)
	for i := range level.Map {
		level.Map[i] = make([]Tile, cols)
		level.Visible[i] = make([]bool, cols)
	}
	for y := 0; y < rows; y++ {
		line := strings.Split(levelLines[y], ",")
		for x := 0; x < len(line); x++ {
			c, _ := strconv.Atoi(line[x])
			var t Tile
			switch c {
			case -2:
				t.TileType = Blank
			case 8:
				t.TileType = Hole
			case 42:
				t.TileType = WallS
			case 65, 129:
				t.TileType = Wall
			case 64:
				t.TileType = WallSW
			case 67:
				t.TileType = WallSW
			case 66:
				t.TileType = WallSE
			case 68:
				t.TileType = WallNWE
			case 69:
				t.TileType = WallNSE
			case 73:
				t.TileType = WallE
			case 74:
				t.TileType = WallSWE
			case 75:
				t.TileType = WallW
			case 96, 98:
				t.TileType = WallNS
			case 106:
				t.TileType = WallN
			case 128:
				t.TileType = WallNW
			case 130:
				t.TileType = WallNE
			case -1:
				t.TileType = Floor
			case 71:
				t.TileType = ClosedDoorH
			case 102:
				t.TileType = ClosedDoorV
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
	//checkLevelWallOrientation(level)
	//checkLevelDoors(level)
	return level
}

func (level *Level) resetVisibility(v bool) {
	for y := 0; y < level.Height; y++ {
		for x := 0; x < level.Width; x++ {
			level.Visible[y][x] = v
		}
	}
}

func checkLevelDoorOrientation(level *Level) {
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

func checkLevelWallOrientation(level *Level) {
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
		level.Map[pos.Y][pos.X].TileType == WallN ||
		level.Map[pos.Y][pos.X].TileType == WallW ||
		level.Map[pos.Y][pos.X].TileType == WallS ||
		level.Map[pos.Y][pos.X].TileType == WallE ||
		level.Map[pos.Y][pos.X].TileType == WallNS ||
		level.Map[pos.Y][pos.X].TileType == WallNE ||
		level.Map[pos.Y][pos.X].TileType == WallNW ||
		level.Map[pos.Y][pos.X].TileType == WallSW ||
		level.Map[pos.Y][pos.X].TileType == WallSE ||
		level.Map[pos.Y][pos.X].TileType == WallSWE ||
		level.Map[pos.Y][pos.X].TileType == WallNSE ||
		level.Map[pos.Y][pos.X].TileType == WallNSW ||
		level.Map[pos.Y][pos.X].TileType == WallNWE {
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

func checkHole(pos Position, level *Level) {
	t := level.getTileType(pos)
	if t == Hole {
		fmt.Println("you check the hole")
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

func (level *Level) getTileType(pos Position) TileType {
	return level.Map[pos.Y][pos.X].TileType
}

func (level *Level) getRandomPosition() Position {
	t := time.Now().Nanosecond()
	rand.Seed(int64(t))
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
