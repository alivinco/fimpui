package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
)

type LoopNode struct {
	BaseNode
	ctx       *model.Context
	transport *fimpgo.MqttTransport
	config    LoopNodeConfig
	counter   int64
	countUp   bool
}

type LoopNodeConfig struct {
	StartValue     int64
	EndValue       int64
	Step           int64
	SaveToVariable bool
}

func NewLoopNode(flowOpCtx *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node {
	node := LoopNode{ctx: ctx, transport: transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	return &node
}

func (node *LoopNode) LoadNodeConfig() error {
	defValue := LoopNodeConfig{}
	err := mapstructure.Decode(node.meta.Config, &defValue)
	if err != nil {
		log.Error(node.flowOpCtx.FlowId+"<LoopNode> Can't decode configuration", err)
	} else {
		node.config = defValue
		if node.config.Step == 0 {
			node.config.Step = 1
		}
		if defValue.EndValue > defValue.StartValue {
			node.countUp = true
			if node.config.Step > defValue.EndValue {
				node.config.Step = defValue.EndValue
			}
		} else {
			if node.config.Step < defValue.EndValue {
				node.config.Step = defValue.EndValue
			}
		}
	}
	return nil
}

func (node *LoopNode) OnInput(msg *model.Message) ([]model.NodeID, error) {
	log.Debug(node.flowOpCtx.FlowId+"<LoopNode> Executing LoopNode . Name = ", node.meta.Label)
	if node.countUp {
		node.counter = node.counter + node.config.Step
	} else {
		node.counter = node.counter - node.config.Step
	}
	log.Debug(node.flowOpCtx.FlowId+"<LoopNode> value = ", node.counter)
	if (node.countUp && node.counter >= node.config.EndValue) || (!node.countUp && node.counter <= node.config.EndValue) {
		node.counter = node.config.StartValue
		return []model.NodeID{node.meta.ErrorTransition}, nil
	}
	return []model.NodeID{node.meta.SuccessTransition}, nil
}
