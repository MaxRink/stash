package videophash

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/corona10/goimagehash"
	"golang.org/x/image/bmp"
)

func TestScreenshotDecodersProduceEquivalentPerceptualHash(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			src.SetRGBA(x, y, color.RGBA{R: uint8(x * 7), G: uint8(y * 7), B: uint8((x + y) * 3), A: 255})
		}
	}
	var bmpBytes, pngBytes bytes.Buffer
	if err := bmp.Encode(&bmpBytes, src); err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(&pngBytes, src); err != nil {
		t.Fatal(err)
	}
	bmpImage, err := bmp.Decode(bytes.NewReader(bmpBytes.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	pngImage, err := png.Decode(bytes.NewReader(pngBytes.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	bmpHash, err := goimagehash.PerceptionHash(bmpImage)
	if err != nil {
		t.Fatal(err)
	}
	pngHash, err := goimagehash.PerceptionHash(pngImage)
	if err != nil {
		t.Fatal(err)
	}
	if bmpHash.GetHash() != pngHash.GetHash() {
		t.Fatalf("BMP and PNG phash differ: %x != %x", bmpHash.GetHash(), pngHash.GetHash())
	}
}
