package transcoder

import (
	"reflect"
	"testing"
)

func TestScreenshotTimeDefaultUsesFastSeek(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
	}

	got := ScreenshotTime("input.webm", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-threads", "1",
		"-filter_threads", "1",
		"-filter_complex_threads", "1",
		"-ss", "12.5",
		"-i", "input.webm",
		"-frames:v", "1",
		"-f", "image2",
		"out.jpg",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeSlowSeek(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
		SlowSeek:   true,
	}

	got := ScreenshotTime("input.webm", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-threads", "1",
		"-filter_threads", "1",
		"-filter_complex_threads", "1",
		"-i", "input.webm",
		"-ss", "12.5",
		"-frames:v", "1",
		"-f", "image2",
		"out.jpg",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeAccurateSeeksAfterInput(t *testing.T) {
	args := ScreenshotTimeAccurate("clip.mkv", 12.5, ScreenshotOptions{OutputPath: "-", OutputType: ScreenshotOutputTypeBMP})
	want := []string{"-v", "error", "-y", "-threads", "1", "-filter_threads", "1", "-filter_complex_threads", "1", "-i", "clip.mkv", "-ss", "12.5", "-frames:v", "1", "-c:v", "bmp", "-f", "rawvideo", "-"}
	if !reflect.DeepEqual([]string(args), want) {
		t.Fatalf("args = %#v, want %#v", []string(args), want)
	}
}
