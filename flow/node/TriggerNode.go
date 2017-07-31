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

func NewTriggerNode(meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport,activeSubscriptions *[]string,msgInStream model.MsgPipeline) model.Node {
	node := TriggerNode{ctx:ctx,transport:transport,activeSubscriptions:activeSubscriptions}
	node.isStartNode = true
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
				if !node.ctx.IsFlowRunning {
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


//func TriggerNode(ctx *model.Context,nodes []model.MetaNode, msgInStream model.MsgPipeline, transport *fimpgo.MqttTransport, activeSubscriptions *[]string) (model.Message, *model.MetaNode, error) {
//	for i := range nodes {
//		if nodes[i].Type == "trigger" {
//			log.Info("<Node> TriggerNode is listening for events . Name = ", nodes[i].Label)
//			needToSubscribe := true
//			for i := range *activeSubscriptions {
//				if (*activeSubscriptions)[i] == nodes[i].Address {
//					needToSubscribe = false
//					break
//				}
//			}
//			if needToSubscribe {
//				log.Info("<Node> Subscribing for service by address :", nodes[i].Address)
//				transport.Subscribe(nodes[i].Address)
//				*activeSubscriptions = append(*activeSubscriptions, nodes[i].Address)
//			}
//		}
//	}
//
//	for msg := range msgInStream {
//		log.Info("<Node> New message from msgInStream")
//		if !ctx.IsFlowRunning {
//			break
//		}
//		for i := range nodes {
//			if nodes[i].Type == "trigger" {
//				if (msg.AddressStr == nodes[i].Address || nodes[i].Address == "*") &&
//					(msg.Payload.Service == nodes[i].Service || nodes[i].Service == "*") &&
//					(msg.Payload.Type == nodes[i].ServiceInterface || nodes[i].ServiceInterface == "*") {
//					//log.Info("New message.")
//					return msg, &nodes[i], nil
//				}
//			}
//
//		}
//	}
//	return model.Message{}, nil, nil
//}