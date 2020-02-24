package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/wehard/hive-master/game"
)

type UI2d struct {
	WindowTitle string
}

const (
	winWidth, winHeight = 1920, 1080
)

var renderer *sdl.Renderer
var window *sdl.Window
var textureAtlas *sdl.Texture
var textureIndex map[game.TileType]sdl.Rect
var keyboardState []uint8
var prevKeyboardState []uint8
var centerX int
var centerY int
var offsetX int32
var offsetY int32
var characterLabels map[*game.Character]Label
var tileSize int32 = 32

func (ui *UI2d) GetTextureIndex(tileType game.TileType) *sdl.Rect {
	i := textureIndex[tileType]
	return &i
}

func (ui *UI2d) GetTextureAtlas() *sdl.Texture {
	return textureAtlas
}

func (ui *UI2d) NewCharacterLabel(character *game.Character) {
	characterLabels[character] = NewLabel(character.Name, renderer)
}

func init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("Hive Master", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	textureAtlas, err = img.LoadTexture(renderer, "ui/assets/dungeon.png")
	if err != nil {
		fmt.Println(err)
	}
	textureAtlas.SetBlendMode(sdl.BLENDMODE_BLEND)
	loadTextureIndex("ui/assets/texture_index.txt")

	keyboardState = sdl.GetKeyboardState()
	prevKeyboardState = make([]uint8, len(keyboardState))
	for i, v := range keyboardState {
		prevKeyboardState[i] = v
	}
	centerX = -1
	centerY = -1
	characterLabels = make(map[*game.Character]Label)
}

func loadTextureIndex(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	textureIndex = make(map[game.TileType]sdl.Rect)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		var tile game.Tile
		tile.TileType = getTileType(rune(line[0]))
		xy := line[1:]
		split := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(split[0]), 10, 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseInt(strings.TrimSpace(split[1]), 10, 64)
		if err != nil {
			panic(err)
		}
		w, err := strconv.ParseInt(strings.TrimSpace(split[2]), 10, 64)
		if err != nil {
			panic(err)
		}
		h, err := strconv.ParseInt(strings.TrimSpace(split[3]), 10, 64)
		if err != nil {
			panic(err)
		}
		tileIndexX := x
		tileIndexY := y
		tileRect := sdl.Rect{X: int32(tileIndexX * 16), Y: int32(tileIndexY * 16), W: int32(w), H: int32(h)}
		textureIndex[tile.TileType] = tileRect
	}
}

func getTileType(r rune) game.TileType {
	return game.TileType(r)
}

func (ui UI2d) Draw(level *game.Level) {
	if centerX == -1 && centerY == -1 {
		centerX = level.Player.Pos.X
		centerY = level.Player.Pos.Y
	}
	moveThreshold := 4
	if level.Player.Pos.X > centerX+moveThreshold {
		centerX++
	} else if level.Player.Pos.X < centerX-moveThreshold {
		centerX--
	} else if level.Player.Pos.Y > centerY+moveThreshold {
		centerY++
	} else if level.Player.Pos.Y < centerY-moveThreshold {
		centerY--
	}

	offsetX = int32((winWidth / 2) - int32(centerX)*tileSize)
	offsetY = int32((winHeight / 2) - int32(centerY)*tileSize)
	renderer.Clear()
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.TileType != game.Blank {
				srcRect := textureIndex[level.Map[y][x].TileType]
				destRect := sdl.Rect{
					X: int32(x*int(tileSize)) + offsetX,
					Y: int32(y*int(tileSize)) + offsetY,
					W: tileSize,
					H: tileSize,
				}
				pos := game.Position{x, y}
				if level.Debug[pos] {
					textureAtlas.SetColorMod(128, 255, 128)
				} else {
					textureAtlas.SetColorMod(255, 255, 255)
				}
				floorRect := textureIndex[game.Floor]
				renderer.Copy(textureAtlas, &floorRect, &destRect)
				renderer.Copy(textureAtlas, &srcRect, &destRect)
			}
		}
	}
	for _, enemy := range level.Enemies {
		if !enemy.IsDead {
			enemy.Draw(renderer, tileSize, offsetX, offsetY)
			label := characterLabels[&enemy.Character]
			label.Draw(enemy.Pos)
		}
	}
	level.Player.Draw(renderer, tileSize, offsetX, offsetY)
	label := characterLabels[&level.Player.Character]
	label.Draw(level.Player.Pos)
	renderer.Present()
}

func (ui *UI2d) GetInput() *game.Input {
	for {
		var input game.Input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return &game.Input{Type: game.Quit}
			}
		}
		if keyboardState[sdl.SCANCODE_ESCAPE] == 1 && prevKeyboardState[sdl.SCANCODE_ESCAPE] == 0 {
			input.Type = game.Quit
		} else if keyboardState[sdl.SCANCODE_UP] == 1 && prevKeyboardState[sdl.SCANCODE_UP] == 0 {
			input.Type = game.Up
		} else if keyboardState[sdl.SCANCODE_DOWN] == 1 && prevKeyboardState[sdl.SCANCODE_DOWN] == 0 {
			input.Type = game.Down
		} else if keyboardState[sdl.SCANCODE_LEFT] == 1 && prevKeyboardState[sdl.SCANCODE_LEFT] == 0 {
			input.Type = game.Left
		} else if keyboardState[sdl.SCANCODE_RIGHT] == 1 && prevKeyboardState[sdl.SCANCODE_RIGHT] == 0 {
			input.Type = game.Right
		} else if keyboardState[sdl.SCANCODE_SPACE] == 1 && prevKeyboardState[sdl.SCANCODE_SPACE] == 0 {
			input.Type = game.Action
		} else if keyboardState[sdl.SCANCODE_KP_PLUS] == 1 && prevKeyboardState[sdl.SCANCODE_KP_PLUS] == 0 {
			input.Type = game.ZoomIn
			tileSize++
		} else if keyboardState[sdl.SCANCODE_KP_MINUS] == 1 && prevKeyboardState[sdl.SCANCODE_KP_MINUS] == 0 {
			input.Type = game.ZoomOut
			tileSize--
		}
		for i, v := range keyboardState {
			prevKeyboardState[i] = v
		}
		if input.Type != game.None {
			return &input
		}
	}
}
