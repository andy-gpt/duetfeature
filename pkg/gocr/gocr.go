package gocr

type GoCr struct {
}

func NewGoCr() *GoCr {
	return &GoCr{}
}

func (g *GoCr) GetText(blob []byte) (string, error) {
	return "ocr test", nil
}
