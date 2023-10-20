package ocr

import (
	"duetfeature/internal/domain"
	"github.com/rs/zerolog/log"
)

type Engine interface {
	GetText(blob []byte) (string, error)
}

type Processor struct {
	engine     Engine
	inputChan  <-chan []byte
	doneSignal chan struct{}
}

func NewProcessor(engine Engine, inputChan <-chan []byte) *Processor {
	return &Processor{
		engine:     engine,
		inputChan:  inputChan,
		doneSignal: make(chan struct{}),
	}
}

func (p *Processor) Start() {
	go func() {
		for imgReq := range p.inputChan {
			text, err := p.engine.GetText(imgReq)
			if err != nil {
				log.Error().Err(err)
			}
			log.Info().Msgf("%s", text)
		}
		close(p.doneSignal)
	}()
}

func (p *Processor) Accept(v domain.Visitor) {
	v.VisitOCR(p)
}

func (p *Processor) Stop() {
	<-p.doneSignal
	log.Info().Msgf("OCR Stopped ...")
}
