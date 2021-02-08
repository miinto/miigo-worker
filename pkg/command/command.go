package command

import (
	"encoding/json"
)

type Command struct {
	Type string											`json:"_type"`
	Payload json.RawMessage								`json:"_data"`
}

func CreateFromJson(json_payload string) (*Command, error) {
	cmd := &Command{}
	er := json.Unmarshal([]byte(json_payload), &cmd)
	if er != nil {
		return cmd, er
	}
	return cmd, nil
}