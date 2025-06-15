package Frame_Processing

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/KononK/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const ASCII = ".:-=+*#%@"
const newWidth = 300
const charWidth = 8
const charHeight = 12

func ImageResizeAndLoad() (image.Image, error) {

	File, err := os.Open("test.png")
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

	character_Pixel_Aspect_Ratio := float64(charWidth) / float64(charHeight)

	newHeight := uint((float64(originalHeight) / float64(originalWidth)) * float64(newWidth) * character_Pixel_Aspect_Ratio)
	if newHeight == 0 && originalHeight > 0 {

		newHeight = 1
	}

	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	out,err := os.Create("resized.png")
	if err != nil {

		log.Println("ERROR OCCURED WHILE CREATING THE FILE: ",err)
		return nil,err
	}
	defer out.Close()

	err = png.Encode(out,resizedImage)
	if err != nil {

		log.Println("ERROR OCCURED WHILE ENCODING THE IMAGE: ",err)
		return nil,err
	}

	return resizedImage, nil

}

func ProcessImageForAscii(img image.Image) ([][]uint8, [][]color.RGBA) {

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

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			rgbaRow[x] = color.RGBA{R: r8, G: g8, B: b8, A: 255}
			Brightness := uint8(0.2126*float64(r8) + 0.7152*float64(g8) + 0.0722*float64(b8))
			row[x] = Brightness

		}
		rgbaValues[y] = rgbaRow
		Pixels[y] = row

	}

	return Pixels, rgbaValues
}

func BrightnessToASCII(brightnessValue uint8) byte {

	var asciiValue byte

	asciiIndex := int(float64((len(ASCII) - 1)) * float64(brightnessValue) / 255.0)

	asciiValue = ASCII[asciiIndex]

	return asciiValue
}

func GrayScaleToAscii(Pixels [][]uint8, rgbaValues [][]color.RGBA) (*image.RGBA, error) {

	asciiMatrix := make([][]byte, len(Pixels))

	for rowNum, row := range Pixels {

		asciiRow := make([]byte, len(row))
		cellNumber := 0

		for _, cell := range row {

			asciiRow[cellNumber] = BrightnessToASCII(cell)
			cellNumber++
		}
		asciiMatrix[rowNum] = asciiRow
	}

	fontFile, err := os.Open("Font.ttf")
	if err != nil {

		log.Println("ERROR OCCURED WHILE OPENING FILE: ", err)
		return nil, err
	}
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {

		log.Println("ERROR OCCURED WHILE READING FONTFILE: ", err)
	}

	ttfFont, err := opentype.Parse(fontBytes)

	if err != nil {

		log.Println("ERROR OCCURED WHILE PARSING FONT: ", err)
		return nil, err
	}

	SourceCode, _ := opentype.NewFace(ttfFont, &opentype.FaceOptions{

		Size:    16.0,
		DPI:     150.0,
		Hinting: font.HintingNone,
	})

	Img := image.NewRGBA(image.Rect(0, 0, len(Pixels[0])*charWidth, len(Pixels)*charHeight))
	draw.Draw(Img, Img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	Drawer := &font.Drawer{

		Dst:  Img,
		Src:  nil,
		Face: SourceCode,
	}

	for y, row := range asciiMatrix {

		for x, ch := range row {

			Drawer.Src = image.NewUniform(rgbaValues[y][x])

			Drawer.Dot = fixed.P(x*charWidth, (y+1)*charHeight)
			Drawer.DrawString(string(ch))
		
		}
	}

	return Img, nil
}

func SaveImage(Img *image.RGBA) error {

	newImage, err := os.Create("output.png")

	if err != nil {

		log.Println("Couldnt create File the following error occured: ", err)
		return err
	}
	defer newImage.Close()

	err = png.Encode(newImage, Img)

	if err != nil {

		log.Println("Couldnt convert to grayScale: ", err)
		return err
	}

	return nil

}
