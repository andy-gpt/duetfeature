package screenshot

import (
	"github.com/kbinani/screenshot"
	"testing"
)

func TestCaptureNotEmpty(t *testing.T) {
	gen := ImageGenerator{}

	capture, err := gen.Capture()
	if err != nil {
		t.Errorf("Failed to capture screen: %s", err)
		return
	}

	for i, imgBytes := range capture {
		if len(imgBytes) == 0 {
			t.Errorf("Capture for display %d is empty", i)
		}
	}

}

func TestIfActiveDisplaysFound(t *testing.T) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		t.Error("No active displays found")
		return
	}
}
