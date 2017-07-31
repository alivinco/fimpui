package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
)

type DefaultValue struct {
	Value     interface{}
	ValueType string
}

func ActionNode(node *model.MetaNode, msg *model.Message, transport *fimpgo.MqttTransport) error {
	log.Info("<Node> Executing ActionNode . Name = ", node.Label)
	fimpMsg := fimpgo.FimpMessage{Type: node.ServiceInterface, Service: node.Service}
	defaultValue, ok := node.Config.(DefaultValue)
	if ok {

		if defaultValue.Value == "" || defaultValue.ValueType == ""{
			fimpMsg.Value = msg.Payload.Value
			fimpMsg.ValueType = msg.Payload.ValueType
		}else {
			fimpMsg.Value = defaultValue.Value
			fimpMsg.ValueType = defaultValue.ValueType
		}
	} else {
		fimpMsg.Value = msg.Payload.Value
		fimpMsg.ValueType = msg.Payload.ValueType
	}

	msgBa, err := fimpMsg.SerializeToJson()
	if err != nil {
		return err
	}
	log.Debug("<Node> Action message :", fimpMsg)
	transport.PublishRaw(node.Address, msgBa)
	return nil
}
