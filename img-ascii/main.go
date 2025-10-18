package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"os/user"
	"path"
)

const asciiChars = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/|()1{}[]?-_+~<>i!lI;:,\"^'\\`. "

type PixMap []byte

func grayToASCII(gray byte) string {
	return string(asciiChars[int(math.Round(float64(((len(asciiChars)-1)*int(gray))/255)))])
}

func colorToGray(rgb []byte) (gray []byte) {
	gray = make([]byte, 0, len(rgb)/3)

	for i := 0; i < len(rgb); i = i + 3 {
		r := float64(rgb[i])
		g := float64(rgb[i+1])
		b := float64(rgb[i+2])

		grayScale := r*0.2126 + g*0.7152 + b*0.0722

		gray = append(gray, byte(math.Round(grayScale)))
	}

	return gray
}

func (pm PixMap) reduce(fn func(ac string, val int, idx int) string, initVal string) string {
	ac := initVal
	for i, v := range pm {
		ac = fn(ac, int(v), i)
	}

	return ac
}

func getRGB(img image.Image) (rgb []byte) {
	bounds := img.Bounds()

	rgb = make([]byte, 0, bounds.Dx()*bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := byte(r >> 8)
			g8 := byte(g >> 8)
			b8 := byte(b >> 8)

			rgb = append(rgb, byte(r8), byte(g8), byte(b8))
		}
	}

	return rgb
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Please input a valid path")
	}

	filePath := args[0]

	userInfo, err := user.Current()
	if err != nil {
		panic("Could not get user's info")
	}

	absFilePath := path.Join(userInfo.HomeDir, filePath)

	//#nosec
	imgFile, err := os.Open(absFilePath)
	if err != nil {
		panic("Could not open image")
	}

	defer func() {
		if err := imgFile.Close(); err != nil {
			os.Exit(1)
		}
	}()

	img, format, err := image.Decode(imgFile)
	if err != nil {
		panic("Could not decode image")
	}

	if format != "png" && format != "jpeg" {
		log.Fatalf("Expected format jpeg or png, instead received format: %v", format)
	}

	rgb := getRGB(img)
	grayScale := PixMap(colorToGray(rgb))
	imgWidth := img.Bounds().Dx()

	ascii := grayScale.reduce(func(ac string, val int, idx int) string {
		nextChars := grayToASCII(byte(val))

		if idx+1%imgWidth == 0 {
			nextChars += "\n"
		}

		return ac + nextChars
	}, "")
	fmt.Println(ascii)
}
