package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/alivinco/fimpui/flow/node"
)

type Flow struct {
	Id                  string
	Name                string
	Description         string
	FlowMeta 			*model.FlowMeta
	globalContext       *model.Context
	opContext           model.FlowOperationalContext
	currentNodeId       model.NodeID
	currentMsg          *model.Message
	currentNode         model.Node
	Nodes               []model.Node
	msgPipeline         model.MsgPipeline
	msgTransport        *fimpgo.MqttTransport
	activeSubscriptions []string
	msgInStream         model.MsgPipeline
	localMsgInStream    map[model.NodeID]model.MsgPipeline
	TriggerCounter      int64
	ErrorCounter        int64
}

func NewFlow(metaFlow model.FlowMeta, globalContext *model.Context, msgTransport *fimpgo.MqttTransport) *Flow {
	flow := Flow{globalContext: globalContext}
	flow.msgPipeline = make(model.MsgPipeline)
	flow.Nodes = make([]model.Node, 0)
	flow.msgTransport = msgTransport
	flow.globalContext = globalContext
	flow.opContext = model.FlowOperationalContext{}
	flow.initFromMetaFlow(&metaFlow)
	flow.globalContext.RegisterFlow(flow.Id)
	return &flow
}

func (fl *Flow) CleanupBeforeDelete() {
	fl.globalContext.UnregisterFlow(fl.Id)
}

func (fl *Flow) initFromMetaFlow(meta *model.FlowMeta) {
	fl.Id = meta.Id
	fl.Name = meta.Name
	fl.Description = meta.Description
	fl.FlowMeta = meta
	fl.opContext.FlowId = meta.Id
	fl.localMsgInStream = make(map[model.NodeID]model.MsgPipeline)
}

func (fl *Flow) InitAllNodes() {
	log.Infof("<Flow> ---------Initializing Flow Id = %s , Name = %s -----------",fl.Id,fl.Name)
	for _,metaNode := range fl.FlowMeta.Nodes {
		var newNode model.Node
		log.Infof("<Flow> Loading node . Type = %s , Label = %s",metaNode.Type,metaNode.Label)
		constructor ,ok := node.Registry[metaNode.Type]
		if ok {
			newNode = constructor(&fl.opContext,metaNode,fl.globalContext,fl.msgTransport)
			if newNode.IsMsgReactorNode() {
				// Creating channel for each message reactor
				nodeChannel := make(model.MsgPipeline)
				fl.localMsgInStream[metaNode.Id] = nodeChannel
				// Configuring input message stream for nodes like trigger , receive , etc.
				newNode.ConfigureInStream(&fl.activeSubscriptions,nodeChannel)
			}
		}else {
			log.Errorf("<Flow> Node type = %s isn't supported",metaNode.Type)
		}
		err := newNode.LoadNodeConfig()
		if err == nil {
			fl.AddNode(newNode)
			log.Info("<Flow> Node is loaded.")
		}else {
			log.Errorf("<Flow> Node type %s can't be loaded . Error : %s",metaNode.Type,err)
		}
	}
}

func (fl*Flow) GetContext()*model.Context {
	return fl.globalContext
}

func (fl *Flow) SetNodes(nodes []model.Node) {
	fl.Nodes = nodes
}

func (fl *Flow) ReloadNodes(nodes []model.Node) {
	fl.Stop()
	fl.Nodes = nodes
	fl.Start()
}

func (fl *Flow)GetCurrentNode()model.Node {
	return fl.currentNode
}

func (fl *Flow) AddNode(node model.Node) {
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
		if fl.Nodes[i].GetMetaNode().Id == transitionNodeId {
			return true
		}
	}
	log.Error("<Flow> Transition node doesn't exist")
	return false
}

