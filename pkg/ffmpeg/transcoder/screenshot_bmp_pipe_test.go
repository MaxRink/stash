package transcoder

import (
	"reflect"
	"testing"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

func TestScreenshotOutputTypeBMPUsesImagePipe(t *testing.T) {
	got := ScreenshotOutputTypeBMP.Args()
	want := []string{"-c:v", "bmp", "-f", "image2pipe"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BMP output args = %v, want %v", got, want)
	}
	_ = ffmpeg.FormatImage2Pipe
}
