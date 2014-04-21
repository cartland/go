/*
 * Copyright 2014 Chris Cartland
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
