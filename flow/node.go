package flow

import "github.com/alivinco/fimpgo"
import (
	log "github.com/Sirupsen/logrus"
	//"time"
	"github.com/pkg/errors"
	"time"
)

type NodeID string

type MetaNode struct {
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

func TriggerNode(nodes []MetaNode, ctx *Context, msgInStream MsgPipeline, transport *fimpgo.MqttTransport, activeSubscriptions *[]string) (Message, *MetaNode, error) {
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
		if !ctx.isFlowRunning {
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
	return Message{}, nil, nil
}

type DefaultValue struct {
	Value     interface{}
	ValueType string
}

func ActionNode(node *MetaNode, msg *Message, transport *fimpgo.MqttTransport) error {
	log.Info("<Node> Executing ActionNode . Name = ", node.Label)
	fimpMsg := fimpgo.FimpMessage{Type: node.ServiceInterface, Service: node.Service}
	defaultValue, ok := node.Config.(DefaultValue)
	if ok {
		fimpMsg.Value = defaultValue.Value
		fimpMsg.ValueType = defaultValue.ValueType
	} else {
		fimpMsg.Value = msg.Payload.Value
		fimpMsg.ValueType = msg.Payload.ValueType
	}

	msgBa, err := fimpMsg.SerializeToJson()
	if err != nil {
		return err
	}
	log.Debug("<Node> Action message :", fimpMsg)
	transport.PublishRaw(node.Address, msgBa)
	return nil
}

func WaitNode(node *MetaNode) error {
	delayMilisec, ok := node.Config.(int)
	if ok {
		log.Info("<Node> Waiting  for = ", delayMilisec)
		time.Sleep(time.Millisecond * time.Duration(delayMilisec))
	} else {
		log.Error("<Node> Wrong time format")
	}

	return nil
}

type IFExpressions struct {
	Expression      []IFExpression
	TrueTransition  NodeID
	FalseTransition NodeID
}

type IFExpression struct {
	Operand         string // eq , gr , lt
	Value           interface{}
	ValueType       string
	BooleanOperator string // and , or , not
}

func IfNode(node *MetaNode, msg *Message) error {
	conf, ok := node.Config.(IFExpressions)
	if ok {
		booleanOperator := ""
		var finalResult bool
		for _, item := range conf.Expression {
			if item.ValueType != msg.Payload.ValueType {
				return errors.New("Incompatible value type")
			}
			var result bool
			log.Info("<Node> Operand = ", item.Operand)
			switch item.Operand {
			case "eq":
				result = item.Value == msg.Payload.Value
			case "gt":
				val, err := msg.Payload.GetIntValue()
				if err != nil {
					return err
				}
				result = float64(val) > item.Value.(float64)
			case "lt":
				val, err := msg.Payload.GetIntValue()
				if err != nil {
					return err
				}
				result = float64(val) < item.Value.(float64)

			}
			switch booleanOperator {
			case "":
				finalResult = result
			case "and":
				finalResult = finalResult && result
				booleanOperator = item.BooleanOperator
			case "or":
				finalResult = finalResult || result
				booleanOperator = item.BooleanOperator
			case "not":
				finalResult = !result
			}
		}
		if finalResult {
			node.SuccessTransition = conf.TrueTransition
		} else {
			node.SuccessTransition = conf.FalseTransition
		}
		return nil
	} else {
		log.Error(node.Config)
		return errors.New("Incompatible node configuration format")
	}

}
