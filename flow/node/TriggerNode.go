package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
	"time"
	"github.com/alivinco/fimpui/flow/utils"
)

type TriggerNode struct {
	BaseNode
	ctx                 *model.Context
	transport           *fimpgo.MqttTransport
	activeSubscriptions *[]string
	msgInStream         model.MsgPipeline
	config              TriggerConfig
}

type TriggerConfig struct {
	Timeout int64 // in seconds
	ValueFilter model.Variable
	InputVariableType string
	IsValueFilterEnabled bool
	RegisterAsVirtualService bool
	VirtualServiceGroup string
}

func NewTriggerNode(flowOpCtx *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node {
	node := TriggerNode{ctx: ctx, transport: transport}
	node.isStartNode = true
	node.isMsgReactor = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.config = TriggerConfig{}
	node.SetupBaseNode()
	return &node
}

func (node *TriggerNode) ConfigureInStream(activeSubscriptions *[]string, msgInStream model.MsgPipeline) {
	node.getLog().Info("Configuring Stream")
	node.activeSubscriptions = activeSubscriptions
	node.msgInStream = msgInStream
	node.initSubscriptions()
}

func (node *TriggerNode) initSubscriptions() {
	if node.meta.Type == "trigger" {
		node.getLog().Info("TriggerNode is listening for events . Name = ", node.meta.Label)
		needToSubscribe := true
		for i := range *node.activeSubscriptions {
			if (*node.activeSubscriptions)[i] == node.meta.Address {
				needToSubscribe = false
				break
			}
		}
		if needToSubscribe {
			if node.meta.Address != ""{
				node.getLog().Info("Subscribing for service by address :", node.meta.Address)
				node.transport.Subscribe(node.meta.Address)
				*node.activeSubscriptions = append(*node.activeSubscriptions, node.meta.Address)
			}else {
				node.getLog().Error(" Can't subscribe to service with empty address")
			}

		}
	}
}

func (node *TriggerNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		node.getLog().Error("Error while decoding node configs.Err:",err)
	}
	return err
}


func (node *TriggerNode) WaitForEvent(nodeEventStream chan model.ReactorEvent) {
	node.isReactorRunning = true
	defer func() {
		node.isReactorRunning = false
		node.getLog().Debug(" WaitForEvent is stopped ")
	}()
	node.getLog().Debug( "Waiting for event . chan size = ",len(node.msgInStream))
	start := time.Now()
	timeout := node.config.Timeout
	if timeout == 0 {
		timeout = 86400 // 24 hours
	}
	for {
		select {
		case newMsg := <-node.msgInStream:
			if newMsg.CancelOp {
				return
			}
			node.getLog().Debug("New message from InStream ")
			if utils.RouteIncludesTopic(node.meta.Address,newMsg.AddressStr) &&
				(newMsg.Payload.Service == node.meta.Service || node.meta.Service == "*") &&
				(newMsg.Payload.Type == node.meta.ServiceInterface || node.meta.ServiceInterface == "*") {

				if !node.config.IsValueFilterEnabled {
					newEvent := model.ReactorEvent{Msg:newMsg,TransitionNodeId:node.meta.SuccessTransition}
					select {
					case nodeEventStream <- newEvent:
						return
					default:
						node.getLog().Debug("Message is dropped (no listeners) ")
					}
				}else if newMsg.Payload.Value == node.config.ValueFilter.Value {
					newEvent := model.ReactorEvent{Msg:newMsg,TransitionNodeId:node.meta.SuccessTransition}
					select {
					case nodeEventStream <- newEvent:
						return
					default:
						node.getLog().Debug("Message is dropped (no listeners) ")
					}
				}
			}
			if node.config.Timeout > 0 {
				elapsed := time.Since(start)
				timeout =  timeout - int64(elapsed.Seconds())
			}
			node.getLog().Debug("Not interested .")

		case <-time.After(time.Second * time.Duration(timeout)):
			node.getLog().Debug("Timeout ")
			newEvent := model.ReactorEvent{TransitionNodeId:node.meta.TimeoutTransition}
			select {
			case nodeEventStream <- newEvent:
				return
			default:
				node.getLog().Debug("Message is dropped (no listeners) ")
			}
		case signal := <-node.flowOpCtx.NodeControlSignalChannel:
			node.getLog().Debug("Control signal ")
			if signal == model.SIGNAL_STOP {
				return
			}
		}
	}
}

func (node *TriggerNode) OnInput(msg *model.Message) ([]model.NodeID, error) {
	return nil,nil
}
