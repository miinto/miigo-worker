package command

import (
	"encoding/json"
	interfaces "miinto.com/miigo/worker/pkg"
)

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

func NewGenericCommand(json_payload string) (interfaces.Command, error) {
	cmd := &genericCommand{}
	er := json.Unmarshal([]byte(json_payload), &cmd)
	if er != nil {
		return cmd, er
	}
	return cmd, nil
}