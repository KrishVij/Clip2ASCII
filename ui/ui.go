package main

import (
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/inconsolata"
)

type Game struct {
	ui *ebitenui.UI
}

func NewGame() *Game {

	button1 := widget.NewButton(

		widget.ButtonOpts.Text("ToASCII", text.NewGoXFace(inconsolata.Bold8x16), &widget.ButtonTextColor{

			Idle:    colornames.White,
			Hover:   colornames.Azure,
			Pressed: colornames.Olive,
		}),

		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Left:   10,
			Right:  10,
			Bottom: 10,
		}),

		widget.ButtonOpts.TextPosition(

			widget.TextPositionCenter,
			widget.TextPositionCenter,
		),
		widget.ButtonOpts.Image(&widget.ButtonImage{

			Idle:    DefaultPill(colornames.Darkslategray),
			Hover:   DefaultPill(Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4)),
			Pressed: PressedPill(Mix(colornames.Darkslategray, colornames.Black, 0.4)),
		}),
		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	button2 := widget.NewButton(

		widget.ButtonOpts.Text("ToASCII", text.NewGoXFace(inconsolata.Bold8x16), &widget.ButtonTextColor{

			Idle:    colornames.White,
			Hover:   colornames.Azure,
			Pressed: colornames.Olive,
		}),
		
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    10,
			Left:   10,
			Right:  10,
			Bottom: 10,
		}),

		widget.ButtonOpts.TextPosition(

			widget.TextPositionCenter,
			widget.TextPositionCenter,
		),
		widget.ButtonOpts.Image(&widget.ButtonImage{

			Idle:    DefaultPill(colornames.Darkslategray),
			Hover:   DefaultPill(Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4)),
			Pressed: PressedPill(Mix(colornames.Darkslategray, colornames.Black, 0.4)),
		}),
		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	root := widget.NewContainer(

		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(colornames.Black),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	root.AddChild(button1)
	root.AddChild(button2)

	return &Game{

		ui: &ebitenui.UI{Container: root},
	}

}

func (g *Game) Update() error {

	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.ui.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func DefaultPill(base color.Color) *image.NineSlice {

	img := ebiten.NewImage(64, 64)
	img.Fill(base)
	return image.NewNineSliceBorder(img, 16)

}

func PressedPill(base color.Color) *image.NineSlice {

	img := ebiten.NewImage(64, 64)
	img.Fill(Mix(base, colornames.Black, 0.2))
	return image.NewNineSliceBorder(img, 16)

}

func Mix(a, b color.Color, percent float64) color.Color {
	rgba := func(c color.Color) (r, g, b, a uint8) {
		r16, g16, b16, a16 := c.RGBA()
		return uint8(r16 >> 8), uint8(g16 >> 8), uint8(b16 >> 8), uint8(a16 >> 8)
	}
	lerp := func(x, y uint8) uint8 {
		return uint8(math.Round(float64(x) + percent*(float64(y)-float64(x))))
	}
	r1, g1, b1, a1 := rgba(a)
	r2, g2, b2, a2 := rgba(b)

	return color.RGBA{
		R: lerp(r1, r2),
		G: lerp(g1, g2),
		B: lerp(b1, b2),
		A: lerp(a1, a2),
	}
}

func RoundedRectPath(x, y, w, h, tl, tr, br, bl float32) *vector.Path {
	path := &vector.Path{}

	path.Arc(x+w-tr, y+tr, tr, 3*math.Pi/2, 0, vector.Clockwise)
	path.LineTo(x+w, y+h-br)
	path.Arc(x+w-br, y+h-br, br, 0, math.Pi/2, vector.Clockwise)
	path.LineTo(x+bl, y+h)
	path.Arc(x+bl, y+h-bl, bl, math.Pi/2, math.Pi, vector.Clockwise)
	path.LineTo(x, y+tl)
	path.Arc(x+tl, y+tl, tl, math.Pi, 3*math.Pi/2, vector.Clockwise)
	path.Close()

	return path
}

func main() {

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewGame()); err != nil {

		panic(err)
	}
}
