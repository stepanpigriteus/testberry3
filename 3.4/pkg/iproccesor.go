package pkg

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"

	"io"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func ProcessImage(src image.Image, watermarkPath, format string) ([]io.Reader, error) {
	watermarkImg, err := imaging.Open(watermarkPath)
	if err != nil {
		return nil, fmt.Errorf("open watermark: %w", err)
	}

	thumb := imaging.Thumbnail(src, 150, 150, imaging.Lanczos)
	resized := imaging.Resize(src, 500, 0, imaging.Lanczos)

	bounds := src.Bounds()
	wb := watermarkImg.Bounds()
	scale := float64(bounds.Dx()) * 0.5 / float64(wb.Dx())
	watermarkScaled := imaging.Resize(watermarkImg, int(float64(wb.Dx())*scale), 0, imaging.Lanczos)

	offsetX := (bounds.Dx() - watermarkScaled.Bounds().Dx()) / 2
	offsetY := (bounds.Dy() - watermarkScaled.Bounds().Dy()) / 2

	watermarked := imaging.Clone(src)
	draw.Draw(
		watermarked,
		image.Rect(offsetX, offsetY, offsetX+watermarkScaled.Bounds().Dx(), offsetY+watermarkScaled.Bounds().Dy()),
		watermarkScaled,
		image.Point{},
		draw.Over,
	)

	images := []image.Image{thumb, resized, watermarked}
	var readers []io.Reader

	for _, img := range images {
		buf := new(bytes.Buffer)
		if err := encodeImage(buf, img, format); err != nil {
			return nil, err
		}
		readers = append(readers, buf)
	}

	return readers, nil
}

func encodeImage(w io.Writer, img image.Image, format string) error {
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(w, img)
	case "gif":
		return gif.Encode(w, img, nil)
	case "bmp":
		return bmp.Encode(w, img)
	case "tiff":
		return tiff.Encode(w, img, nil)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
