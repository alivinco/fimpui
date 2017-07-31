package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	log "github.com/Sirupsen/logrus"
)


func TriggerNode(nodes []model.MetaNode, ctx *model.Context, msgInStream model.MsgPipeline, transport *fimpgo.MqttTransport, activeSubscriptions *[]string) (model.Message, *model.MetaNode, error) {
	for i := range nodes {
		if nodes[i].Type == "trigger" {
			log.Info("<Node> TriggerNode is listening for events . Name = ", nodes[i].Label)
			needToSubscribe := true
			for i := range *activeSubscriptions {
				if (*activeSubscriptions)[i] == nodes[i].Address {
					needToSubscribe = false
					break
				}
			}
			if needToSubscribe {
				log.Info("<Node> Subscribing for service by address :", nodes[i].Address)
				transport.Subscribe(nodes[i].Address)
				*activeSubscriptions = append(*activeSubscriptions, nodes[i].Address)
			}
		}
	}

	for msg := range msgInStream {
		log.Info("<Node> New message from msgInStream")
		if !ctx.IsFlowRunning {
			break
		}
		for i := range nodes {
			if nodes[i].Type == "trigger" {
				if (msg.AddressStr == nodes[i].Address || nodes[i].Address == "*") &&
					(msg.Payload.Service == nodes[i].Service || nodes[i].Service == "*") &&
					(msg.Payload.Type == nodes[i].ServiceInterface || nodes[i].ServiceInterface == "*") {
					//log.Info("New message.")
					return msg, &nodes[i], nil
				}
			}

		}
	}
	return model.Message{}, nil, nil
}