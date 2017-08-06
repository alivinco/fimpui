package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	log "github.com/Sirupsen/logrus"
	//"github.com/mitchellh/mapstructure"
)

type TriggerNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	activeSubscriptions *[]string
	msgInStream model.MsgPipeline
}

func NewTriggerNode(flowOpCtx *model.FlowOperationalContext ,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport,activeSubscriptions *[]string,msgInStream model.MsgPipeline) model.Node {
	node := TriggerNode{ctx:ctx,transport:transport,activeSubscriptions:activeSubscriptions}
	node.isStartNode = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.msgInStream = msgInStream
	node.initSubscriptions()
	return &node
}

func (node *TriggerNode) initSubscriptions() {
		if node.meta.Type == "trigger" {
			log.Info("<Node> TriggerNode is listening for events . Name = ", node.meta.Label)
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
}


func (node *TriggerNode) LoadNodeConfig() error {
	return nil
}

func (node *TriggerNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	for inMsg := range node.msgInStream {
				log.Info("<Node> New message from msgInStream")
				if !node.flowOpCtx.IsFlowRunning {
					break
				}
				if node.meta.Type == "trigger" {
						if (inMsg.AddressStr == node.meta.Address || node.meta.Address == "*") &&
							(inMsg.Payload.Service == node.meta.Service || node.meta.Service == "*") &&
							(inMsg.Payload.Type == node.meta.ServiceInterface || node.meta.ServiceInterface == "*") {
							//log.Info("New message.")
							*msg = inMsg
							return []model.NodeID{node.meta.SuccessTransition}, nil
						}
					}
			}
	return nil, nil
}

