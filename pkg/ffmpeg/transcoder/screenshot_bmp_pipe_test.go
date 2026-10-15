package transcoder

import (
	"reflect"
	"testing"
)

func TestScreenshotOutputTypeBMPUsesImagePipe(t *testing.T) {
	got := ScreenshotOutputTypeBMP.Args()
	want := []string{"-c:v", "bmp", "-f", "image2pipe"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BMP output args = %v, want %v", got, want)
	}
}
