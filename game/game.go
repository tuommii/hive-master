package game

import (
	"math"
	"sort"

	"github.com/wehard/ftapi"

	"github.com/veandco/go-sdl2/sdl"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
	GetTextureIndex(TileType) *sdl.Rect
	GetTextureAtlas() *sdl.Texture
	NewCharacterLabel(character *Character)
}

type Position struct {
	X int
	Y int
}

type TileType string

type Tile struct {
	TileType TileType
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
	Blank       TileType = "blank"
	Wall        TileType = "wall"
	WallEW      TileType = "wall_ew"
	WallSW      TileType = "wall_sw"
	WallNS      TileType = "wall_ns"
	WallNW      TileType = "wall_nw"
	WallNE      TileType = "wall_ne"
	WallSE      TileType = "wall_se"
	WallN       TileType = "wall_n"
	WallS       TileType = "wall_s"
	WallE       TileType = "wall_e"
	WallW       TileType = "wall_w"
	WallSWE     TileType = "wall_swe"
	WallNSW     TileType = "wall_nsw"
	WallNSE     TileType = "wall_nse"
	WallNWE     TileType = "wall_nwe"
	Floor       TileType = "floor"
	Hole        TileType = "hole"
	ClosedDoorV TileType = "door_closed_v"
	OpenDoorV   TileType = "door_open_v"
	ClosedDoorH TileType = "door_closed_h"
	OpenDoorH   TileType = "door_open_h"
	ClosedChest TileType = "chest_closed"
	OpenChest   TileType = "chest_open"
)

type Path []Position

func BreadthFirstSearch(level *Level, startPos Position) map[Position]bool {
	frontier := make([]Position, 0, 8)
	frontier = append(frontier, startPos)
	visited := make(map[Position]bool)
	visited[startPos] = true
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
	return visited
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

func get_some_key(m map[string]ftapi.UserData) string {
	for k := range m {
		return k
	}
	return ""
}

func checkVisibility(level *Level, character *Character) {
	v := BreadthFirstSearch(level, character.Pos)

	for p := range v {
		level.Visible[p.Y][p.X] = true
	}
}

func Run(gameUI GameUI) {
	userData, _ := ftapi.LoadUserData("game/users.json")
	level := LoadLevelFromCSVFile("ui/assets/dungeon_csv_Wall.csv")

	playerUser := userData["wkorande"]

	level.Player = NewPlayer(playerUser.Login, playerUser.CursusUsers[0].Level, level.getRandomPosition(), gameUI.GetTextureAtlas(), gameUI.GetTextureIndex("player"))
	gameUI.NewCharacterLabel(&level.Player.Character)

	level.Enemies = make([]*Enemy, 0)
	for i := 0; i < 50; i++ {
		user := userData[get_some_key(userData)]
		pos := level.getRandomPosition()
		enemy := NewEnemy(user.Login, user.CursusUsers[0].Level, pos, gameUI.GetTextureAtlas(), gameUI.GetTextureIndex("enemy"))
		level.Enemies = append(level.Enemies, enemy)
		gameUI.NewCharacterLabel(&enemy.Character)
	}

	for {
		// Clear dead enemies
		for i := len(level.Enemies) - 1; i >= 0; i-- {
			if level.Enemies[i].IsDead {
				level.Enemies = append(level.Enemies[:i], level.Enemies[i+1:]...)
			}
		}

		// Update enemies
		for _, e := range level.Enemies {
			if !e.IsDead {
				e.Update(level)
			}
		}

		// Check visibility
		level.resetVisibility(false)
		checkVisibility(level, &level.Player.Character)

		gameUI.Draw(level)
		input := gameUI.GetInput()
		if input.Type == Quit {
			return
		}
		if input.Type == Action {
			checkHole(level.Player.Pos, level)
			//level.Debug = make(map[Position]bool)
			//astar(level, level.Player.Pos, getRandomPositionInsideCircle(5, level.Player.Pos))
		}
		handleInput(level, input)
	}
}
