package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
)

type SetVariableNode struct {
	BaseNode
	ctx *model.Context
	nodeConfig SetVariableNodeConfig
	transport *fimpgo.MqttTransport
}

type SetVariableNodeConfig struct {
	Name string
	UpdateGlobal bool                    // true - update global variable ; false - update local variable
	UpdateInputMsg bool               // true - update input message  ; false - update context variable
	DefaultValue model.Variable
}

func NewSetVariableNode(meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := SetVariableNode{ctx:ctx,transport:transport}
	node.meta = meta
	return &node
}

func (node *SetVariableNode) LoadNodeConfig() error {
	defValue := SetVariableNodeConfig{}
	err := mapstructure.Decode(node.meta.Config,&defValue)
	if err != nil{
		log.Error("<SetVarNode> Can't decode configuration",err)
	}else {
		node.nodeConfig = defValue
		node.meta.Config = defValue
	}
	return nil
}

func (node *SetVariableNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info("<Node> Executing SetVariableNode . Name = ", node.meta.Label)
	// set input message value as variable value
	if node.nodeConfig.DefaultValue.ValueType == "" {
		if node.nodeConfig.UpdateInputMsg {
			msg.Payload.Value = msg.Payload.Value
			msg.Payload.ValueType = msg.Payload.ValueType
		}else {
			if node.nodeConfig.UpdateGlobal {
				node.ctx.GetParentContext().SetVariable(node.nodeConfig.Name,msg.Payload.ValueType,msg.Payload.Value)
			}else {
				node.ctx.SetVariable(node.nodeConfig.Name,msg.Payload.ValueType,msg.Payload.Value)

			}
		}

	}else {
		// set variable value to default value
		if node.nodeConfig.UpdateGlobal {
			node.ctx.GetParentContext().SetVariable(node.nodeConfig.Name,node.nodeConfig.DefaultValue.ValueType,node.nodeConfig.DefaultValue.Value)
		}else {
			node.ctx.SetVariable(node.nodeConfig.Name,node.nodeConfig.DefaultValue.ValueType,node.nodeConfig.DefaultValue.Value)
			log.Info(node.ctx.GetVariable(node.nodeConfig.Name))
		}
	}

	return []model.NodeID{node.meta.SuccessTransition},nil
}

