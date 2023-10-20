package domain

type Visitor interface {
	VisitCapture(CapturerInterface)
	VisitOCR(ProcessorInterface)
}
