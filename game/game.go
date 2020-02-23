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
	GetTextureAtlas() *sdl.Texture
	NewCharacterLabel(character *Character)
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
	WallSW      rune = 'F'
	WallNS      rune = 'I'
	WallNW      rune = 'L'
	WallNE      rune = 'J'
	WallSE      rune = 'T'
	WallN       rune = 'N'
	WallS       rune = 'S'
	WallE       rune = 'E'
	WallW       rune = 'W'
	WallSWE     rune = 'X'
	Floor       rune = '.'
	ClosedDoor  rune = '|'
	OpenDoor    rune = '/'
	ClosedChest rune = '*'
	OpenChest   rune = 'o'
)

type Path []Position

func bfs(ui GameUI, level *Level, startPos Position) {
	frontier := make([]Position, 0, 8)
	frontier = append(frontier, startPos)
	visited := make(map[Position]bool)
	visited[startPos] = true
	level.Debug = visited
	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		ns, _ := getNeighbors(level, current)
		for _, next := range ns {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
}

func astar(level *Level, start Position, goal Position) Path {
	var path Path
	path = make(Path, 0)

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
				//level.Debug[p] = true
				path = append(path, p)
				p = cameFrom[p]
			}
			//path = append(path, p)
			//level.Debug[p] = true

			// Reverse slice
			for i := len(path)/2 - 1; i >= 0; i-- {
				opp := len(path) - 1 - i
				path[i], path[opp] = path[opp], path[i]
			}
			return path
		}

		frontier = frontier[1:]
		ns, _ := getNeighbors(level, current.Position)
		for _, next := range ns {
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
	return nil
}

func Run(gameUI GameUI) {
	level := LoadLevelFromFile("game/maps/level1.map")
	level.Player = NewPlayer("wkorande", 5.0, Position{5, 5}, gameUI.GetTextureAtlas(), gameUI.GetTextureIndex('@'))
	gameUI.NewCharacterLabel(&level.Player.Character)

	level.Enemies = make([]*Enemy, 0)
	enemy := NewEnemy("bocal", 1.0, Position{5, 13}, gameUI.GetTextureAtlas(), gameUI.GetTextureIndex('E'))
	level.Enemies = append(level.Enemies, enemy)
	gameUI.NewCharacterLabel(&enemy.Character)

	for {
		level.Debug = make(map[Position]bool)
		for _, e := range level.Enemies {
			e.Update(level)
		}
		gameUI.Draw(level)
		input := gameUI.GetInput()
		if input.Type == Quit {
			return
		}
		if input.Type == Search {
			level.Debug = make(map[Position]bool)
			//astar(level, level.Player.Pos, getRandomPositionInsideCircle(5, level.Player.Pos))
		}

		handleInput(level, input)

	}
}
