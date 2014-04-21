// This example demonstrates decoding a JPEG image and examining its pixels.
package main

import (
	"github.com/cartland/go/imagic"
	"image"
	"log"
	"os"
	// _ "image/gif"
	_ "image/jpeg"
	"image/png"
)

func main() {
	// Decode the JPEG data.
	reader, err := os.Open("testdata/Chefchaouen.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	bg, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	reader, err = os.Open("testdata/borrodepth.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	dm, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	configWallEyed := imagic.Config{60, 100, false}
	configCrossEyed := imagic.Config{100, 160, true}

	wall := imagic.Imagic(dm, bg, configWallEyed)
	writer, err := os.Create("testdata/wallOutput.png")
	png.Encode(writer, wall)

	cross := imagic.Imagic(dm, bg, configCrossEyed)
	writer, err = os.Create("testdata/crossOutput.png")
	png.Encode(writer, cross)
}
