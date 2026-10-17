package transcoder

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestScreenshotThreadLimitsPreserveCFRAndVFRFrames(t *testing.T) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("ffmpeg is required for the integration fixture")
	}

	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(ffmpeg, args...)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("ffmpeg %v failed: %v\n%s", args, err, output)
		}
	}
	codec := "mpeg4"
	if output, err := exec.Command(ffmpeg, "-hide_banner", "-encoders").CombinedOutput(); err == nil && strings.Contains(string(output), "libx264") {
		codec = "libx264"
	}
	encode := func(source, duration, output string) {
		args := []string{"-v", "error", "-f", "lavfi", "-i", source, "-t", duration, "-c:v", codec}
		if codec == "libx264" {
			args = append(args, "-preset", "ultrafast", "-crf", "28", "-pix_fmt", "yuv420p")
		} else {
			args = append(args, "-q:v", "5")
		}
		run(append(args, "-y", output)...)
	}

	cfr := filepath.Join(dir, "cfr.mp4")
	vfrA := filepath.Join(dir, "vfr-a.mp4")
	vfrB := filepath.Join(dir, "vfr-b.mp4")
	vfrList := filepath.Join(dir, "vfr.txt")
	vfr := filepath.Join(dir, "vfr.mkv")
	encode("testsrc2=size=160x90:rate=10", "2", cfr)
	encode("testsrc2=size=160x90:rate=5", "1", vfrA)
	encode("testsrc2=size=160x90:rate=20", "1", vfrB)
	if err := os.WriteFile(vfrList, []byte("file '"+vfrA+"'\nfile '"+vfrB+"'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	run("-v", "error", "-f", "concat", "-safe", "0", "-i", vfrList, "-fps_mode", "vfr", "-c", "copy", "-y", vfr)

	for _, fixture := range []struct {
		name  string
		input string
		times []string
	}{
		{name: "cfr", input: cfr, times: []string{"0.25", "1.25"}},
		{name: "vfr", input: vfr, times: []string{"0.25", "1.2", "1.8"}},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			for index, timestamp := range fixture.times {
				for _, mode := range []struct {
					name        string
					uncappedPre bool
					capped      func(string) []string
				}{
					{name: "fast", uncappedPre: true, capped: func(output string) []string {
						return ScreenshotTime(fixture.input, mustFloat(t, timestamp), ScreenshotOptions{OutputPath: output, OutputType: ScreenshotOutputTypeImage2})
					}},
					{name: "accurate", capped: func(output string) []string {
						return ScreenshotTimeAccurate(fixture.input, mustFloat(t, timestamp), ScreenshotOptions{OutputPath: output, OutputType: ScreenshotOutputTypeImage2})
					}},
					{name: "slow", capped: func(output string) []string {
						return ScreenshotTime(fixture.input, mustFloat(t, timestamp), ScreenshotOptions{OutputPath: output, OutputType: ScreenshotOutputTypeImage2, SlowSeek: true})
					}},
				} {
					uncapped := filepath.Join(dir, fixture.name+"-"+timestamp+"-"+mode.name+"-uncapped.bmp")
					capped := filepath.Join(dir, fixture.name+"-"+timestamp+"-"+mode.name+"-capped.bmp")
					uncappedArgs := []string{"-v", "error"}
					if mode.uncappedPre {
						uncappedArgs = append(uncappedArgs, "-ss", timestamp)
					}
					uncappedArgs = append(uncappedArgs, "-i", fixture.input)
					if !mode.uncappedPre {
						uncappedArgs = append(uncappedArgs, "-ss", timestamp)
					}
					uncappedArgs = append(uncappedArgs, "-frames:v", "1", "-f", "image2", "-y", uncapped)
					cappedArgs := mode.capped(capped)
					for _, args := range [][]string{uncappedArgs, cappedArgs} {
						cmd := exec.Command(ffmpeg, args...)
						if output, err := cmd.CombinedOutput(); err != nil {
							t.Fatalf("ffmpeg %s screenshot %d failed: %v\n%s", mode.name, index, err, output)
						}
					}
					if got, want := fileDigest(t, capped), fileDigest(t, uncapped); got != want {
						t.Fatalf("timestamp %s mode %s changed frame digest: capped=%s uncapped=%s", timestamp, mode.name, got, want)
					}
				}
			}
		})
	}
}

func mustFloat(t *testing.T, value string) float64 {
	t.Helper()
	var result float64
	if _, err := fmt.Sscan(value, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func fileDigest(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatalf("empty screenshot output: %s", path)
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
