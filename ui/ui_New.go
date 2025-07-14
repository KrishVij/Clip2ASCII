package main

import (
	"bytes"
	"fmt"

	// stdImage "image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	ui  *ebitenui.UI
	btn *widget.Button
}

func loadImageInvisible() (*widget.ButtonImage, error) {

	idle := image.NewBorderedNineSliceColor(color.NRGBA{R: 144, G: 122, B: 169, A: 255}, color.NRGBA{144, 122, 169, 255}, 3)

	hover := image.NewBorderedNineSliceColor(color.NRGBA{R: 144, G: 122, B: 169, A: 255}, color.NRGBA{144, 122, 169, 255}, 3)

	pressed := image.NewAdvancedNineSliceColor(color.NRGBA{R: 144, G: 122, B: 169, A: 255}, image.NewBorder(3, 2, 2, 2, color.NRGBA{246, 193, 119, 255}))

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadImage() (*widget.ButtonImage, error) {

	idle := image.NewBorderedNineSliceColor(color.NRGBA{R: 234, G: 157, B: 52, A: 255}, color.NRGBA{234, 157, 52, 255}, 3)

	hover := image.NewBorderedNineSliceColor(color.NRGBA{R: 234, G: 157, B: 52, A: 255}, color.NRGBA{234, 157, 52, 255}, 3)

	pressed := image.NewAdvancedNineSliceColor(color.NRGBA{R: 246, G: 193, B: 119, A: 255}, image.NewBorder(3, 2, 2, 2, color.NRGBA{246, 193, 119, 255}))

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (text.Face, error) {

	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("Error loading font: %w", err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func main() {

	rootContainer := widget.NewContainer(

		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{

			0x90, 0x7a, 0xa9, 0xff,
		})),

		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	ButtonImageInvisible, _ := loadImageInvisible()
	ButtonImage, _ := loadImage()
	Face, _ := loadFont(30)

	buttonGroup1 := widget.NewContainer(

		widget.ContainerOpts.Layout(widget.NewRowLayout(
			
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
		
		widget.ContainerOpts.WidgetOpts(
			
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	buttonGroup2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Spacing(0))),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	btn_CHTOFD:= widget.NewButton(

		widget.ButtonOpts.Image(ButtonImage),

		widget.ButtonOpts.Text("ClickHereToOpenFileDialog", Face, &widget.ButtonTextColor{

			Idle: color.NRGBA{0, 0, 0, 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(false),

		widget.ButtonOpts.TextPadding(widget.Insets{

			Left:   60,
			Right:  60,
			Top:    10,
			Bottom: 10,
		}),

		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
		),
	)

	btn_Invisible := widget.NewButton(

		widget.ButtonOpts.Image(ButtonImageInvisible),

		widget.ButtonOpts.Text("ToASCII", Face, &widget.ButtonTextColor{

			Idle: color.NRGBA{144, 122, 169, 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(false),

		widget.ButtonOpts.TextPadding(widget.Insets{

			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	btn_ToASCII := widget.NewButton(

		widget.ButtonOpts.Image(ButtonImage),

		widget.ButtonOpts.Text("ToASCII", Face, &widget.ButtonTextColor{

			Idle: color.NRGBA{0, 0, 0, 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(false),

		widget.ButtonOpts.TextPadding(widget.Insets{

			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	btn_FromASCII := widget.NewButton(

		widget.ButtonOpts.Image(ButtonImage),

		widget.ButtonOpts.Image(ButtonImage),

		widget.ButtonOpts.Text("FromASCII", Face, &widget.ButtonTextColor{

			Idle: color.NRGBA{0, 0, 0, 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(false),

		widget.ButtonOpts.TextPadding(widget.Insets{

			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				HorizontalPosition: widget.AnchorLayoutPositionStart,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	buttonGroup1.AddChild(btn_CHTOFD)
	rootContainer.AddChild(buttonGroup1)
	buttonGroup2.AddChild(btn_ToASCII)
	buttonGroup2.AddChild(btn_Invisible)
	buttonGroup2.AddChild(btn_FromASCII)
	// rootContainer.AddChild(buttonGroup1)
	rootContainer.AddChild(buttonGroup2)

	ui := ebitenui.UI{

		Container: rootContainer,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Clip2ASCII")

	game := Game{
		ui:  &ui,
		btn: btn_ToASCII,
	}

	// game.ButtonTraversal(ButtonImageInvisible)

	err := ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
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
