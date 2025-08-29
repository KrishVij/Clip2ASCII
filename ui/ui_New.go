package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
	"sync"
	"strings"
	"path/filepath"

	stdImage "image"
	"image/jpeg"
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
	done                                  bool
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

				return
				
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

	defer FFmpegutils.Delete_Thumbnail_Folder()

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

				widget.TextOpts.Text("Loading...", Face, rosePinePine),
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
			if game.resized_thumbnail_ebiten_image_format != nil {

				game.resized_thumbnail_ebiten_image_format.Deallocate()
			}
			rootContainer.AddChild(progressBarsContainer)

			go func() {

				for Current < Max {

					Current++
					progressBar.SetCurrent(Current)
					time.Sleep(50 * time.Millisecond)
				}

				text_to_notify_video_is_converted := widget.NewText(

					widget.TextOpts.Text("Starting Conversion..Please Wait Patiently !!", Face, rosePinePine),
					widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
				)

				rootContainer.RemoveChildren()
				textContainer.RemoveChild(text)
				textContainer.AddChild(text_to_notify_video_is_converted)
				rootContainer.AddChild(textContainer)

				videoPath := selected_video_path
				slice_of_splitted_video_path := []string{}
				slice_of_splitted_video_path = strings.Split(selected_video_path, "\\")
				path_to_output_ascii_video_file := filepath.Join(user_home_directory, slice_of_splitted_video_path[len(slice_of_splitted_video_path) - 1])
				outputPATH := path_to_output_ascii_video_file

				fmt.Println("Extracting frames from video...")
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
						str := fmt.Sprintf("Frame: %d processed successfully\n", count-1)
						text_Print_Progress := widget.NewText(
							widget.TextOpts.Text(str, Face, rosePinePine),
							widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter))

						rootContainer.RemoveChildren()
						textContainer.RemoveChildren()
						textContainer.AddChild(text_Print_Progress)
						rootContainer.AddChild(textContainer)
					}(i + 1, item.Name())

				}

				wg.Wait()

				creating_final_video_string := "Creating final video..."
				text_Print_final_video_string := widget.NewText(
					widget.TextOpts.Text(creating_final_video_string, Face, rosePinePine),
					widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter))

				rootContainer.RemoveChildren()
				textContainer.RemoveChildren()
				textContainer.AddChild(text_Print_final_video_string)
				rootContainer.AddChild(textContainer)
				
				FFmpegutils.StitchFramesToVideo(outputPATH)

				creating_generated_successfully_string := "ASCII Video Generated successfully!! Close The App !!"

				text_Print_generated_successfully_string := widget.NewText(
					widget.TextOpts.Text(creating_generated_successfully_string, Face, rosePinePine),
					widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter))

				rootContainer.RemoveChildren()
				textContainer.RemoveChildren()
				textContainer.AddChild(text_Print_generated_successfully_string)
				rootContainer.AddChild(textContainer)

				defer FFmpegutils.Delete_Generated_Folders()
				
			}()
				
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
	rootContainer.AddChild(buttonGroup2)
	
	ui := ebitenui.UI{

		Container: rootContainer,
	}

	path_to_image_icon := filepath.Join(defaultPath, "Clip2ASCII", "ui", "imageIcon.jpg")
	fmt.Println(path_to_image_icon)
	icon_image, err := os.Open(path_to_image_icon)
	if err != nil {

		log.Fatalf("Error Occured While Opening Icon Image: %v", err)
	}

	jpeg_decoded_icon_image, err := jpeg.Decode(icon_image)
	if err != nil {

		log.Fatalf("Error Occured While Decoding The Icon Image : %v", err)
	}
	defer icon_image.Close()
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowIcon([]stdImage.Image{jpeg_decoded_icon_image})
	ebiten.SetWindowTitle("Clip2ASCII")
	
	game.ui = &ui
	game.btn = btn_ToASCII
	game.btn_CHTOFD = game.btn_CHTOFD

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
