package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
)

type ActionNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
}

func NewActionNode(meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := ActionNode{ctx:ctx,transport:transport}
	node.meta = meta
	return &node
}

func (node *ActionNode) LoadNodeConfig() error {
	defValue := model.Variable{}
	err := mapstructure.Decode(node.meta.Config,&defValue)
	if err != nil{
		log.Error(err)
	}else {
		node.meta.Config = defValue
	}
	return nil
}

func (node *ActionNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info("<Node> Executing ActionNode . Name = ", node.meta.Label)
	fimpMsg := fimpgo.FimpMessage{Type: node.meta.ServiceInterface, Service: node.meta.Service}
	defaultValue, ok := node.meta.Config.(model.Variable)
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
		return nil,err
	}
	log.Debug("<Node> Action message :", fimpMsg)
	node.transport.PublishRaw(node.meta.Address, msgBa)
	return []model.NodeID{node.meta.SuccessTransition},nil
}

