package main

import (
	"image/png"
	"os"
)

type texture struct {
	pixels      []byte
	w, h, pitch int
}

func (tex *texture) draw(position pos, screenPixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(position.y) - tex.h/2
			screenX := x + int(position.x) - tex.w/2
			if screenX >= 0 && screenX < windowWidth && screenY >= 0 && screenY < windowHeight {
				// Index to read the data from the texture
				texIndex := y*tex.pitch + x*4
				// Index to copy the data to the screen
				screenIndex := screenY*windowWidth*4 + screenX*4

				screenPixels[screenIndex] = tex.pixels[texIndex]
				screenPixels[screenIndex+1] = tex.pixels[texIndex+1]
				screenPixels[screenIndex+2] = tex.pixels[texIndex+2]
				screenPixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func (tex *texture) drawAlpha(position pos, pixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(position.y) - tex.h/2
			screenX := x + int(position.x) - tex.w/2
			if screenX >= 0 && screenX < windowWidth && screenY >= 0 && screenY < windowHeight {
				// Index to read the data from the texture
				texIndex := y*tex.pitch + x*4
				// Index to copy the data to the screen
				screenIndex := screenY*windowWidth*4 + screenX*4

				srcR := int(tex.pixels[texIndex])
				srcG := int(tex.pixels[texIndex+1])
				srcB := int(tex.pixels[texIndex+2])
				srcA := int(tex.pixels[texIndex+3])

				dstR := int(pixels[screenIndex])
				dstG := int(pixels[screenIndex+1])
				dstB := int(pixels[screenIndex+2])

				rstR := (srcR*255 + dstR*(255-srcA)) / 255
				rstG := (srcG*255 + dstG*(255-srcA)) / 255
				rstB := (srcB*255 + dstB*(255-srcA)) / 255

				pixels[screenIndex] = byte(rstR)
				pixels[screenIndex+1] = byte(rstG)
				pixels[screenIndex+2] = byte(rstB)
				// pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func loadTexture(filepath string, invertedHorizontally bool) texture {
	infile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	texturePixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			tempX := x
			if invertedHorizontally {
				tempX = w - x
			}
			r, g, b, a := img.At(tempX, y).RGBA()
			texturePixels[bIndex] = byte(r / 256)
			bIndex++
			texturePixels[bIndex] = byte(g / 256)
			bIndex++
			texturePixels[bIndex] = byte(b / 256)
			bIndex++
			texturePixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	return texture{texturePixels, w, h, w * 4}
}
