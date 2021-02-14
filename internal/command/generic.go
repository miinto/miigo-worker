package command

import (
	"bytes"
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
	decoder := json.NewDecoder(bytes.NewReader([]byte(json_payload)))
	decoder.DisallowUnknownFields()
	er := decoder.Decode(&cmd)

	if er != nil {
		return &genericCommand{}, er
	}

	return cmd, nil
}