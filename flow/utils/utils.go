package utils

import (
	"github.com/alivinco/fimpui/flow/model"
	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

func GenerateId(len int) string {
	return uniuri.NewLen(len)
}

func ConfigValueToNumber(valueType string,value interface{})(float64,error){
	if valueType == "int" {
		intVal,ok := value.(int64)
		if ok {
			return float64(intVal),nil
		}else {
			return 0, errors.New("Can't convert interface{} to int64")
		}
	}else
	if valueType == "float" {
		floatVal,ok := value.(float64)
		if ok {
			return floatVal , nil
		}else {
			return 0, errors.New("Can't convert interface{} to float64")
		}
	}
	return 0,errors.New("Not numeric value type")
}

func MsgValueToNumber(msg *model.Message)(float64,error) {
	if msg.Payload.ValueType == "int" {
		intValue , err := msg.Payload.GetIntValue()
		if err  == nil {
			return float64(intValue) ,nil
		}
	}else if msg.Payload.ValueType == "float" {
		return msg.Payload.GetFloatValue()
	}
	return 0 , errors.New("Not numeric value type")
}