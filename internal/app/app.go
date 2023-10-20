package app

import (
	"duetfeature/config"
	"duetfeature/internal/usecase/capture"
	"duetfeature/internal/usecase/filesaver"
	"duetfeature/internal/usecase/ocr"
	"duetfeature/pkg/gocr"
	"duetfeature/pkg/screenshot"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func Run(conf *config.ServerConfiguration) {
	log.Info().Str("port", conf.BrokerPort).Msgf("starting web server ...")

	// Incoming images
	imageRequestChan := make(chan []byte)

	// Screen Capture
	capt := capture.NewCapturer(capture.WithWorker(
		&screenshot.ImageGenerator{}),
		capture.WithOutputChan(imageRequestChan),
		capture.WithInterval(2),
	)

	capt.Start()

	// OCR
	ocrInst := ocr.NewProcessor(gocr.NewGoCr(), imageRequestChan)
	ocrInst.Start()

	// Example save file functionality using visitor pattern
	capt.Accept(&filesaver.FileSaver{Filename: "test-filename.jpg"})
	ocrInst.Accept(&filesaver.FileSaver{Filename: "test-filename.jpg"})

	// Waiting syscall signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Msgf("syscall signal: " + s.String())
	}

	// Shutdown
	capt.Stop()
	close(imageRequestChan)
	ocrInst.Stop()
}
