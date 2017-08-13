package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
)

type CounterNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	config CounterNodeConfig
	counter int64
	countUp bool
}

type CounterNodeConfig struct {
	StartValue int64
	EndValue int64
	EndValueTransition model.NodeID
	SaveToVariable bool

}

func NewCounterNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := CounterNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	return &node
}

func (node *CounterNode) LoadNodeConfig() error {
	defValue := CounterNodeConfig{}
	err := mapstructure.Decode(node.meta.Config,&defValue)
	if err != nil{
		log.Error("<CounterNode> Can't decode configuration",err)
	}else {
		node.config = defValue
		if defValue.EndValue > defValue.StartValue {
			node.countUp = true
		}
	}
	return nil
}

func (node *CounterNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Debug("<CounterNode> Executing CounterNode . Name = ", node.meta.Label)
	if node.countUp{
		node.counter++
	} else {
		node.counter--
	}
	log.Debug("<CounterNode> New counter value = ",node.counter )
	log.Debug("<CounterNode> End value = ",node.config.EndValue )
	if (node.countUp && node.counter >= node.config.EndValue) || (!node.countUp && node.counter <= node.config.EndValue) {
		node.counter = node.config.StartValue
		log.Debug("<CounterNode> Doing counter reset ")
		return []model.NodeID{node.config.EndValueTransition},nil
	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

