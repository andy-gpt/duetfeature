package tesseract

import "github.com/otiai10/gosseract/v2"

type Gosseract struct {
	client *gosseract.Client
}

func NewGosseract() *Gosseract {
	return &Gosseract{
		client: gosseract.NewClient(),
	}
}

func (g *Gosseract) SetImage(blob []byte) error {
	return g.client.SetImageFromBytes(blob)
}

func (g *Gosseract) GetText(blob []byte) (string, error) {
	g.client.SetImageFromBytes(blob)
	return g.client.Text()
}
