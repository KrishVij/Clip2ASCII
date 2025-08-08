package Frame_Processing

import (
	img "image"
	"image/color"
	"image/draw"
	"io"
	"log"
	"os"
	"slices"

	// "github.com/hexops/gotextdiff/difftest"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Char_Metrics struct {
	ink        float64
	brightness float64
}

var lookup_table [256]string

func Calculate_Ink_Required_For_Drawing_ASCII_Chars() map[string]Char_Metrics {

	fontFile, err := os.Open("Font.ttf")
	if err != nil {

		log.Println("ERROR OCCURED WHILE OPENING FILE: ", err)
	}

	//Reads the entire font file into memory
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {

		log.Println("ERROR OCCURED WHILE READING FONTFILE: ", err)
	}

	//Parses the font data into a usable font object
	ttfFont, err := opentype.Parse(fontBytes)
	if err != nil {

		log.Println("ERROR OCCURED WHILE PARSING FONT: ", err)
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

	metrics_for_each_char_from_ASCII_string := make(map[string]Char_Metrics)

	for i := 0; i < len(ASCII); i++ {

		newFontImage := img.NewRGBA(img.Rectangle{

			Min: img.Point{X: 0, Y: 0},
			Max: img.Point{X: 32, Y: 32},
		})

		draw.Draw(newFontImage, newFontImage.Bounds(), &img.Uniform{color.White}, img.Point{}, draw.Src)

		Drawer := &font.Drawer{

			Dst:  newFontImage,
			Src:  &img.Uniform{color.Black},
			Face: SourceCode,
			Dot:  fixed.P(0, 24),
		}

		Drawer.DrawString(string(ASCII[i]))

		black_pixels_of_drawn_font := 0
		total_no_of_pixels := 0
		total_brightness := 0.0

		for y := 0; y < newFontImage.Bounds().Dy(); y++ {

			for x := 0; x < newFontImage.Bounds().Dx(); x++ {

				r, g, b, _ := newFontImage.At(x, y).RGBA()

				r8 := uint8(r >> 8)
				g8 := uint8(g >> 8)
				b8 := uint8(b >> 8)

				if r != 0xFFFF || g != 0xFFFF || b != 0xFFFF {

					black_pixels_of_drawn_font++
				}

				brightness := 0.2126*(float64)(r8) + 0.7152*(float64)(g8) + 0.0722*(float64)(b8)
				total_brightness += brightness / 255.0
				total_no_of_pixels++

			}
		}

		total_pixels := float64(newFontImage.Bounds().Dy() * newFontImage.Bounds().Dx())
		ink := (float64)(black_pixels_of_drawn_font) / (float64)(total_pixels)
		average_brightness := total_brightness/float64(total_no_of_pixels)

		metrics_for_each_char_from_ASCII_string[(string)(ASCII[i])] = Char_Metrics{ink: ink, brightness: average_brightness}

	}

	return metrics_for_each_char_from_ASCII_string
}

func Generate_ASCII_Lookup_Table(ink_for_each_char_from_ASCII_string map[string]Char_Metrics) [256]string {

	type Score_Of_Character struct {

		char string
		score float64
	}

	sorted_chars := make([]Score_Of_Character, 0, 256)

	for char, metrics_required := range ink_for_each_char_from_ASCII_string {

		score := 0.7*metrics_required.ink + (1 - 0.7)*(1 - metrics_required.brightness)

		sorted_chars = append(sorted_chars, Score_Of_Character{

			char: char,
			score: score,
		})
	}

	slices.SortFunc(sorted_chars, func(a, b Score_Of_Character) int {

		if a.score < b.score {

			return -1

		} else if a.score > b.score {

			return 1
		}

		return 0
	})

	for i := 0; i < 256; i++ {

		gray_value_to_ink := (float64)(i) / 255.0

		index := int(gray_value_to_ink * float64(len(sorted_chars)-1))
		lookup_table[i] = sorted_chars[index].char

		// log.Printf("index : %d and char = %s", i, lookup_table[i])
	}

	return lookup_table

}
