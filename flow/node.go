package flow

import "github.com/alivinco/fimpgo"
import (
	log "github.com/Sirupsen/logrus"
	//"time"
	"time"
)

type NodeID string

type Node struct {
	Id                NodeID
	Type              string
	Label             string
	SuccessTransition NodeID
	TimeoutTransition NodeID
	ErrorTransition   NodeID
	Address           string
	Service           string
	ServiceInterface  string
	Config            interface{}
}

func Trigger(nodes []Node, msgInStream MsgPipeline, transport *fimpgo.MqttTransport, activeSubscriptions *[]string) (Message , *Node , error) {
	for i := range nodes {
		if nodes[i].Type == "trigger" {
			log.Info("Trigger is listening for events . Name = ", nodes[i].Label)
			needToSubscribe := true
			for i := range *activeSubscriptions {
				if (*activeSubscriptions)[i] == nodes[i].Address {
					needToSubscribe = false
					break
				}
			}
			if needToSubscribe {
				log.Info("Subscribing for service by address :", nodes[i].Address)
				transport.Subscribe(nodes[i].Address)
				*activeSubscriptions = append(*activeSubscriptions, nodes[i].Address)
			}
		}
	}


	for msg := range msgInStream {
		log.Info("New message from msgInStream")
		for i := range nodes {
			if nodes[i].Type == "trigger" {
				if (msg.AddressStr == nodes[i].Address || nodes[i].Address == "*") &&
					(msg.Payload.Service == nodes[i].Service || nodes[i].Service == "*") &&
					(msg.Payload.Type == nodes[i].ServiceInterface || nodes[i].ServiceInterface == "*") {
					//log.Info("New message.")
					return msg,&nodes[i],nil
				}
			}

		}
	}
	return Message{},nil,nil
}

func Action(node *Node, msg *Message, transport *fimpgo.MqttTransport) error {
	log.Info("Executing Action . Name = ",node.Label)
	fimpMsg := fimpgo.FimpMessage{Type:node.ServiceInterface,Service:node.Service,Value:msg.Payload.Value,ValueType:msg.Payload.ValueType}
	msgBa ,err := fimpMsg.SerializeToJson()
	if err != nil {
		return err
	}
	transport.PublishRaw(node.Address,msgBa)
	return nil
}

func Wait(node *Node) error {
	delayMilisec,ok:= node.Config.(int)
	if ok {
		log.Info("Waiting  for = ",delayMilisec)
		time.Sleep(time.Millisecond*time.Duration(delayMilisec))
	}else {
		log.Error("Wrong time format")
	}

	return nil
}
