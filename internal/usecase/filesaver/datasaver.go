package filesaver

import (
	"duetfeature/internal/domain"
	"fmt"
)

type FileSaver struct {
	Filename string
}

func (fs *FileSaver) VisitCapture(c domain.CapturerInterface) {
	fmt.Println("file saved to disk...")
}

func (fs *FileSaver) VisitOCR(p domain.ProcessorInterface) {
	fmt.Println("file saved to disk...")
}
