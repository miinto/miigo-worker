package interfaces

type Command interface {
	GetType() string
	GetPayload() map[string]interface{}
}