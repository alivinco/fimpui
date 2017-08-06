package node

import "time"
import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/alivinco/fimpgo"
)

type WaitNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
}

func NewWaitNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := WaitNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx  = flowOpCtx
	return &node
}

func (node *WaitNode) LoadNodeConfig() error {
	delay ,ok := node.meta.Config.(float64)
	if ok {
		node.meta.Config = int(delay)
	}else {
		log.Error("<FlMan> Can't cast Wait node delay value")
	}

	return nil
}

func (node *WaitNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	delayMilisec, ok := node.meta.Config.(int)
	if ok {
		log.Info("<Node> Waiting  for = ", delayMilisec)
		time.Sleep(time.Millisecond * time.Duration(delayMilisec))
	} else {
		log.Error("<Node> Wrong time format")
	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

//func WaitNode(ctx *model.Context,node *model.MetaNode) error {
//	delayMilisec, ok := node.Config.(int)
//	if ok {
//		log.Info("<Node> Waiting  for = ", delayMilisec)
//		time.Sleep(time.Millisecond * time.Duration(delayMilisec))
//	} else {
//		log.Error("<Node> Wrong time format")
//	}
//
//	return nil
//}

