package node

import (
"github.com/alivinco/fimpgo"
"github.com/alivinco/fimpui/flow/model"
log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"time"
	"github.com/alivinco/fimpui/flow/utils"
)

type ReceiveNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	activeSubscriptions *[]string
	msgInStream model.MsgPipeline
	config ReceiveConfig

}

type ReceiveConfig struct {
	Timeout int64 // in seconds
	ValueFilter model.Variable
	IsValueFilterEnabled bool
}

func NewReceiveNode(flowOpCtx *model.FlowOperationalContext ,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := ReceiveNode{ctx:ctx,transport:transport}
	node.isStartNode = false
	node.isMsgReactor = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.config = ReceiveConfig{}
	return &node
}

func (node *ReceiveNode) ConfigureInStream(activeSubscriptions *[]string,msgInStream model.MsgPipeline) {
	node.activeSubscriptions = activeSubscriptions
	node.msgInStream = msgInStream
	node.initSubscriptions()
}

func (node *ReceiveNode) initSubscriptions() {
	log.Info(node.flowOpCtx.FlowId+"<Node> ReceiveNode is listening for events . Name = ", node.meta.Label)
	needToSubscribe := true
	for i := range *node.activeSubscriptions {
			if (*node.activeSubscriptions)[i] == node.meta.Address {
				needToSubscribe = false
				break
			}
	}
	if needToSubscribe {
			log.Info(node.flowOpCtx.FlowId+"<ReceiveNode> Subscribing for service by address :", node.meta.Address)
			node.transport.Subscribe(node.meta.Address)
			*node.activeSubscriptions = append(*node.activeSubscriptions, node.meta.Address)
	}
}


func (node *ReceiveNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		log.Error(err)
	}
	return err
}

func (node *ReceiveNode) WaitForEvent(nodeEventStream chan model.ReactorEvent) {
	node.isReactorRunning = true
	defer func() {
		node.isReactorRunning = false
		log.Debug("<ReceiveNode> Reactor-WaitForEvent is stopped ")
	}()
	log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Reactor-Waiting for event .chan size = ",len(node.msgInStream))
	start := time.Now()
	timeout := node.config.Timeout
	if timeout == 0 {
		timeout = 86400 // 24 hours
	}

	for {
		select {
		case newMsg := <-node.msgInStream:
			log.Info(node.flowOpCtx.FlowId+"<ReceiveNode> New message :")
			if newMsg.CancelOp {
				return
			}
			if utils.RouteIncludesTopic(node.meta.Address,newMsg.AddressStr) &&
				(newMsg.Payload.Service == node.meta.Service || node.meta.Service == "*") &&
				(newMsg.Payload.Type == node.meta.ServiceInterface || node.meta.ServiceInterface == "*") {
				if !node.config.IsValueFilterEnabled {
					newEvent := model.ReactorEvent{Msg:newMsg,TransitionNodeId:node.meta.SuccessTransition}
					select {
					case nodeEventStream <- newEvent:
						return
					default:
							log.Debug("<ReceiveNode> Message is dropped (no listeners) ")
					}

				} else if newMsg.Payload.Value == node.config.ValueFilter.Value {
					newEvent := model.ReactorEvent{Msg:newMsg,TransitionNodeId:node.meta.SuccessTransition}
					select {
					case nodeEventStream <- newEvent:
						return
					default:
						log.Debug("<ReceiveNode> Message is dropped (no listeners) ")
					}
				}
			}
			if node.config.Timeout > 0 {
				elapsed := time.Since(start)
				timeout = timeout - int64(elapsed.Seconds())
			}
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Not interested .")

		case <-time.After(time.Second * time.Duration(timeout)):
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Timeout ")
			newEvent := model.ReactorEvent{}
			newEvent.TransitionNodeId = node.meta.TimeoutTransition
			select {
			case nodeEventStream <- newEvent:
				return
			default:
				log.Debug("<ReceiveNode> Message is dropped (no listeners) ")
			}
		case signal := <-node.flowOpCtx.NodeControlSignalChannel:
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Control signal ")
			if signal == model.SIGNAL_STOP {
				return
			}
		}
	}
}

func (node *ReceiveNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	return nil,nil
}
