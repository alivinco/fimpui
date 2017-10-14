package flow

import (
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/alivinco/fimpui/flow/node"
	"time"
	"github.com/alivinco/fimpui/flow/utils"
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
	TriggerCounter      int64
	ErrorCounter        int64
	StartedAt           time.Time
	WaitingSince 		time.Time
	LastExecutionTime   time.Duration
}

func NewFlow(metaFlow model.FlowMeta, globalContext *model.Context, msgTransport *fimpgo.MqttTransport) *Flow {
	flow := Flow{globalContext: globalContext}
	flow.msgPipeline = make(model.MsgPipeline)
	flow.Nodes = make([]model.Node, 0)
	flow.msgTransport = msgTransport
	flow.globalContext = globalContext
	flow.opContext = model.FlowOperationalContext{NodeIsReady:make(chan bool),NodeControlSignalChannel:make(chan int)}
	flow.initFromMetaFlow(&metaFlow)
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
	//fl.localMsgInStream = make(map[model.NodeID]model.MsgPipeline,10)
	fl.globalContext.RegisterFlow(fl.Id)
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
				log.Debug("<Flow> Creating a channel for node id = ",metaNode.Id)
				newNode.ConfigureInStream(&fl.activeSubscriptions,fl.msgInStream)
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


func (fl *Flow)GetFlowStats() (*model.FlowStatsReport) {
	stats := model.FlowStatsReport{}
	stats.CurrentNodeId = fl.currentNode.GetMetaNode().Id
	stats.CurrentNodeLabel = fl.currentNode.GetMetaNode().Label
	stats.IsAtStartingPoint = fl.currentNode.IsStartNode()
	stats.StartedAt = fl.StartedAt
	stats.WaitingSince = fl.WaitingSince
	stats.LastExecutionTime = int64(fl.LastExecutionTime/time.Millisecond)
	return &stats
}


func (fl *Flow) AddNode(node model.Node) {
	fl.Nodes = append(fl.Nodes, node)
}

func (fl *Flow) IsNodeIdValid(currentNodeId model.NodeID, transitionNodeId model.NodeID) bool {
	if transitionNodeId == ""{
		return true
	}

	if currentNodeId == transitionNodeId {
		log.Error(fl.Id+"<Flow> Transition node can't be the same as current")
		return false
	}
	for i := range fl.Nodes {
		if fl.Nodes[i].GetMetaNode().Id == transitionNodeId {
			return true
		}
	}
	log.Error(fl.Id+"<Flow> Transition node doesn't exist")
	return false
}


func (fl *Flow) Run() {
	var transitionNodeId model.NodeID
	defer func() {
		if r := recover(); r != nil {
			log.Error(fl.Id+"<Flow> Flow process CRASHED with error : ",r)
			log.Errorf(fl.Id+"<Flow> Crashed while processing message from Current Node = %d Next Node = %d ",fl.currentNodeId, transitionNodeId)
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
				var err error
				log.Infof(fl.Id+"<Flow> ------Flow %s is waiting for triggering event----------- ",fl.Name)
				fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				fl.currentNode = fl.Nodes[i]
				newMsg := model.Message{}
				fl.WaitingSince = time.Now()
				fl.LastExecutionTime = time.Since(fl.StartedAt)
				nextNodes, err := fl.Nodes[i].OnInput (&newMsg)
				fl.StartedAt = time.Now()
				fl.currentMsg = &newMsg
				if err != nil {
					log.Error(fl.Id+"<Flow> TriggerNode failed with error :", err)
					fl.currentNodeId = ""
				}
				if !fl.opContext.IsFlowRunning {
					break
				}
				fl.TriggerCounter++
				//fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				transitionNodeId = nextNodes[0]
				if !fl.IsNodeIdValid(fl.currentNodeId, transitionNodeId) {
					log.Errorf(fl.Id+"<Flow> Unknown transition mode %s.Switching back to first node", transitionNodeId)
					transitionNodeId = ""
				}
				log.Debug("<Flow> Transition from Trigger to node = ", transitionNodeId)
			} else if fl.Nodes[i].GetMetaNode().Id == transitionNodeId {
				var err error
				fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				fl.currentNode = fl.Nodes[i]
				nextNodes, err := fl.Nodes[i].OnInput(fl.currentMsg)
				if err != nil {
					fl.ErrorCounter++
					log.Errorf(fl.Id+"<Flow> Node executed with error . Doing error transition to %s. Error : %s", transitionNodeId,err)
				}
				if len(nextNodes)>0 {
					transitionNodeId = nextNodes[0]
				}else {
					transitionNodeId = ""
				}
				if !fl.IsNodeIdValid(fl.currentNodeId, transitionNodeId) {
					log.Errorf(fl.Id+"<FLow> Unknown transition mode %s.Switching back to first node", transitionNodeId)
					transitionNodeId = ""
				}
				log.Debug(fl.Id+"<FLow> Transition to node = ",transitionNodeId)

			} else if transitionNodeId == "" {
				// Flow is finished . Returning to first step.
				fl.currentNodeId = ""
			}
		}

	}
	fl.opContext.State = "STOPPED"
	log.Infof(fl.Id+"<Flow> Runner for flow %s stopped.", fl.Name)

}