func (fl *Flow) Run() {
	var transitionNodeId model.NodeID
	defer func() {
		if r := recover(); r != nil {
			log.Error("<Flow> Flow process CRASHED with error : ",r)
			log.Errorf("<Flow> Crashed while processing message from Current Node = %d Next Node = %d ",fl.currentNodeId, transitionNodeId)
			transitionNodeId = ""
		}
	}()

	for {
		if !fl.opContext.IsFlowRunning {
			break
		}
		for i := range fl.Nodes {
			if !fl.opContext.IsFlowRunning {
				break
			}
			if fl.currentNodeId == "" && fl.Nodes[i].IsStartNode() {
				log.Infof("<Flow> ------Flow %s is waiting for triggering event----------- ",fl.Name)
				fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				fl.currentNode = fl.Nodes[i]
				var err error
				newMsg := model.Message{}
				nextNodes, err := fl.Nodes[i].OnInput (&newMsg)
				fl.currentMsg = &newMsg
				if err != nil {
					log.Error("<Flow> TriggerNode failed with error :", err)
					fl.currentNodeId = ""
				}
				if !fl.opContext.IsFlowRunning {
					break
				}
				fl.TriggerCounter++
				fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				transitionNodeId = nextNodes[0]
				if !fl.IsNodeIdValid(fl.currentNodeId, transitionNodeId) {
					log.Errorf("Unknown transition mode %s.Switching back to first node", transitionNodeId)
					transitionNodeId = ""
				}
				log.Debug("<Flow> Transition from Trigger to node = ", transitionNodeId)
			} else if fl.Nodes[i].GetMetaNode().Id == transitionNodeId {
				var err error
				nextNodes, err := fl.Nodes[i].OnInput(fl.currentMsg)
				fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				fl.currentNode = fl.Nodes[i]
				if err != nil {
					fl.ErrorCounter++
					log.Errorf("<Flow> Node executed with error . Doing error transition to %s. Error : %s", transitionNodeId,err)
				}
				if len(nextNodes)>0 {
					transitionNodeId = nextNodes[0]
				}else {
					transitionNodeId = ""
				}
				if !fl.IsNodeIdValid(fl.currentNodeId, transitionNodeId) {
					log.Errorf("Unknown transition mode %s.Switching back to first node", transitionNodeId)
					transitionNodeId = ""
				}

			} else if transitionNodeId == "" {
				// Flow is finished . Returning to first step.
				fl.currentNodeId = ""
			}
		}

	}
	fl.opContext.State = "STOPPED"
	log.Infof("Flow was %s stopped.", fl.Name)

}

func (fl *Flow) InStreamMsgRouter() {
	// fetching all messages
	for inMsg := range fl.msgInStream {
		//log.Debug("<Flow> Router : New message from msgInStream")
		if !fl.opContext.IsFlowRunning {
			break
		}
		currNode := fl.currentNode.GetMetaNode()
		// doing filtering
		if (inMsg.AddressStr == currNode.Address || currNode.Address == "*") &&
			(inMsg.Payload.Service == currNode.Service || currNode.Service == "*") &&
			(inMsg.Payload.Type == currNode.ServiceInterface || currNode.ServiceInterface == "*") {
			// sending message to each channel
			select {
				case fl.localMsgInStream[currNode.Id] <- inMsg:
					log.Debug("<Flow> Router: Message was sent to Node Id = ", fl.currentNodeId)
				default:
					log.Debug("<Flow> Router: Message is dropped (no listeners).")
			}
		}
	}
}

// Starts Flow loop in its own goroutine and sets isFlowRunning flag to true
func (fl *Flow) Start() error {
	log.Info("<Flow> Starting flow : ", fl.Name)
	fl.opContext.State = "STARTING"
	fl.opContext.IsFlowRunning = true
	isFlowValid := false
	// Starting flow loop for every trigger.
	for i := range fl.Nodes {
		if fl.Nodes[i].IsStartNode() {
			go fl.Run()
			go fl.InStreamMsgRouter()
			isFlowValid = true
			fl.opContext.State = "RUNNING"
			log.Infof("<Flow> Flow %s is running", fl.Name)
		}
	}
	if !isFlowValid{
		fl.opContext.State = "STOPPED"
		log.Errorf("<Flow> Flow %s is not valid and will not be started.Flow should have at least one trigger or wait node ",fl.Name)
		return errors.New("Flow should have at least one trigger or wait node")
	}
	return nil
}
// Terminates flow loop , stops goroutine .
func (fl *Flow) Stop() {
	log.Info("<Flow> Stopping flow  ", fl.Name)
	fl.opContext.IsFlowRunning = false
	fl.msgInStream <- model.Message{}
}

func (fl *Flow) SetMessageStream(msgInStream model.MsgPipeline) {
	fl.msgInStream = msgInStream
}
