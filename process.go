package gateway

type Processor interface {
	Process(handler Handler, req Request, respond bool) bool
}