func (fl *Flow) IsFlowInterestedInMessage(topic string ) bool {
	for i :=range fl.activeSubscriptions {
		if utils.RouteIncludesTopic(fl.activeSubscriptions[i],topic) {
			return true
		}else {
			//log.Debug(fl.Id+"<Flow> Not interested in topic : ",topic)
		}
	}
	return false
}

// Starts Flow loop in its own goroutine and sets isFlowRunning flag to true
func (fl *Flow) Start() error {
	log.Info(fl.Id+"<Flow> Starting flow : ", fl.Name)
	fl.opContext.State = "STARTING"
	fl.opContext.IsFlowRunning = true
	isFlowValid := false
	// Init all nodes
	for i := range fl.Nodes{
		fl.Nodes[i].Init()
	}

	// Validating flow is it has at least one start node .
	for i := range fl.Nodes {
		if fl.Nodes[i].IsStartNode() {
			isFlowValid = true
			fl.opContext.State = "RUNNING"
			log.Infof(fl.Id+"<Flow> Flow %s is running", fl.Name)
		}
	}

	if isFlowValid{
		go fl.Run()
	}else {
		fl.opContext.State = "STOPPED"
		log.Errorf(fl.Id+"<Flow> Flow %s is not valid and will not be started.Flow should have at least one trigger or wait node ",fl.Name)
		return errors.New("Flow should have at least one trigger or wait node")
	}
	return nil
}
// Terminates flow loop , stops goroutine .
func (fl *Flow) Stop() error {
	log.Info(fl.Id+"<Flow> Stopping flow  ", fl.Name)

	// is invoked when node flow is stopped
	for _,topic := range fl.activeSubscriptions {
		log.Info(fl.Id+"<Flow> Unsubscribing from topic : ",topic)
		fl.msgTransport.Unsubscribe(topic)
	}

	fl.opContext.IsFlowRunning = false
	select {
	case fl.opContext.NodeControlSignalChannel <- model.SIGNAL_STOP:
	default:
		log.Debug(fl.Id+"<Flow> No signal listener.")
	}
	fl.msgInStream <- model.Message{}
	for i := range fl.Nodes{
		fl.Nodes[i].Cleanup()
	}
	log.Info(fl.Id+"<Flow> Stopped .  ", fl.Name)
	return nil
}

func (fl *Flow) GetFlowState() string {
	return fl.opContext.State
}

func (fl *Flow) SetMessageStream(msgInStream model.MsgPipeline) {
	fl.msgInStream = msgInStream
}
