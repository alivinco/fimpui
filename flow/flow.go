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
	FlowMeta            *model.FlowMeta
	globalContext       *model.Context
	opContext           model.FlowOperationalContext
	currentNodeIds      [] model.NodeID
	currentMsg          model.Message
	Nodes               []model.Node
	nodeInboundStreams  map[model.NodeID]model.MsgPipeline
	nodeOutboundStream  chan model.ReactorEvent
	msgTransport        *fimpgo.MqttTransport
	activeSubscriptions []string
	msgInStream         model.MsgPipeline
	TriggerCounter      int64
	ErrorCounter        int64
	StartedAt           time.Time
	WaitingSince        time.Time
	LastExecutionTime   time.Duration
}

func NewFlow(metaFlow model.FlowMeta, globalContext *model.Context, msgTransport *fimpgo.MqttTransport) *Flow {
	flow := Flow{globalContext: globalContext}
	flow.Nodes = make([]model.Node, 0)
	flow.currentNodeIds = make([]model.NodeID,1)
	flow.nodeInboundStreams = make(map[model.NodeID]model.MsgPipeline)
	flow.nodeOutboundStream = make(chan model.ReactorEvent)
	flow.msgTransport = msgTransport
	flow.globalContext = globalContext
	flow.opContext = model.FlowOperationalContext{NodeIsReady:make(chan bool),NodeControlSignalChannel:make(chan int)}
	flow.initFromMetaFlow(&metaFlow)
	return &flow
}

