package domain

type ProcessorInterface interface {
	Accept(Visitor)
}

type CapturerInterface interface {
	Accept(Visitor)
}
