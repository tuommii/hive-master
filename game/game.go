package game

import (
	"math"
	"sort"

	"github.com/veandco/go-sdl2/sdl"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
	GetTextureIndex(rune) *sdl.Rect
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

func astar(level *Level, start Position, goal Position) {
	frontier := make(priorityArray, 0, 0)
	frontier = append(frontier, priorityPos{start, 1})
	cameFrom := make(map[Position]Position)
	cameFrom[start] = start
	costSoFar := make(map[Position]int)
	costSoFar[start] = 0
	for len(frontier) > 0 {
		sort.Stable(frontier)
		current := frontier[0]
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
			}
		}
	}
}

func Run(gameUI GameUI) {
	level := LoadLevelFromFile("game/maps/level1.map")
	level.Player = NewPlayer("name", 5.0, Position{5, 5}, gameUI.GetTextureIndex('@'))
	level.Debug = make(map[Position]bool)
	level.Player.Pos.X = 5
	level.Player.Pos.Y = 5
	for {
		gameUI.Draw(level)
		input := gameUI.GetInput()
		if input.Type == Quit {
			return
		}
		if input.Type == Search {
			level.Debug = make(map[Position]bool)
			astar(level, level.Player.Pos, Position{9, 4})
		}
		handleInput(level, input)
	}
}