func (fl *Flow) SetStoragePath(path string) {
	fl.opContext.StoragePath = path
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
	defer func() {
		if r := recover(); r != nil {
			log.Error(fl.Id+"<Flow> Flow process CRASHED with error while doing node configuration : ",r)
			fl.opContext.State = "INIT_FAIL"
		}
	}()
	log.Infof("<Flow> ---------Initializing Flow Id = %s , Name = %s -----------",fl.Id,fl.Name)
	for _,metaNode := range fl.FlowMeta.Nodes {
		var newNode model.Node
		log.Infof("<Flow> Loading node . Type = %s , Label = %s",metaNode.Type,metaNode.Label)
		constructor ,ok := node.Registry[metaNode.Type]
		if ok {
			newNode = constructor(&fl.opContext,metaNode,fl.globalContext,fl.msgTransport)
			if newNode.IsMsgReactorNode() {
				log.Debug("<Flow> Creating a channel for Reactor-Node id = ",metaNode.Id)
				// Each reactor node gets its own channel.
				fl.nodeInboundStreams[metaNode.Id] = make(model.MsgPipeline)
				newNode.ConfigureInStream(&fl.activeSubscriptions,fl.nodeInboundStreams[metaNode.Id])
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
	fl.StartMsgStreamRouter()
}

func (fl*Flow) GetContext()*model.Context {
	return fl.globalContext
}

func (fl*Flow) GetCurrentMessage()*model.Message {
	return &fl.currentMsg
}

func (fl *Flow) SetNodes(nodes []model.Node) {
	fl.Nodes = nodes
}

func (fl *Flow) ReloadNodes(nodes []model.Node) {
	fl.Stop()
	fl.Nodes = nodes
	fl.Start()
}

func (fl *Flow) GetNodeById(id model.NodeID) model.Node {
	for i := range fl.Nodes {
		if fl.Nodes[i].GetMetaNode().Id == id {
			return fl.Nodes[i]
		}
	}
	return nil
}

func (fl *Flow)GetFlowStats() (*model.FlowStatsReport) {
	stats := model.FlowStatsReport{}
	currentNode := fl.GetNodeById(fl.currentNodeIds[0])
	if currentNode != nil {
		stats.CurrentNodeId = currentNode.GetMetaNode().Id
		stats.CurrentNodeLabel = currentNode.GetMetaNode().Label
		stats.IsAtStartingPoint = currentNode.IsStartNode()
	}
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

func (fl *Flow) IsFlowValid() bool {
	var flowHasStartNode bool
	for i := range fl.Nodes {
		node := fl.Nodes[i].GetMetaNode()
		if node.Type == "trigger" || node.Type == "action" || node.Type == "receive" {
			if node.Address == "" ||  node.ServiceInterface == "" || node.Service == ""	{
				log.Error(fl.Id+"<Flow> Flow is not valid , node is not configured . Node ",node.Label)
				return false
			}
		}
		if fl.Nodes[i].IsStartNode() {
			flowHasStartNode = true
		}
	}
	if !flowHasStartNode {
		log.Error(fl.Id+"<Flow> Flow is not valid, start node not found")
		return false
	}
	return true
}

func (fl *Flow) Run() {
	var transitionNodeId model.NodeID
	defer func() {
		if r := recover(); r != nil {
			log.Error(fl.Id+"<Flow> Flow process CRASHED with error : ",r)
			log.Errorf(fl.Id+"<Flow> Crashed while processing message from Current Node = %d Next Node = %d ",fl.currentNodeIds[0], transitionNodeId)
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
			if fl.currentNodeIds[0] == "" && fl.Nodes[i].IsStartNode() {
				fl.LastExecutionTime = time.Since(fl.StartedAt)
				log.Infof(fl.Id+"<Flow> ------Flow %s is waiting for triggering event----------- ",fl.Name)
				// Initial message received by trigger , which is passed further throughout the flow.
				fl.WaitingSince = time.Now()
				// Starting all Start nodes and waiting for event from one of them
				fl.currentNodeIds = fl.currentNodeIds[:0]
				for si := range fl.Nodes {
					if fl.Nodes[si].IsStartNode() {
						if ! fl.Nodes[si].IsReactorRunning(){
							go fl.Nodes[si].WaitForEvent(fl.nodeOutboundStream)
						}
						fl.currentNodeIds = append(fl.currentNodeIds,fl.Nodes[si].GetMetaNode().Id )
					}

				}
				// Blocking wait
				reactorEvent :=<- fl.nodeOutboundStream
				log.Debug(fl.Id+"<Flow> New event from reactor node.")
				fl.StartedAt = time.Now()
				fl.currentMsg = reactorEvent.Msg
				//log.Debug("<Flow> msg.payload : ",fl.currentMsg)
				if reactorEvent.Err != nil {
					log.Error(fl.Id+"<Flow> TriggerNode failed with error :", reactorEvent.Err)
					fl.currentNodeIds[0] = ""
				}
				if !fl.opContext.IsFlowRunning {
					break
				}
				fl.TriggerCounter++
				//fl.currentNodeId = fl.Nodes[i].GetMetaNode().Id
				transitionNodeId = reactorEvent.TransitionNodeId
				log.Debug(fl.Id+"<Flow> Transition node id = ",transitionNodeId)
				log.Debug(fl.Id+"<Flow> Current nodes = ",fl.currentNodeIds)
				if !fl.IsNodeIdValid(fl.currentNodeIds[0], transitionNodeId) {
					log.Errorf(fl.Id+"<Flow> Unknown transition node %s from first node.Switching back to first node", transitionNodeId)
					transitionNodeId = ""
				}
				log.Debug("<Flow> Transition from Trigger to node = ", transitionNodeId)
			} else if fl.Nodes[i].GetMetaNode().Id == transitionNodeId {
				var err error
				var nextNodes []model.NodeID
				fl.currentNodeIds[0] = fl.Nodes[i].GetMetaNode().Id
				if fl.Nodes[i].IsMsgReactorNode() {
					if ! fl.Nodes[i].IsReactorRunning(){
						go fl.Nodes[i].WaitForEvent(fl.nodeOutboundStream)
					}
					// Blocking wait
					reactorEvent :=<- fl.nodeOutboundStream
					fl.currentMsg = reactorEvent.Msg
					transitionNodeId = reactorEvent.TransitionNodeId
					err = reactorEvent.Err
					log.Debug(fl.Id+"<Flow> New event from reactor node.")
					//log.Debug("<Flow> msg.payload : ",fl.currentMsg)

				}else {
					nextNodes, err = fl.Nodes[i].OnInput(&fl.currentMsg)
					if len(nextNodes)>0 {
						transitionNodeId = nextNodes[0]
					}else {
						transitionNodeId = ""
					}
				}

				if err != nil {
					fl.ErrorCounter++
					log.Errorf(fl.Id+"<Flow> Node executed with error . Doing error transition to %s. Error : %s", transitionNodeId,err)
				}

				if !fl.IsNodeIdValid(fl.currentNodeIds[0], transitionNodeId) {
					log.Errorf(fl.Id+"<FLow> Unknown transition node %s.Switching back to first node", transitionNodeId)
					transitionNodeId = ""
				}
				log.Debug(fl.Id+"<FLow> Transition to node = ",transitionNodeId)

			} else if transitionNodeId == "" {
				// Flow is finished . Returning to first step.
				fl.currentNodeIds[0] = ""
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
	isFlowValid := fl.IsFlowValid()

	if isFlowValid{
		// Init all nodes
		for i := range fl.Nodes{
			fl.Nodes[i].Init()
		}
		fl.opContext.State = "RUNNING"
		log.Infof(fl.Id+"<Flow> Flow %s is running", fl.Name)
		go fl.Run()
	}else {
		fl.opContext.State = "NOT_CONFIGURED"
		log.Errorf(fl.Id+"<Flow> Flow %s is not valid and will not be started.Flow should have at least one trigger or wait node ",fl.Name)
		return errors.New("Flow should have at least one trigger or wait node")
	}
	return nil
}

func (fl *Flow) CancelAllRunningNodes() {
	for _,stream := range fl.nodeInboundStreams {
		cancelMsg := model.Message{CancelOp:true}
		select {
		case stream <- cancelMsg:
			continue
		default:
			continue
		}
		close(stream)
	}
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
	fl.CancelAllRunningNodes()
	close(fl.nodeOutboundStream)
	log.Info(fl.Id+"<Flow> All streams and running goroutins were closed  ")
	log.Info(fl.Id+"<Flow> Stopped .  ", fl.Name)
	return nil
}

func (fl *Flow) GetFlowState() string {
	return fl.opContext.State
}

func (fl *Flow) IsNodeCurrentNode(nodeId model.NodeID) bool {
	for i := range fl.currentNodeIds {
		if fl.currentNodeIds[i] == nodeId {
			return true
		}
	}
	return false
}

func (fl *Flow) StartMsgStreamRouter() {
	// Message broadcast from flow incomming stream to reactor nodes
	go func() {
		log.Info(fl.Id+"<Flow> Starting flow msg router ")
		defer func() {
			log.Info(fl.Id+"<Flow> Router is stopped ")
		}()
		for msg := range fl.msgInStream {
			for nodeId,stream := range fl.nodeInboundStreams {
				if fl.IsNodeCurrentNode(nodeId){
					select {
					case stream <- msg:
						log.Debug("<Flow> Message was sent to node with Id = ", nodeId)
					default:
						log.Debug("<Flow> Message is dropped (no listeners) nodeId = ", nodeId)
					}
				}
			}
			if msg.CancelOp {
				return
			}

		}
	}()

}

func (fl *Flow) SetMessageStream(msgInStream model.MsgPipeline) {
	fl.msgInStream = msgInStream
}
