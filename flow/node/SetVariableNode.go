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
	Description string
	UpdateGlobal bool                    // true - update global variable ; false - update local variable
	UpdateInputMsg bool               // true - update input message  ; false - update context variable
	PersistOnUpdate bool              // true - is saved on disk ; false - in memory only
	DefaultValue model.Variable
}

func NewSetVariableNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := SetVariableNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	return &node
}

func (node *SetVariableNode) LoadNodeConfig() error {
	defValue := SetVariableNodeConfig{}
	err := mapstructure.Decode(node.meta.Config,&defValue)
	if err != nil{
		log.Error(node.flowOpCtx.FlowId+"<SetVarNode> Can't decode configuration",err)
	}else {
		node.nodeConfig = defValue
		node.meta.Config = defValue
	}
	return nil
}

func (node *SetVariableNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info(node.flowOpCtx.FlowId+"<Node> Executing SetVariableNode . Name = ", node.meta.Label)
	// set input message value to variable value
	if node.nodeConfig.DefaultValue.ValueType == "" {
		if node.nodeConfig.UpdateInputMsg {
			msg.Payload.Value = msg.Payload.Value
			msg.Payload.ValueType = msg.Payload.ValueType
		}else {
			if node.nodeConfig.UpdateGlobal {
				node.ctx.SetVariable(node.nodeConfig.Name,msg.Payload.ValueType,msg.Payload.Value,node.nodeConfig.Description,"global",false)
			}else {
				node.ctx.SetVariable(node.nodeConfig.Name,msg.Payload.ValueType,msg.Payload.Value,node.nodeConfig.Description,node.flowOpCtx.FlowId,false)

			}
		}

	}else {
		// set variable value to default value
		if node.nodeConfig.UpdateGlobal {
			node.ctx.SetVariable(node.nodeConfig.Name,node.nodeConfig.DefaultValue.ValueType,node.nodeConfig.DefaultValue.Value,node.nodeConfig.Description,"global",false)
		}else {
			node.ctx.SetVariable(node.nodeConfig.Name,node.nodeConfig.DefaultValue.ValueType,node.nodeConfig.DefaultValue.Value,node.nodeConfig.Description,node.flowOpCtx.FlowId,false)
		}
	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

