package command

import (
	"encoding/json"
	"miinto.com/miigo/worker/pkg"
)

func NewGenericCommand(json_payload string) (interfaces.Command, error) {
	cmd := &genericCommand{}
	er := json.Unmarshal([]byte(json_payload), &cmd)
	if er != nil {
		return cmd, er
	}
	return cmd, nil
}