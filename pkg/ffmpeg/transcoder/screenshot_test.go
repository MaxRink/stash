package transcoder

import (
	"reflect"
	"testing"
)

func TestScreenshotTimeBoundsDecoderAndFilterThreads(t *testing.T) {
	got := ScreenshotTime("input.webm", 12.5, ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
	})
	want := []string{
		"-v", "error", "-y",
		"-threads", "1", "-filter_threads", "1", "-filter_complex_threads", "1",
		"-ss", "12.5", "-i", "input.webm", "-frames:v", "1", "-f", "image2", "out.jpg",
	}
	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotFrameBoundsDecoderAndFilterThreads(t *testing.T) {
	got := ScreenshotFrame("input.webm", 12, ScreenshotOptions{
		OutputPath: "out.bmp",
		OutputType: ScreenshotOutputTypeBMP,
	})
	want := []string{
		"-v", "error", "-y",
		"-threads", "1", "-filter_threads", "1", "-filter_complex_threads", "1",
		"-i", "input.webm", "-frames:v", "1", "-vsync", "0", "-vf", "select=eq(n\\,12)",
		"-c:v", "bmp", "-f", "rawvideo", "out.bmp",
	}
	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotFrame() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeAccurateSeeksAfterInput(t *testing.T) {
	args := ScreenshotTimeAccurate("clip.mkv", 12.5, ScreenshotOptions{OutputPath: "-", OutputType: ScreenshotOutputTypeBMP})
	want := []string{"-v", "error", "-y", "-i", "clip.mkv", "-ss", "12.5", "-frames:v", "1", "-vcodec", "bmp", "-f", "rawvideo", "-"}
	if !reflect.DeepEqual([]string(args), want) {
		t.Fatalf("args = %#v, want %#v", []string(args), want)
	}
}
