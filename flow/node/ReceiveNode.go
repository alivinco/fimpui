package node

import (
"github.com/alivinco/fimpgo"
"github.com/alivinco/fimpui/flow/model"
log "github.com/Sirupsen/logrus"
//"github.com/mitchellh/mapstructure"
)

type ReceiveNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	activeSubscriptions *[]string
	msgInStream model.MsgPipeline
	waitTimeout int;
}

func NewReceiveNode(flowOpCtx *model.FlowOperationalContext ,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := ReceiveNode{ctx:ctx,transport:transport}
	node.isStartNode = false
	node.isMsgReactor = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	return &node
}

func (node *ReceiveNode) ConfigureInStream(activeSubscriptions *[]string,msgInStream model.MsgPipeline) {
	node.activeSubscriptions = activeSubscriptions
	node.msgInStream = msgInStream
	node.initSubscriptions()
}

func (node *ReceiveNode) initSubscriptions() {
	log.Info("<Node> ReceiveNode is listening for events . Name = ", node.meta.Label)
	needToSubscribe := true
	for i := range *node.activeSubscriptions {
			if (*node.activeSubscriptions)[i] == node.meta.Address {
				needToSubscribe = false
				break
			}
	}
	if needToSubscribe {
			log.Info("<Node> Subscribing for service by address :", node.meta.Address)
			node.transport.Subscribe(node.meta.Address)
			*node.activeSubscriptions = append(*node.activeSubscriptions, node.meta.Address)
	}
}


func (node *ReceiveNode) LoadNodeConfig() error {
	delay ,ok := node.meta.Config.(float64)
	if ok {
		node.waitTimeout = int(delay)
	}else {
		log.Error("<FlMan> Can't cast Wait node delay value")
	}

	return nil
}

func (node *ReceiveNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	for inMsg := range node.msgInStream {
		log.Debug("<Node> New message from msgInStream")
		if !node.flowOpCtx.IsFlowRunning {
			break
		}
		if (inMsg.AddressStr == node.meta.Address || node.meta.Address == "*") &&
			(inMsg.Payload.Service == node.meta.Service || node.meta.Service == "*") &&
			(inMsg.Payload.Type == node.meta.ServiceInterface || node.meta.ServiceInterface == "*") {
			//log.Info("New message.")
			*msg = inMsg
		return []model.NodeID{node.meta.SuccessTransition}, nil
		}
	}
	return nil, nil
}

