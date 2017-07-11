package flow

import "github.com/alivinco/fimpgo"
import (
	log "github.com/Sirupsen/logrus"
)

type MsgPipeline chan Message

type Flow struct {
	Id                  string
	Name                string
	globalContext       *Context  `json:"-"`
	localContext        Context   `json:"-"`
	currentNodeId       NodeID    `json:"-"`
	currentMsg          Message   `json:"-"`
	currentNode         *MetaNode `json:"-"`
	Nodes               []MetaNode
	msgPipeline         MsgPipeline           `json:"-"`
	msgTransport        *fimpgo.MqttTransport `json:"-"`
	activeSubscriptions []string              `json:"-"`
	msgInStream         MsgPipeline           `json:"-"`
}

func NewFlow(Id string, globalContext *Context, msgTransport *fimpgo.MqttTransport) *Flow {
	flow := Flow{globalContext: globalContext}
	flow.msgPipeline = make(MsgPipeline)
	flow.Nodes = make([]MetaNode, 0)
	flow.msgTransport = msgTransport
	flow.localContext = Context{isFlowRunning: true}
	return &flow
}

func (fl *Flow) SetNodes(nodes []MetaNode) {
	fl.Nodes = nodes
}
func (fl *Flow) AddNode(node MetaNode) {
	fl.Nodes = append(fl.Nodes, node)
}

func (fl *Flow) Run() {

	var transitionNode NodeID
	for {
		if !fl.localContext.isFlowRunning {
			break
		}
		for i := range fl.Nodes {
			if !fl.localContext.isFlowRunning {
				break
			}
			if fl.currentNodeId == "" && fl.Nodes[i].Type == "trigger" {
				log.Info("------Flow started and waiting for trigger event----------- ")
				var err error
				fl.currentMsg, fl.currentNode, err = TriggerNode(fl.Nodes, &fl.localContext, fl.msgInStream, fl.msgTransport, &fl.activeSubscriptions)
				if err != nil {
					log.Error("TriggerNode failed with error :", err)
					fl.currentNodeId = ""
				}
				if !fl.localContext.isFlowRunning {
					break
				}
				log.Info("TriggerNode moving forward")
				fl.currentNodeId = fl.currentNode.Id
				transitionNode = fl.currentNode.SuccessTransition
			} else if fl.Nodes[i].Id == transitionNode {
				var err error
				switch fl.Nodes[i].Type {
				case "action":
					log.Info("Executing ActionNode node.")
					err = ActionNode(&fl.Nodes[i], &fl.currentMsg, fl.msgTransport)
				case "wait":
					log.Info("Executing WaitNode node.")
					err = WaitNode(&fl.Nodes[i])
				case "if":
					log.Info("Executing IfNode node.")
					err = IfNode(&fl.Nodes[i], &fl.currentMsg)
				}
				fl.currentNodeId = fl.Nodes[i].Id
				fl.currentNode = &fl.Nodes[i]
				if err == nil {
					transitionNode = fl.Nodes[i].SuccessTransition
				} else {
					log.Info("Node executed with error . Doing error transition. Error :", err)
					transitionNode = fl.Nodes[i].ErrorTransition
				}

			} else if transitionNode == "" {
				// Flow is finished . Returning to first step.
				fl.currentNodeId = ""
			}
		}
	}
	log.Infof("Flow %s stopped.", fl.Name)

}
func (fl *Flow) Start() {
	log.Info("Starting flow  ", fl.Name)
	fl.localContext.isFlowRunning = true
	go fl.Run()
}
func (fl *Flow) Stop() {
	log.Info("Stopping flow  ", fl.Name)
	fl.localContext.isFlowRunning = false
	fl.msgInStream <- Message{}
}
func (fl *Flow) SetMessageStream(msgInStream MsgPipeline) {
	fl.msgInStream = msgInStream
}
