package screenshot

import (
	"bytes"
	"github.com/kbinani/screenshot"
	"image/png"
)

type ImageGeneratorInterface interface {
	Capture() ([][]byte, error)
}

type ImageGenerator struct {
	ch chan<- []byte
}

func (ig *ImageGenerator) Capture() ([][]byte, error) {
	n := screenshot.NumActiveDisplays()
	var capture [][]byte

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		err = png.Encode(&buf, img)
		capture = append(capture, buf.Bytes())
		if err != nil {
			return nil, err
		}

	}

	return capture, nil
}
