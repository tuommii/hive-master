package game

import (
	"bufio"
	"math"
	"os"
	"sort"
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
	Search
	Quit
)

type Input struct {
	Type InputType
}

type Position struct {
	X int
	Y int
}

type Tile struct {
	Rune rune
}

type priorityPos struct {
	Position
	priority int
}

type priorityArray []priorityPos

func (p priorityArray) Len() int           { return len(p) }
func (p priorityArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p priorityArray) Less(i, j int) bool { return p[i].priority < p[j].priority }

const (
	Blank       rune = ' '
	Wall        rune = '#'
	Floor       rune = '.'
	ClosedDoor  rune = '|'
	OpenDoor    rune = '/'
	ClosedChest rune = '*'
	OpenChest   rune = 'o'
)

type Level struct {
	Map    [][]Tile
	Player Player
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

func bfs(ui GameUI, level *Level, startPos Position) {
	frontier := make([]Position, 0, 8)
	frontier = append(frontier, startPos)
	visited := make(map[Position]bool)
	visited[startPos] = true
	level.Debug = visited
	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
}

func astar(ui GameUI, level *Level, start Position, goal Position) {
	frontier := make(priorityArray, 0, 0)
	frontier = append(frontier, priorityPos{start, 1})
	cameFrom := make(map[Position]Position)
	cameFrom[start] = start
	costSoFar := make(map[Position]int)
	costSoFar[start] = 0
	for len(frontier) > 0 {
		sort.Stable(frontier)
		current := frontier[0]
		//level.Debug[current.Position] = true
		if current.Position == goal {
			p := current.Position
			for p != start {
				level.Debug[p] = true
				p = cameFrom[p]
			}
			level.Debug[p] = true
			break
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current.Position) {
			newCost := costSoFar[current.Position] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = append(frontier, priorityPos{next, priority})
				cameFrom[next] = current.Position
				//level.Debug[next] = true
			}
		}
	}
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

func handleInput(level *Level, input *Input) {
	toPos := level.Player.Pos
	switch input.Type {
	case Up:
		toPos.Y--
	case Down:
		toPos.Y++
	case Left:
		toPos.X--
	case Right:
		toPos.X++
	}
	if canMove(toPos, level) {
		level.Player.Move(toPos)
	} else {
		checkDoor(toPos, level)
	}
}

func Run(ui GameUI) {
	level := LoadLevelFromFile("game/maps/level1.map")
	level.Debug = make(map[Position]bool)
	level.Player.Pos.X = 5
	level.Player.Pos.Y = 5
	for {
		ui.Draw(level)
		input := ui.GetInput()
		if input.Type == Quit {
			return
		}
		if input.Type == Search {
			level.Debug = make(map[Position]bool)
			astar(ui, level, level.Player.Pos, Position{9, 4})
		}
		handleInput(level, input)
	}
}
