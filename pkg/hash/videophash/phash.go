package videophash

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/bmp"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

const (
	screenshotSize = 160
	columns        = 5
	rows           = 5
)

func Generate(encoder *ffmpeg.FFMpeg, videoFile *models.VideoFile) (*uint64, error) {
	sprite, err := generateSprite(encoder, videoFile)
	if err != nil {
		return nil, err
	}

	hash, err := goimagehash.PerceptionHash(sprite)
	if err != nil {
		return nil, fmt.Errorf("computing phash from sprite: %w", err)
	}
	hashValue := hash.GetHash()
	return &hashValue, nil
}

func generateSpriteScreenshot(encoder *ffmpeg.FFMpeg, input string, t float64) (image.Image, error) {
	options := transcoder.ScreenshotOptions{
		Width:      screenshotSize,
		OutputPath: "-",
		OutputType: transcoder.ScreenshotOutputTypePNG,
	}

	args := transcoder.ScreenshotTime(input, t, options)
	data, err := encoder.GenerateOutput(context.Background(), args, nil)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	img, err := png.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	return img, nil
}

func combineImages(images []image.Image) image.Image {
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y
	canvasWidth := width * columns
	canvasHeight := height * rows
	montage := imaging.New(canvasWidth, canvasHeight, color.NRGBA{})
	for index := 0; index < len(images); index++ {
		x := width * (index % columns)
		y := height * int(math.Floor(float64(index)/float64(rows)))
		img := images[index]
		montage = imaging.Paste(montage, img, image.Pt(x, y))
	}

	return montage
}

func generateSprite(encoder *ffmpeg.FFMpeg, videoFile *models.VideoFile) (image.Image, error) {
	logger.Infof("[generator] generating phash sprite for %s", videoFile.Path)
	return generateSpriteWithScreenshot(videoFile.Duration, func(t float64) (image.Image, error) {
		return generateSpriteScreenshot(encoder, videoFile.Path, t)
	})
}

func generateSpriteWithScreenshot(duration float64, screenshot func(float64) (image.Image, error)) (image.Image, error) {
	// Generate sprite image offset by 5% on each end to avoid intro/outros.
	chunkCount := columns * rows
	offset := 0.05 * duration
	stepSize := (0.9 * duration) / float64(chunkCount)
	images := make([]image.Image, 0, chunkCount)
	for i := 0; i < chunkCount; i++ {
		time := offset + (float64(i) * stepSize)
		img, err := screenshot(time)
		if err != nil {
			return nil, fmt.Errorf("generating sprite screenshot %d/%d: %w", i+1, chunkCount, err)
		}
		images = append(images, img)
	}

	return combineImages(images), nil
}
