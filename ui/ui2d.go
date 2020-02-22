package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wehard/hive-master/game"
)

type UI2d struct {
	WindowTitle string
}

const (
	winWidth, winHeight = 1920, 1080
	tileSize            = 48
)

var renderer *sdl.Renderer
var window *sdl.Window
var textureAtlas *sdl.Texture
var textureIndex map[rune]sdl.Rect
var keyboardState []uint8
var prevKeyboardState []uint8
var centerX int
var centerY int

func init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
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

	//renderer.SetScale(2, 2)
	textureAtlas, err = img.LoadTexture(renderer, "ui/assets/dungeon_tileset.png")
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
}

func loadTextureIndex(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	textureIndex = make(map[rune]sdl.Rect)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		var tile game.Tile
		tile.Rune = rune(line[0])
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
		tileIndexX := x
		tileIndexY := y
		tileRect := sdl.Rect{X: int32(tileIndexX), Y: int32(tileIndexY), W: 16, H: 16}
		fmt.Println(tileRect)
		textureIndex[tile.Rune] = tileRect
	}
}

func (ui UI2d) Draw(level *game.Level) {
	if centerX == -1 && centerY == -1 {
		centerX = level.Player.Pos.X
		centerY = level.Player.Pos.Y
	}

	//dx := level.Player.Pos.X - centerX
	//dy := level.Player.Pos.Y - centerY

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

	offsetX := int32((winWidth / 2) - centerX*tileSize)
	offsetY := int32((winHeight / 2) - centerY*tileSize)
	renderer.Clear()
	for y, row := range level.Map {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRect := textureIndex[level.Map[y][x].Rune]
				destRect := sdl.Rect{
					X: int32(x*tileSize) + offsetX,
					Y: int32(y*tileSize) + offsetY,
					W: tileSize,
					H: tileSize,
				}
				pos := game.Position{x, y}
				if level.Debug[pos] {
					textureAtlas.SetColorMod(128, 255, 128)
				} else {
					textureAtlas.SetColorMod(255, 255, 255)
				}
				renderer.Copy(textureAtlas, &srcRect, &destRect)
			}
		}
	}
	playerSrcRect := textureIndex['@']
	playerDestRect := sdl.Rect{
		X: int32(level.Player.Pos.X*tileSize) + offsetX,
		Y: int32(level.Player.Pos.Y*tileSize) + offsetY,
		W: tileSize,
		H: tileSize,
	}
	renderer.Copy(textureAtlas, &playerSrcRect, &playerDestRect)
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
			input.Type = game.Search
		}
		for i, v := range keyboardState {
			prevKeyboardState[i] = v
		}
		if input.Type != game.None {
			return &input
		}
	}
}
