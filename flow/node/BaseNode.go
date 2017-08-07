package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
)

type BaseNode struct {

	meta model.MetaNode
	ctx *model.Context
	flowOpCtx *model.FlowOperationalContext
	isStartNode bool   // true - if node is first in a flow
	isMsgReactor bool  // true - node reacts on messages and requires input stream .
	transport *fimpgo.MqttTransport

}

func (node *BaseNode) GetMetaNode()*model.MetaNode {
	return &node.meta
}
func (node *BaseNode) GetNextSuccessNodes()[]model.NodeID {
	return []model.NodeID{node.meta.SuccessTransition}
}

func (node *BaseNode) GetNextErrorNode()model.NodeID {
	return node.meta.ErrorTransition
}

func (node *BaseNode) GetNextTimeoutNode()model.NodeID{
	return node.meta.TimeoutTransition
}

func (node *BaseNode) IsStartNode() bool {
	return node.isStartNode
}

func (node *BaseNode) IsMsgReactorNode() bool {
	return node.isMsgReactor
}

func (node *BaseNode) ConfigureInStream(activeSubscriptions *[]string,msgInStream model.MsgPipeline) {
}