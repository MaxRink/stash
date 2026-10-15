package videophash

import (
	"errors"
	"image"
	"testing"
)

func TestGenerateSpriteRequiresAllCanonicalFrames(t *testing.T) {
	calls := 0
	sprite, err := generateSpriteWithScreenshot(100, func(float64) (image.Image, error) {
		calls++
		return image.NewRGBA(image.Rect(0, 0, 8, 8)), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != columns*rows {
		t.Fatalf("screenshot calls = %d, want %d", calls, columns*rows)
	}
	if got := sprite.Bounds().Size(); got.X != 8*columns || got.Y != 8*rows {
		t.Fatalf("sprite size = %v, want %dx%d", got, 8*columns, 8*rows)
	}
}

func TestGenerateSpriteFailsOnShortSprite(t *testing.T) {
	calls := 0
	_, err := generateSpriteWithScreenshot(100, func(float64) (image.Image, error) {
		calls++
		if calls == columns*rows {
			return nil, errors.New("synthetic decode failure")
		}
		return image.NewRGBA(image.Rect(0, 0, 8, 8)), nil
	})
	if err == nil {
		t.Fatal("expected strict full-frame generation to fail")
	}
	if calls != columns*rows {
		t.Fatalf("screenshot calls = %d, want %d", calls, columns*rows)
	}
}
