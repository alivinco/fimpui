package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/node"
	"github.com/alivinco/fimpui/flow/model"
)

type Flow struct {
	Id                  string
	Name                string
	Description         string
	globalContext       *model.Context  `json:"-"`
	localContext        model.Context   `json:"-"`
	currentNodeId       model.NodeID    `json:"-"`
	currentMsg          model.Message   `json:"-"`
	currentNode         *model.MetaNode `json:"-"`
	Nodes               []model.MetaNode
	msgPipeline         model.MsgPipeline           `json:"-"`
	msgTransport        *fimpgo.MqttTransport `json:"-"`
	activeSubscriptions []string              `json:"-"`
	msgInStream         model.MsgPipeline           `json:"-"`
	TriggerCounter      int64				  `json:"-"`
	ErrorCounter        int64				  `json:"-"`
	isFlowRunning       bool
}

func NewFlow(Id string, globalContext *model.Context, msgTransport *fimpgo.MqttTransport) *Flow {
	flow := Flow{globalContext: globalContext}
	flow.msgPipeline = make(model.MsgPipeline)
	flow.Nodes = make([]model.MetaNode, 0)
	flow.msgTransport = msgTransport
	flow.localContext = model.NewContext()
	flow.localContext.IsFlowRunning = true
	return &flow
}

func (fl *Flow) SetNodes(nodes []model.MetaNode) {
	fl.Nodes = nodes
}

func (fl *Flow) ReloadNodes(nodes []model.MetaNode) {
	fl.Stop()
	fl.Nodes = nodes
	fl.Start()
}

func (fl *Flow)GetCurrentNode()*model.MetaNode {
	return fl.currentNode
}

func (fl *Flow) AddNode(node model.MetaNode) {
	fl.Nodes = append(fl.Nodes, node)
}

func (fl *Flow) IsNodeIdValid(currentNodeId model.NodeID, transitionNodeId model.NodeID) bool {
	if transitionNodeId == ""{
		return true
	}

	if currentNodeId == transitionNodeId {
		log.Error("Transition node can't be the same as current")
		return false
	}
	for i := range fl.Nodes {
		if fl.Nodes[i].Id == transitionNodeId {
			return true
		}
	}
	log.Error("<Flow> Transition node doesn't exist")
	return false
}

func (fl *Flow) Run() {
	var transitionNode model.NodeID
	defer func() {
		if r := recover(); r != nil {
			log.Error("<Flow> Flow process CRASHED with error : ",r)
			log.Errorf("<Flow> Crashed while processing message from Current Node = %d Next Node = %d ",fl.currentNodeId,transitionNode)
			transitionNode = ""
		}
	}()

	for {
		if !fl.isFlowRunning {
			break
		}
		for i := range fl.Nodes {
			if !fl.isFlowRunning {
				break
			}
			if fl.currentNodeId == "" && fl.Nodes[i].Type == "trigger" {
				log.Infof("<Flow> ------Flow %s is waiting for triggering event----------- ",fl.Name)
				var err error
				fl.currentMsg, fl.currentNode, err = node.TriggerNode(fl.Nodes, &fl.localContext, fl.msgInStream, fl.msgTransport, &fl.activeSubscriptions)
				if err != nil {
					log.Error("<Flow> TriggerNode failed with error :", err)
					fl.currentNodeId = ""
				}
				if !fl.isFlowRunning {
					break
				}
				fl.TriggerCounter++
				fl.currentNodeId = fl.currentNode.Id
				transitionNode = fl.currentNode.SuccessTransition
				if !fl.IsNodeIdValid(fl.currentNodeId,transitionNode) {
					log.Errorf("Unknown transition mode %s.Switching back to first node",transitionNode)
					transitionNode = ""
				}
				log.Debug("<Flow> Transition from Trigger to node = ",transitionNode)
			} else if fl.Nodes[i].Id == transitionNode {
				var err error
				switch fl.Nodes[i].Type {
				case "action":
					log.Info("<Flow> Executing ActionNode node.")
					err = node.ActionNode(&fl.Nodes[i], &fl.currentMsg, fl.msgTransport)
				case "wait":
					log.Info("<Flow> Executing WaitNode node.")
					err = node.WaitNode(&fl.Nodes[i])
				case "if":
					log.Info("<Flow> Executing IfNode node.")
					err = node.IfNode(&fl.Nodes[i], &fl.currentMsg)
				}
				fl.currentNodeId = fl.Nodes[i].Id
				fl.currentNode = &fl.Nodes[i]
				if err == nil {
					transitionNode = fl.Nodes[i].SuccessTransition
				} else {
					transitionNode = fl.Nodes[i].ErrorTransition
					fl.ErrorCounter++
					log.Errorf("<Flow> Node executed with error . Doing error transition to %s. Error : %s", transitionNode ,err)
				}
				if !fl.IsNodeIdValid(fl.currentNodeId,transitionNode) {
					log.Errorf("Unknown transition mode %s.Switching back to first node",transitionNode)
					transitionNode = ""
				}

			} else if transitionNode == "" {
				// Flow is finished . Returning to first step.
				fl.currentNodeId = ""
			}
		}

	}
	log.Infof("Flow was %s stopped.", fl.Name)

}

func (fl *Flow) Start() error {
	log.Info("<Flow> Starting flow : ", fl.Name)
	fl.isFlowRunning = true
	isFlowValid := false
	// The Flow should have at least one trigger or wait node to avoid tight loop
	for i := range fl.Nodes {
		if fl.Nodes[i].Type == "wait" || fl.Nodes[i].Type == "trigger" {
			isFlowValid = true
			break
		}
	}
	if isFlowValid{
		go fl.Run()
		log.Infof("<Flow> Flow %s is running", fl.Name)
		return nil
	}
	log.Errorf("<Flow> Flow %s is not valid and will not be started.Flow should have at least one trigger or wait node ",fl.Name)
	return errors.New("Flow should have at least one trigger or wait node")


}
func (fl *Flow) Stop() {
	log.Info("<Flow> Stopping flow  ", fl.Name)
	fl.isFlowRunning = false
	fl.msgInStream <- model.Message{}
}
func (fl *Flow) SetMessageStream(msgInStream model.MsgPipeline) {
	fl.msgInStream = msgInStream
}
