package interfaces

type Handler interface {
	Handle(command Command, logger Logger) (bool,error)
	Validate(payload map[string] interface{}) (bool,error)
}