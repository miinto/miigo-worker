package command

type genericCommand struct {
	Type string											`json:"_type"`
	Payload map[string]interface{}						`json:"_data"`
}

func (c *genericCommand) GetType() string {
	return c.Type
}

func (c *genericCommand) GetPayload() map[string]interface{} {
	return c.Payload
}