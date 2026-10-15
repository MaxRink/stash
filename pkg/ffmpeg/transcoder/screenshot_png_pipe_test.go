package transcoder

import (
	"reflect"
	"testing"
)

func TestScreenshotOutputTypePNGUsesImagePipe(t *testing.T) {
	got := ScreenshotOutputTypePNG.Args()
	want := []string{"-c:v", "png", "-f", "image2pipe"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PNG output args = %v, want %v", got, want)
	}
}
