package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
)

type TriggerNode struct {
	BaseNode
	ctx                 *model.Context
	transport           *fimpgo.MqttTransport
	activeSubscriptions *[]string
	msgInStream         model.MsgPipeline
	config              TriggerConfig
}

type TriggerConfig struct {
	ValueFilter model.Variable
}

func NewTriggerNode(flowOpCtx *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node {
	node := TriggerNode{ctx: ctx, transport: transport}
	node.isStartNode = true
	node.isMsgReactor = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.config = TriggerConfig{}
	return &node
}

func (node *TriggerNode) ConfigureInStream(activeSubscriptions *[]string, msgInStream model.MsgPipeline) {
	log.Info("<TrigNode>Configuring Stream")
	node.activeSubscriptions = activeSubscriptions
	node.msgInStream = msgInStream
	node.initSubscriptions()
}

func (node *TriggerNode) initSubscriptions() {
	if node.meta.Type == "trigger" {
		log.Info("<TrigNode> TriggerNode is listening for events . Name = ", node.meta.Label)
		needToSubscribe := true
		for i := range *node.activeSubscriptions {
			if (*node.activeSubscriptions)[i] == node.meta.Address {
				needToSubscribe = false
				break
			}
		}
		if needToSubscribe {
			log.Info("<TrigNode> Subscribing for service by address :", node.meta.Address)
			node.transport.Subscribe(node.meta.Address)
			*node.activeSubscriptions = append(*node.activeSubscriptions, node.meta.Address)
		}
	}
}

func (node *TriggerNode) LoadNodeConfig() error {
	return nil
}

func (node *TriggerNode) OnInput(msg *model.Message) ([]model.NodeID, error) {
	log.Debug("<TrigNode> Waiting for event ")
	for {
		select {
		case newMsg := <-node.msgInStream:
			log.Debug("<TrigNode> New message from InStream ")
			*msg = newMsg
			if node.config.ValueFilter.ValueType == "" {
				return []model.NodeID{node.meta.SuccessTransition}, nil
			} else if newMsg.Payload.Value == node.config.ValueFilter.Value {
				return []model.NodeID{node.meta.SuccessTransition}, nil
			}
		case signal := <-node.flowOpCtx.NodeControlSignalChannel:
			log.Debug("<TrigNode> Control signal ")
			if signal == model.SIGNAL_STOP {
				return nil, nil
			}
		}
	}
}
