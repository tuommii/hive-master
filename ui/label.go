package ui

import (
	"fmt"

	"github.com/wehard/hive-master/game"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Label interface {
	SetText(text string)
	Draw(pos game.Position)
}

type label struct {
	text string
	r    *sdl.Renderer
	font *ttf.Font
}

func NewLabel(text string, r *sdl.Renderer) *label {
	font, err := ttf.OpenFont("ui/assets/anonymous_pro.ttf", 20)
	if err != nil {
		panic(err)
	}
	return &label{text, r, font}
}

func (l *label) Draw(pos game.Position) {
	s, err := l.font.RenderUTF8Shaded(l.text, sdl.Color{255, 255, 255, 255}, sdl.Color{0, 0, 0, 255})
	if err != nil {
		fmt.Println("failed to create font surface:", err)
		return
	}
	defer s.Free()
	var clipRect sdl.Rect
	s.GetClipRect(&clipRect)
	texture, err := l.r.CreateTextureFromSurface(s)
	if err != nil {
		fmt.Println("failed to create font texture from surface:", err)
		return
	}
	destRect := &sdl.Rect{
		X: int32(pos.X),
		Y: int32(pos.Y),
		W: clipRect.W,
		H: clipRect.H,
	}
	l.r.Copy(texture, nil, destRect)
}

func (l *label) SetText(text string) {
	l.text = text
}
