package capture

import (
	"duetfeature/internal/domain"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"
)

type Worker interface {
	Capture() ([][]byte, error)
}

type Capturer struct {
	worker     Worker
	interval   time.Duration
	ticker     *time.Ticker
	stopChan   chan bool
	outputChan chan<- []byte
}

func NewCapturer(options ...func(*Capturer)) *Capturer {
	concreteCapturer := &Capturer{
		stopChan: make(chan bool),
	}
	for _, o := range options {
		o(concreteCapturer)
	}
	return concreteCapturer
}

func WithWorker(worker Worker) func(*Capturer) {
	return func(s *Capturer) {
		s.worker = worker
	}
}

func WithOutputChan(outputChan chan<- []byte) func(*Capturer) {
	return func(s *Capturer) {
		s.outputChan = outputChan
	}
}

func WithInterval(interval time.Duration) func(*Capturer) {
	return func(s *Capturer) {
		s.interval = interval
	}
}

func (s *Capturer) Start() {
	s.ticker = time.NewTicker(s.interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				data, err := s.worker.Capture()
				if err != nil {
					log.Error().Err(err).Msg("Capturer")
					continue
				}
				for _, v := range data {
					s.outputChan <- v
				}
			case <-s.stopChan:
				fmt.Printf("%s\n", "Capturer stopped ...")
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *Capturer) Accept(v domain.Visitor) {
	v.VisitCapture(s)
}

func (s *Capturer) Stop() {
	close(s.stopChan)
	log.Info().Msgf("Capturer stopped ...")
}
