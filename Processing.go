package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const ASCII = "`^\\\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

func loadImage() *os.File {

	Image, err := os.Open("test.jpg")

	if err != nil {

		log.Println("Error Occured while Opening the image: ", err)

	}

	return Image

}

func getBrightnessMatrix(Image *os.File) [][]uint8 {

	imageInfo, err := jpeg.Decode(Image)

	if err != nil {

		log.Println("Couldn't get Image info: ", err)

	}

	bounds := imageInfo.Bounds()

	Pixels := make([][]uint8, 0)
	Width := bounds.Dx()
	Height := bounds.Dy()

	for y := 0; y < Height; y++ {

		row := make([]uint8, 0)
		for x := 0; x < Width; x++ {

			r, g, b, _ := imageInfo.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			Brightness := uint8(0.2126*float64(r8) + 0.7152*float64(g8) + 0.0722*float64(b8))
			row = append(row, Brightness)

		}

		if len(row) > 0 {

			Pixels = append(Pixels, row)
		}
	}

	return Pixels
}

func brightnessToASCII(brightnessValue uint8) byte {

	var asciiValue byte
	var asciiIndex int

	asciiIndex = int(float64((len(ASCII) - 1)) * float64(brightnessValue) / 255.0)

	asciiValue = ASCII[asciiIndex]

	return asciiValue
}

func grayScaleToAscii(Pixels [][]uint8) *image.RGBA {

	asciiMatrix := make([][]byte, len(Pixels))

	for rowNum, row := range Pixels {

		asciiRow := make([]byte, len(row))
		cellNumber := 0

		for _, cell := range row {

			asciiRow[cellNumber] = brightnessToASCII(cell)
			cellNumber++
		}
		asciiMatrix[rowNum] = asciiRow
	}

	fontFile, err := os.Open("Font.ttf")
	if err != nil {

		log.Println("ERROR OCCURED WHILE OPENING FILE: ", err)
	}
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {

		log.Println("ERROR OCCURED WHILE READING FONTFILE: ", err)
	}

	ttfFont, err := opentype.Parse(fontBytes)

	if err != nil {

		log.Println("ERROR OCCURED WHILE PARSING FONT: ", err)
	}

	SourceCode, _ := opentype.NewFace(ttfFont, &opentype.FaceOptions{

		Size:    12.0,
		DPI:     72.0,
		Hinting: font.HintingFull,
	})

	charWidth := 8
	charHeight := 12

	Img := image.NewRGBA(image.Rect(0, 0, len(Pixels)*charWidth, len(Pixels[0])*charHeight))
	draw.Draw(Img, Img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	Drawer := &font.Drawer{

		Dst:  Img,
		Src:  image.NewUniform(color.Black),
		Face: SourceCode,
	}

	for y, row := range asciiMatrix {

		for x, ch := range row {

			Drawer.Dot = fixed.P(x*charWidth, (y+1)*charHeight)
			Drawer.DrawString(string(ch))
		}
	}

	return Img
}

func saveImage(Img *image.RGBA) {

	newImage, err := os.Create("output.jpg")

	if err != nil {

		log.Println("Couldnt create File the following error occured: ", err)
	}

	err = jpeg.Encode(newImage, Img, nil)

	if err != nil {

		log.Println("Couldnt convert to grayScale: ", err)
	}
}

func main() {

	Image := loadImage()

	if Image != nil {

		defer Image.Close()
	}

	Pixels := getBrightnessMatrix(Image)
	asciiImage := grayScaleToAscii(Pixels)
	saveImage(asciiImage)

}
