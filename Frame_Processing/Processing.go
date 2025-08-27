package Frame_Processing

import (
	
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"path/filepath"
	"io"
	"log"
	"os"
	"math"

	"github.com/KononK/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

/*
Index 0: . (darkest)
Index 8: @ (brightest)
so the characterset is from darkest to brightest
*/

const ASCII = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,\\\"^`'. "
const newWidth = 300
const charWidth = 8
const charHeight = 12
var Path_to_ASCII_FRAMES_delete string

// This function loads a PNG image and resizes it to prepare for ASCII conversion,
// The resize calculation is crucial because it prevents distorted ASCII art.
func LoadAndResizeImage(framePath string) (image.Image, error) {

	File, err := os.Open(framePath)
	if err != nil {

		return nil, err
	}
	defer File.Close()

	img, err := png.Decode(File)
	if err != nil {

		return nil, err
	}

	bounds := img.Bounds()

	originalHeight := bounds.Dy()
	originalWidth := bounds.Dx()

	//ASCII characters are taller than they are wide. Without correction, your ASCII art would look stretched vertically.
	character_Pixel_Aspect_Ratio := float64(charWidth) / float64(charHeight)

	newHeight := uint((float64(originalHeight) / float64(originalWidth)) * float64(newWidth) * character_Pixel_Aspect_Ratio)
	// For very wide images (like panoramas),
	// the calculation might round down to 0. This ensures we always have at least 1 pixel height.
	if newHeight == 0 && originalHeight > 0 {

		newHeight = 1
	}

	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	out, err := os.Create("resized.png")
	if err != nil {

		log.Println("ERROR OCCURED WHILE CREATING THE FILE: ", err)
		return nil, err
	}
	defer out.Close()

	err = png.Encode(out, resizedImage)
	if err != nil {

		log.Println("ERROR OCCURED WHILE ENCODING THE IMAGE: ", err)
		return nil, err
	}

	return resizedImage, nil

}

func ExtractPixelData(img image.Image) ([][]uint8, [][]color.RGBA) {

	bounds := img.Bounds()
	Width := bounds.Dx()
	Height := bounds.Dy()

	Pixels := make([][]uint8, Height)
	rgbaValues := make([][]color.RGBA, Height)

	for y := 0; y < Height; y++ {

		row := make([]uint8, Width)
		rgbaRow := make([]color.RGBA, Width)
		for x := 0; x < Width; x++ {

			r, g, b, _ := img.At(x, y).RGBA()

			gammaCorrectionPixelByPixel := func(color_value uint8) float64{

				new_color_value := math.Pow((float64)(color_value)/ 255.0, 2.2)

				return new_color_value
			}

			/*
				Go's RGBA() returns 16-bit values (0-65535)
				But we need 8-bit values (0-255) for standard RGB
				Right shifting by 8 bits (>> 8) divides by 256
				Example: If r = 65535, then r >> 8 = 255
			*/

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			/*
				Human eyes perceive green as brightest (0.7152 = 71.52%)
				Red is moderately bright (0.2126 = 21.26%)
				Blue appears dimmest (0.0722 = 7.22%)
			*/
			rgbaRow[x] = color.RGBA{R: r8, G: g8, B: b8, A: 255}
			Brightness := uint8(255 * (0.2126*float64(gammaCorrectionPixelByPixel(r8)) + 0.7152*float64(gammaCorrectionPixelByPixel(g8)) + 0.0722*float64(gammaCorrectionPixelByPixel(b8))))
			row[x] = Brightness

		}
		rgbaValues[y] = rgbaRow
		Pixels[y] = row

	}

	return Pixels, rgbaValues
}

func Map_Brightness_To_Char(brightnessValue uint8) byte {

	var asciiValue byte
	//float64(brightnessValue) / 255.0,
	//This converts brightness (0-255) to a fraction (0.0-1.0),
	// Multiply by (len(ASCII) - 1) This scales the 0.0-1.0 range to 0.0-8.0.

	asciiIndex := int(float64((len(ASCII) - 1)) * float64(brightnessValue) / 255.0)

	asciiValue = ASCII[asciiIndex]

	return asciiValue
}

func RenderAsciiImage(Pixels [][]uint8, rgbaValues [][]color.RGBA) (*image.RGBA, error) {

	asciiMatrix := make([][]byte, len(Pixels))

	for rowNum, row := range Pixels {

		asciiRow := make([]byte, len(row))
		cellNumber := 0

		for _, cell := range row {

			asciiRow[cellNumber] = Map_Brightness_To_Char(cell)
			cellNumber++
		}
		asciiMatrix[rowNum] = asciiRow
	}

	//Opens a TTF font file from disk
	user_home_directory, err := os.UserHomeDir()
	if err != nil {

		log.Fatalf("Couldnt Find Your Home Directory: %v", err)
	}
	path_to_font_file := filepath.Join(user_home_directory, "Font.ttf")
	fontFile, err := os.Open(path_to_font_file)
	if err != nil {

		log.Println("ERROR OCCURED WHILE OPENING FILE: ", err)
		return nil, err
	}
	defer fontFile.Close()
	
	//Reads the entire font file into memory
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {

		log.Println("ERROR OCCURED WHILE READING FONTFILE: ", err)
	}

	//Parses the font data into a usable font object
	ttfFont, err := opentype.Parse(fontBytes)
	if err != nil {

		log.Println("ERROR OCCURED WHILE PARSING FONT: ", err)
		return nil, err
	}

	//Create Font face with specific settings:
	//Size: 16.0 - Font size in points
	//DPI: 150.0 - Dots per inch (higher = sharper text)
	//Hinting: None - No font hinting (pixel grid alignment)
	SourceCode, _ := opentype.NewFace(ttfFont, &opentype.FaceOptions{

		Size:    16.0,
		DPI:     150.0,
		Hinting: font.HintingNone,
	})

	//Your code creates a canvas like this:
	//Canvas width: 3 × 8 = 24 pixels
	//Canvas height: 3 × 12 = 36 pixels

	/*	Canvas Grid (24×36 pixels):
		0    8    16   24
		├────┼────┼────┤ 0
		│    │    │    │
		│    │    │    │ 12
		├────┼────┼────┤
		│    │    │    │
		│    │    │    │ 24
		├────┼────┼────┤
		│    │    │    │
		│    │    │    │ 36
		└────┴────┴────┘
	*/

	//Creates a blank RGBA image
	//Width: number_of_columns × charWidth (e.g., 300 × 8 = 2400 pixels)
	//Height: number_of_rows × charHeight (e.g., 133 × 12 = 1596 pixels)
	//Fills the entire canvas with white background

	Img := image.NewRGBA(image.Rect(0, 0, len(Pixels[0])*charWidth, len(Pixels)*charHeight))
	draw.Draw(Img, Img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	Drawer := &font.Drawer{

		Dst:  Img,
		Src:  nil,
		Face: SourceCode,
	}
	//The Below Loop Follows the following steps :

	/*
		Step-by-step for our 3×3:

		y=0, x=0: Draw '.' at (0, 12) with color from rgbaValues[0][0]
		y=0, x=1: Draw '-' at (8, 12) with color from rgbaValues[0][1]
		y=0, x=2: Draw ':' at (16, 12) with color from rgbaValues[0][2]
		y=1, x=0: Draw '=' at (0, 24) with color from rgbaValues[1][0]
		y=1, x=1: Draw '+' at (8, 24) with color from rgbaValues[1][1]
		y=1, x=2: Draw '*' at (16, 24) with color from rgbaValues[1][2]
		y=2, x=0: Draw '#' at (0, 36) with color from rgbaValues[2][0]
		y=2, x=1: Draw '%' at (8, 36) with color from rgbaValues[2][1]
		y=2, x=2: Draw '@' at (16, 36) with color from rgbaValues[2][2]
	*/

	for y, row := range asciiMatrix {

		for x, ch := range row {

			Drawer.Src = image.NewUniform(rgbaValues[y][x])
			//Here x simply gives us the Horizontal Position
			// while we do Y + 1 to get the vertical baseline(needed for some kind of font mathematics).
			// Refer to this SO link to know more: https://stackoverflow.com/questions/27631736/meaning-of-top-ascent-baseline-descent-bottom-and-leading-in-androids-font
			Drawer.Dot = fixed.P(x*charWidth, (y+1)*charHeight)

			//Converts the ASCII byte to a string
			//Draws it at the calculated position with the set color
			Drawer.DrawString(string(ch))

		}
	}

	//Final Image will be something like :

	/*
		Final 24×36 pixel image:
		0    8    16   24
		├────┼────┼────┤ 0
		│ .  │ -  │ :  │ ← Row 0 characters
		│    │    │    │ 12
		├────┼────┼────┤
		│ =  │ +  │ *  │ ← Row 1 characters
		│    │    │    │ 24
		├────┼────┼────┤
		│ #  │ %  │ @  │ ← Row 2 characters
		│    │    │    │ 36
		└────┴────┴────┘
	*/

	return Img, nil
}

func SaveImage(Img *image.RGBA,Count int) (error error) {

	user_home_directory, err := os.UserHomeDir()
	if err != nil {

		log.Fatalf("Couldnt Find Your Home Directory: %v", err)
	}

	path_to_ASCII_Frames := filepath.Join(user_home_directory, "ASCII_Frames")
	if err := os.MkdirAll(path_to_ASCII_Frames, 0750); err != nil {

		log.Println("Error Occured While Creating ASCII Frames Directory: %v", err)
	}
	s := filepath.Join(path_to_ASCII_Frames, fmt.Sprintf("ASCII_Frames%03d.png", Count))
	newImage,err := os.Create(s)

	if err != nil {

		log.Println("Couldnt create File the following error occured: ", err)
		return err
		
	}

	defer func() {
		if cerr := newImage.Close(); cerr != nil {
			if err == nil {
				err = fmt.Errorf("error closing file %s: %w", s, cerr)
			} else {
				log.Printf("warning: error closing file %s: %v", s, cerr)
			}
		}
	}()
	
	if err = png.Encode(newImage, Img); err != nil {

		log.Println("Couldnt convert to grayScale: ", err)
		return err
		
	}
	
	Path_to_ASCII_FRAMES_delete = path_to_ASCII_Frames
	
	return nil

}
