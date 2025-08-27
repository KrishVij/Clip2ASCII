package main

import (
	"bytes"
	"fmt"
	"os"
	// "os/exec"
	"time"
	"sync"
	"path/filepath"

	stdImage "image"
	"image/color"
	"log"

	"github.com/KononK/resize"
	"github.com/hajimehoshi/ebiten/v2"

	FFmpegutils "github.com/KrishVij/clip2ASCII/FFmpeg_Utils"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/ncruces/zenity"
	"golang.org/x/image/font/gofont/goregular"
	"github.com/KrishVij/clip2ASCII/Frame_Processing"
)

var my_color = color.NRGBA{R: 242, G: 233, B: 225, A: 255}
var rosePinePine color.Color = my_color
var selected_video_path string

// var newPage = image.NewImageColor(rosePinePurple)

type Game struct {
	ui                                    *ebitenui.UI
	btn                                   *widget.Button
	thumbnail                             *ebiten.Image
	thumbnail_image_image_format          stdImage.Image
	resized_thumbnail_image_image_format  stdImage.Image
	resized_thumbnail_ebiten_image_format *ebiten.Image
	btn_CHTOFD                            *widget.Button
	screen_height, screen_width           int
	thumbnail_widget                      *widget.Graphic
}

func loadImageInvisible() (*widget.ButtonImage, error) {

	idle := image.NewBorderedNineSliceColor(color.NRGBA{R: 42, G: 39, B: 63, A: 255}, color.NRGBA{42, 39, 63, 255}, 3)

	hover := image.NewBorderedNineSliceColor(color.NRGBA{R: 42, G: 39, B: 63, A: 255}, color.NRGBA{42, 39, 63, 255}, 3)

	pressed := image.NewAdvancedNineSliceColor(color.NRGBA{R: 42, G: 39, B: 63, A: 255}, image.NewBorder(3, 2, 2, 2, color.NRGBA{42, 39, 63, 255}))

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadImage() (*widget.ButtonImage, error) {

	idle := image.NewBorderedNineSliceColor(color.NRGBA{R: 215, G: 130, B: 126, A: 255}, color.NRGBA{234, 157, 52, 255}, 3)

	hover := image.NewBorderedNineSliceColor(color.NRGBA{R: 215, G: 130, B: 126, A: 255}, color.NRGBA{234, 157, 52, 255}, 3)

	pressed := image.NewAdvancedNineSliceColor(color.NRGBA{R: 215, G: 130, B: 126, A: 255}, image.NewBorder(5, 5, 5, 5, color.NRGBA{246, 193, 119, 255}))

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

	user_home_directory, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Couldnt Find Your Home Directory: %v", err)
	}

	defaultPath := user_home_directory
	game := &Game{}

	Min, Current, Max := 0, 0, 10

	rootContainer := widget.NewContainer(

		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{

			0x20, 0x27, 0x3f, 0xff,
		})),

		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	ButtonImageInvisible, _ := loadImageInvisible()
	ButtonImage, _ := loadImage()
	Face, _ := loadFont(30)

	buttonGroup1 := widget.NewContainer(

		widget.ContainerOpts.Layout(widget.NewRowLayout(

			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(80),
		)),

		widget.ContainerOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionStart,
			}),
		),
	)

	buttonGroup2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Spacing(0))),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),
	)
	
	btn_Invisible_Two := widget.NewButton(

		widget.ButtonOpts.Image(ButtonImageInvisible),

		widget.ButtonOpts.Text("", Face, &widget.ButtonTextColor{

			Idle: color.NRGBA{42, 63, 39, 255},
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

	game.btn_CHTOFD = widget.NewButton(

		widget.ButtonOpts.Image(ButtonImage),

		widget.ButtonOpts.Text("Select Video", Face, &widget.ButtonTextColor{

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

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

			var err error
			selected_video_path, err = zenity.SelectFile(

				zenity.Filename(defaultPath),
				zenity.FileFilters{{

					Name:     "Video files",
					Patterns: []string{
						"*.mp4",
						"*.avi",
						"*.mkv",
						"*.mpeg",
						"*.mpg",
						"*.mov",
					},
					CaseFold: true,
				}},
			)

			if err != nil {

				log.Fatalf("Error Ocuured while Opening file dialogs : %v", err)
			}
			
			if FFmpegutils.Check_Duration(selected_video_path) == true {

				fmt.Print("That Your Video is correct Duration!!")
			}
			
			if FFmpegutils.Check_Duration(selected_video_path) == false {

				textContainer_for_btn_CHTOFD := widget.NewContainer(
					// The container will use a vertical row layout to lay out the progress
					// bars in a vertical row.
					widget.ContainerOpts.Layout(widget.NewRowLayout(
						widget.RowLayoutOpts.Direction(widget.DirectionVertical),
						widget.RowLayoutOpts.Spacing(180),
					)),
					// Set the required anchor layout data to determine where in the root
					// container to place the progress bars.
					widget.ContainerOpts.WidgetOpts(
						widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
							HorizontalPosition: widget.AnchorLayoutPositionCenter,
							VerticalPosition:   widget.AnchorLayoutPositionStart,
						}),
					),
				)

				text_for_video_too_long := widget.NewText(

					widget.TextOpts.Text("Your Video is Too Long!!", Face, rosePinePine),
					widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
				)

				textContainer_for_btn_CHTOFD.AddChild(text_for_video_too_long)
				rootContainer.RemoveChildren()
				rootContainer.AddChild(textContainer_for_btn_CHTOFD)
				
			}
			
			thumbnail_file_path := FFmpegutils.Extract_Thumbnail(selected_video_path)

			f, err := os.Open(thumbnail_file_path)
			if err != nil {

				log.Fatalf("Error Occured While Opening the image : %v",err)
			}
			defer f.Close()

			img, _, err := stdImage.Decode(f)
			if err != nil {

				log.Fatalf("Error While Decoding image  Dceode Image: %v")
			}

			game.thumbnail_image_image_format = img

			original_width, original_height := img.Bounds().Dx(), img.Bounds().Dy()
			if original_width == 0 || original_height == 0 {

				log.Println("Invalid thumbnail dimensions")
				return
			}

			const max_size = 400
			ratio := float64(original_width) / float64(original_height)
			var new_width, new_height int
			if float64(game.screen_width)/float64(game.screen_height) < ratio {
				new_width = max_size
				new_height =  int(float64(max_size) / ratio)
			} else {
				new_height = max_size
				new_width = int(float64(max_size) * ratio)
			}

			resizedImg := resize.Resize(uint(new_width), uint(new_height), game.thumbnail_image_image_format, resize.Lanczos3)
			game.resized_thumbnail_ebiten_image_format = ebiten.NewImageFromImage(resizedImg)

			buttonGroup1.RemoveChild(game.btn_CHTOFD)

		},
	))

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

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

			textContainer := widget.NewContainer(
				// The container will use a vertical row layout to lay out the progress
				// bars in a vertical row.
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(180),
				)),
				// Set the required anchor layout data to determine where in the root
				// container to place the progress bars.
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
						HorizontalPosition: widget.AnchorLayoutPositionCenter,
						VerticalPosition:   widget.AnchorLayoutPositionStart,
					}),
				),
			)

			progressBarsContainer := widget.NewContainer(
				// The container will use a vertical row layout to lay out the progress
				// bars in a vertical row.
				widget.ContainerOpts.Layout(widget.NewRowLayout(
					widget.RowLayoutOpts.Direction(widget.DirectionVertical),
					widget.RowLayoutOpts.Spacing(20),
				)),
				// Set the required anchor layout data to determine where in the root
				// container to place the progress bars.
				widget.ContainerOpts.WidgetOpts(
					widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
						HorizontalPosition: widget.AnchorLayoutPositionCenter,
						VerticalPosition:   widget.AnchorLayoutPositionCenter,
					}),
				),
			)

			text := widget.NewText(

				widget.TextOpts.Text("Loading.....", Face, rosePinePine),
				widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
			)

			progressBar := widget.NewProgressBar(

				widget.ProgressBarOpts.Direction(widget.DirectionHorizontal),

				widget.ProgressBarOpts.WidgetOpts(

					widget.WidgetOpts.MinSize(600, 40),
				),

				widget.ProgressBarOpts.Images(

					&widget.ProgressBarImage{
						Idle: image.NewNineSliceColor(color.NRGBA{40, 105, 131, 255}),
					},

					&widget.ProgressBarImage{
						Idle: image.NewNineSliceColor(color.NRGBA{235, 111, 146, 255}),
					},
				),

				widget.ProgressBarOpts.Values(Min, Max, Current),

				widget.ProgressBarOpts.TrackPadding(widget.Insets{
					Top:    2,
					Bottom: 2,
				}),
			)

			rootContainer.RemoveChildren()
			textContainer.AddChild(btn_Invisible_Two)
			textContainer.AddChild(text)
			rootContainer.AddChild(textContainer)
			progressBarsContainer.AddChild(progressBar)
			game.resized_thumbnail_ebiten_image_format.Deallocate()
			rootContainer.AddChild(progressBarsContainer)
			
			
			go func() {

				for Current < Max {

					Current++
					progressBar.SetCurrent(Current)
					time.Sleep(50 * time.Millisecond)
				}

				text_to_notify_video_is_converted := widget.NewText(

					widget.TextOpts.Text("Your ASCII Video Is Ready Check Your File Explorer!!", Face, rosePinePine),
					widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
				)

				rootContainer.RemoveChildren()
				textContainer.RemoveChild(text)
				textContainer.AddChild(text_to_notify_video_is_converted)
				rootContainer.AddChild(textContainer)

				videoPath := selected_video_path
				// if FFmpegutils.Check_Duration(videoPath) == false {

				// 	text_container_for_btn_TOASCII := widget.NewContainer(

				// 		widget.WidgetOpts.ContainerOpts(widget.AnchorLayout{

				// 			widget.ContainerOpts.AnchorLayoutData{

				// 				HorizontalPosition:
				// 			}
				// 		})
				// 	)
				// 	text_to_notify_video_is_too_long := widget.NewText(

				// 		widget.TextOpts.Text("Your Video Is Too Long!!", Face, rosePinePine),
				// 		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
				// 	)
				
				// }
				path_to_output_ascii_video_file := filepath.Join(user_home_directory, "output.mp4")
				outputPATH := path_to_output_ascii_video_file
				
				// mpp := Frame_Processing.Calculate_Ink_Required_For_Drawing_ASCII_Chars()
				// Frame_Processing.Generate_ASCII_Lookup_Table(mpp)

				// path_to_ASCII_Frames := filepath.Join(user_home_directory, "ASCII_Frames")
				// if err := os.MkdirAll(path_to_ASCII_Frames, 0750); err != nil {

				// 	log.Println("Error Occured While Creating ASCII Frames Directory: %v", err)
				// }

				// Count := 1
				
				result := FFmpegutils.ExtractFramesFromVideo(videoPath)

				directory, err := os.ReadDir(result)
				if err != nil {

					log.Fatalf("Error Occured While Reading Contents of The Frames Folder: %v", err)
				}

				var wg sync.WaitGroup

				for i, item := range directory {

					wg.Add(1)

					go func(count int, filename string) {

						defer wg.Done()

						framePath := filepath.Join(result, filename)

						Image, err := Frame_Processing.LoadAndResizeImage(framePath)

						if err != nil {

							log.Fatalf("ERROR OCCURED WHILE LOADING IMAGE: %v", err)
						}

						Pixels, rgbaValues := Frame_Processing.ExtractPixelData(Image)
						asciiImage, err := Frame_Processing.RenderAsciiImage(Pixels, rgbaValues)
						if err != nil {

							log.Fatalf("ERROR OCCURED WHILE CONVERTING TO ASCII: %v", err)

						}
						err = Frame_Processing.SaveImage(asciiImage, count)
						if err != nil {

							log.Fatalf("Error occured while saving the image: %v", err)
						}
						
						fmt.Printf("Frame: %d processed successfully\n", count-1)
					}(i + 1, item.Name())

				}

				wg.Wait()

				log.Println("ASCII frames folder generated successfully!")

				FFmpegutils.StitchFramesToVideo(outputPATH)

				fmt.Println("ASCII Video Generated successfully")
			}()

			
			// text_to_notify_video_is_converted := widget.NewText(

			// 	widget.TextOpts.Text("Your ASCII Video is ready check your File Explorer", Face, rosePinePine),
			// 	widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
			// )

			// textContainer.AddChild(text_to_notify_video_is_converted)
			// rootContainer.AddChild(textContainer)
			

			println("Buttons is Clicked")

		}),

		widget.ButtonOpts.WidgetOpts(

			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{

				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),
	)

	buttonGroup1.AddChild(btn_Invisible_Two)
	buttonGroup1.AddChild(game.btn_CHTOFD)
	rootContainer.AddChild(buttonGroup1)
	buttonGroup2.AddChild(btn_ToASCII)
	// buttonGroup2.AddChild(btn_Invisible)
	// buttonGroup2.AddChild(btn_FromASCII)
	// rootContainer.AddChild(buttonGroup1)
	rootContainer.AddChild(buttonGroup2)
	//	rootContainer.AddChild(thumbnail_container)

	ui := ebitenui.UI{

		Container: rootContainer,
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Clip2ASCII")

	// game := Game{
	// 	ui:  &ui,
	// 	btn: btn_ToASCII,
	// }
	game.ui = &ui
	game.btn = btn_ToASCII
	game.btn_CHTOFD = game.btn_CHTOFD
	// game.ButtonTraversal(ButtonImageInvisible)

	err = ebiten.RunGame(game)
	if err != nil {
		log.Println(err)
	}

}

func (g *Game) Update() error {

	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// screen.DrawImage(new_ebiten_image, nil)
	g.ui.Draw(screen)
	if g.resized_thumbnail_ebiten_image_format != nil {
		op := &ebiten.DrawImageOptions{}
		img_width, img_height := g.resized_thumbnail_ebiten_image_format.Size()
		op.GeoM.Translate(float64((g.screen_width - img_width)/2), float64((g.screen_height - img_height)/2))
		screen.DrawImage(g.resized_thumbnail_ebiten_image_format, op)
	}
	
}

func (g *Game) Layout(w, h int) (int, int) {
	g.screen_width = w
	g.screen_height = h
	return w, h
}
