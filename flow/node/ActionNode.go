package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
)

type ActionNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	config ActionNodeConfig
}

type ActionNodeConfig struct {
	DefaultValue model.Variable
	VariableName string
	VariableType string
	IsVariableGlobal bool
	Props fimpgo.Props
	RegisterAsVirtualService bool
	VirtualServiceGroup string
}

func NewActionNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := ActionNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	node.config = ActionNodeConfig{DefaultValue:model.Variable{}}
	node.SetupBaseNode()
	return &node
}

func (node *ActionNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		node.getLog().Error("Can't decode config.Err:",err)

	}
	return err
}

func (node *ActionNode) WaitForEvent(responseChannel chan model.ReactorEvent) {

}

func (node *ActionNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	node.getLog().Info("Executing ActionNode . Name = ", node.meta.Label)
	fimpMsg := fimpgo.FimpMessage{Type: node.meta.ServiceInterface, Service: node.meta.Service,Properties:node.config.Props}
	if node.config.VariableName != "" {
		flowId := node.flowOpCtx.FlowId
		if node.config.IsVariableGlobal {
			flowId = "global"
		}
		variable,err := node.ctx.GetVariable(node.config.VariableName,flowId)
		if err != nil {
			node.getLog().Error("Can't get variable . Error:",err)
			return nil , err
		}
		fimpMsg.ValueType = variable.ValueType
		fimpMsg.Value = variable.Value
	}else {
		if node.config.DefaultValue.Value == "" || node.config.DefaultValue.ValueType == ""{
			fimpMsg.Value = msg.Payload.Value
			fimpMsg.ValueType = msg.Payload.ValueType
		}else {
			fimpMsg.Value = node.config.DefaultValue.Value
			fimpMsg.ValueType = node.config.DefaultValue.ValueType
		}
	}

	msgBa, err := fimpMsg.SerializeToJson()
	if err != nil {
		return nil,err
	}
	node.getLog().Debug(" Action message :", fimpMsg)
	node.transport.PublishRaw(node.meta.Address, msgBa)
	return []model.NodeID{node.meta.SuccessTransition},nil
}

