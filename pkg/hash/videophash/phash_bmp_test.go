package videophash

import (
	"bytes"
	"encoding/base64"
	"image"
	"testing"
)

func TestBMPDecoderRegistration(t *testing.T) {
	const fixture = "Qk1GAAAAAAAAADYAAAAoAAAAAgAAAAIAAAABABgAAAAAABAAAAATCwAAEwsAAAAAAAAAAAAAAAD/AP8AAAD/AAD///8AAA=="
	data, err := base64.StdEncoding.DecodeString(fixture)
	if err != nil {
		t.Fatal(err)
	}
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode BMP fixture: %v", err)
	}
	if format != "bmp" {
		t.Fatalf("format = %q, want bmp", format)
	}
	if got := img.Bounds().Size(); got.X != 2 || got.Y != 2 {
		t.Fatalf("size = %v, want 2x2", got)
	}
}
