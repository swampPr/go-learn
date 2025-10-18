// Package main provides main  INFO:
package main

import (
	"fmt"
	"image/png"
	"os"
)

func main() {
	file, _ := os.Open("./img_test.png")

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	img, err := png.Decode(file)
	bounds := img.Bounds()
	if err != nil {
		panic(err)
	}
	rgb := make([]byte, 0, bounds.Dx()*bounds.Dy()*3)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgb = append(rgb, byte(r), byte(g), byte(b))
		}
	}

	fmt.Println(string(rgb))
}
