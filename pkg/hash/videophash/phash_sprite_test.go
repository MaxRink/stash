package videophash

import (
	"errors"
	"image"
	"testing"
)

func TestGenerateSpriteRequiresAllCanonicalFrames(t *testing.T) {
	calls := 0
	sprite, err := generateSpriteWithScreenshot("test.mp4", 100, func(float64, bool) (image.Image, error) {
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
	_, err := generateSpriteWithScreenshot("test.mp4", 100, func(float64, bool) (image.Image, error) {
		calls++
		if calls >= columns*rows {
			return nil, errors.New("synthetic decode failure")
		}
		return image.NewRGBA(image.Rect(0, 0, 8, 8)), nil
	})
	if err == nil {
		t.Fatal("expected strict full-frame generation to fail")
	}
	// the failing last frame is retried once with accurate seek
	if calls != columns*rows+1 {
		t.Fatalf("screenshot calls = %d, want %d", calls, columns*rows+1)
	}
}

func TestGenerateSpriteRetriesWithAccurateSeek(t *testing.T) {
	var seeks []bool
	_, err := generateSpriteWithScreenshot("test.mp4", 100, func(_ float64, slowSeek bool) (image.Image, error) {
		seeks = append(seeks, slowSeek)
		if len(seeks) == 2 && !slowSeek {
			return nil, errors.New("synthetic fast seek failure")
		}
		return image.NewRGBA(image.Rect(0, 0, 8, 8)), nil
	})
	if err != nil {
		t.Fatalf("accurate seek retry should succeed: %v", err)
	}
	if len(seeks) != columns*rows+1 || seeks[0] || !seeks[2] || !seeks[len(seeks)-1] {
		t.Fatalf("unexpected seek sequence %v", seeks)
	}
}
