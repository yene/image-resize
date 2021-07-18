package main

import (
	"image"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
)

func main() {
	// config
	autorotate := true

	// a5, dpi 300
	targetHeight := 1748
	targetWidth := 2480
	targetRatio := float64(targetWidth) / float64(targetHeight)

	src, err := imaging.Open("tests/input.jpg")
	srcWidth := src.Bounds().Max.X
	srcHeight := src.Bounds().Max.Y
	srcRatio := float64(srcWidth) / float64(srcHeight)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	dst := imaging.New(targetWidth, targetHeight, color.White)

	if autorotate && srcHeight > srcWidth {
		src = imaging.Rotate90(src)
		srcWidth = src.Bounds().Max.X
		srcHeight = src.Bounds().Max.Y
		srcRatio = float64(srcWidth) / float64(srcHeight)
	}

	if targetRatio < srcRatio {
		img1 := imaging.Resize(src, targetWidth, 0, imaging.Lanczos)
		img1Height := img1.Bounds().Max.Y
		offsetTop := (targetHeight - img1Height) / 2
		dst = imaging.Paste(dst, img1, image.Pt(0, offsetTop))
	} else {
		img1 := imaging.Resize(src, 0, targetHeight, imaging.Lanczos)
		img1Width := img1.Bounds().Max.X
		offsetLeft := (targetWidth - img1Width) / 2
		dst = imaging.Paste(dst, img1, image.Pt(offsetLeft, 0))
	}

	// Save the resulting image as JPEG.
	err = imaging.Save(dst, "tests/out_example.jpg")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
}
