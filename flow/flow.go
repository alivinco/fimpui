package flow

import "github.com/alivinco/fimpgo"
import (
	log "github.com/Sirupsen/logrus"
)

type Definition struct {
	Id NodeID
	Type string
	Label string
	Address string

}

type MsgPipeline chan Message

type Flow struct {
	Label string
	globalContext *Context
	localContext Context
	currentNodeId NodeID
	currentMsg Message
	currentNode *Node
	nodes []Node
	msgPipeline MsgPipeline
	msgTransport *fimpgo.MqttTransport
	activeSubscriptions []string
	msgInStream MsgPipeline
	isRunning bool
}


func NewFlow(globalContext *Context,msgTransport *fimpgo.MqttTransport) *Flow  {
	flow := Flow{globalContext:globalContext}
	flow.msgPipeline = make (MsgPipeline)
	flow.nodes = make([]Node,0)
	flow.msgTransport = msgTransport
	return &flow
}

func (fl *Flow) SetNodes(nodes []Node) {
	fl.nodes = nodes
}
func (fl *Flow) AddNode(node Node) {
	fl.nodes = append(fl.nodes,node)
}

func (fl *Flow) Run() {

	var transitionNode NodeID
	for {
		if !fl.isRunning {
			break
		}
		for i := range fl.nodes {
			if fl.currentNodeId == "" && fl.nodes[i].Type == "trigger" {
				log.Info("------Flow started and waiting for trigger event----------- ")
				var err error
				fl.currentMsg ,fl.currentNode ,err = Trigger(fl.nodes,fl.msgInStream,fl.msgTransport,&fl.activeSubscriptions)
				if err != nil {
					log.Error("Trigger failed with error :",err)
					fl.currentNodeId = ""
				}
				log.Info("Trigger moving forward")
				fl.currentNodeId = fl.currentNode.Id
				transitionNode = fl.currentNode.SuccessTransition
			}else if fl.nodes[i].Id == transitionNode {
				switch fl.nodes[i].Type {
				case "action":
					log.Info("Executing Action node.")
					Action(&fl.nodes[i],&fl.currentMsg,fl.msgTransport)
				case "wait":
					log.Info("Executing Wait node.")
					Wait(&fl.nodes[i])
				}
				fl.currentNodeId = fl.nodes[i].Id
				fl.currentNode = &fl.nodes[i]
				transitionNode = fl.nodes[i].SuccessTransition
			}else if transitionNode == "" {
				// Flow is finished . Returning to first step.
				fl.currentNodeId = ""
			}
		}
	}

}
func (fl *Flow) Start(){
	log.Info("Starting flow  ",fl.Label)
	fl.isRunning = true
	go fl.Run()
}
func (fl *Flow) Stop(){
	log.Info("Stopping flow  ",fl.Label)
	fl.isRunning = false
}
func (fl *Flow) SetMessageStream(msgInStream MsgPipeline){
	fl.msgInStream = msgInStream
}
