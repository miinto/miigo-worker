package command

import (
	"reflect"
	"strconv"
	"testing"
)

func TestNewGenericCommand_Success(t *testing.T) {
	input := []struct{
		Json string; Type string; Payload map[string]interface{}
	}{
		{`{"_type":"Command\\BasicCommand","_data":{"foo":"bar"}}`, `Command\BasicCommand`, map[string]interface{}{"foo": "bar"}},
		{`{"_type":"Command\\BasicCommand","_data":{}}`, `Command\BasicCommand`, map[string]interface{}{}},
		{`{"_type":"Command\\BasicCommand", "_data":{"foo":"bar"}}`, `Command\BasicCommand`, map[string]interface{}{"foo": "bar"}},
		{`{"_type":"Command\\BasicCommand","_data":{"foąśo":"barąłęð"}}`, `Command\BasicCommand`, map[string]interface{}{"foąśo": "barąłęð"}},
		{`{"_type":"Command\\BasicCommand","_data":{"foo":{"innerfoo":"bar"}}}`, `Command\BasicCommand`, map[string]interface{}{"foo": map[string]interface{}{"innerfoo": "bar"}}},
		{
			`{"_type":"Command\\BasicCommand","_data":{"foo":[{"innerfoo":"bar"},{"innerfoo2":"bar2"}]}}`,
			`Command\BasicCommand`,
			map[string]interface{}{"foo": []interface{}{map[string]interface{}{"innerfoo": "bar"}, map[string]interface{}{"innerfoo2": "bar2"}}},
		},
	}

	for index,val := range input {
		t.Run("TestNewGenericCommand_Success ["+strconv.Itoa(index)+"]", func(t *testing.T){
			cmd,_ := NewGenericCommand(val.Json)
			if cmd.GetType() != val.Type {
				t.Error("incorrect command type")
			}

			if reflect.DeepEqual(cmd.GetPayload(), val.Payload) != true {
				t.Error("incorrect payload")
			}
		})
	}
}

func TestNewGenericCommand_Fail(t *testing.T) {
	input := []struct{
		Json string; Type string; Payload map[string]interface{}
	}{
		{`{"_type":"Command\\BasicCommand","_data":[]}`, `Command\BasicCommand`, map[string]interface{}{}},
		{`{"type":"Command\\BasicCommand","_data":{"foo":"bar"}}`, `Command\BasicCommand`, map[string]interface{}{"foo": "bar"}},
		{`{"type":"Command\\BasicCommand"},"_data":{"foo":"bar"}`, `Command\BasicCommand`, map[string]interface{}{"foo": "bar"}},
		{`"_data":{"foo":"bar"}`, `Command\BasicCommand`, map[string]interface{}{"foo": "bar"}},
	}

	for index,val := range input {
		t.Run("TestNewGenericCommand_Fail ["+strconv.Itoa(index)+"]", func(t *testing.T){
			cmd,err := NewGenericCommand(val.Json)
			if err == nil {
				t.Error("unmarshalling error not returned as expected")
			}

			if cmd.GetType() != "" {
				t.Error("failed command unmarshalling should return an empty command type")
			}

			if reflect.DeepEqual(cmd.GetPayload(), map[string]interface{}{}) {
				t.Error("failed command unmarshalling should return an empty command payload")
			}
		})
	}
}