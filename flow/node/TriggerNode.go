package node

import (
	log "github.com/Sirupsen/logrus"
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
	IsValueFilterEnabled bool
}

func NewTriggerNode(flowOpCtx *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node {
	node := TriggerNode{ctx: ctx, transport: transport}
	node.isStartNode = true
	node.isMsgReactor = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.config = TriggerConfig{}
	return &node
}

func (node *TriggerNode) ConfigureInStream(activeSubscriptions *[]string, msgInStream model.MsgPipeline) {
	log.Info(node.flowOpCtx.FlowId+"<TrigNode>Configuring Stream")
	node.activeSubscriptions = activeSubscriptions
	node.msgInStream = msgInStream
	node.initSubscriptions()
}

func (node *TriggerNode) initSubscriptions() {
	if node.meta.Type == "trigger" {
		log.Info(node.flowOpCtx.FlowId+"<TrigNode> TriggerNode is listening for events . Name = ", node.meta.Label)
		needToSubscribe := true
		for i := range *node.activeSubscriptions {
			if (*node.activeSubscriptions)[i] == node.meta.Address {
				needToSubscribe = false
				break
			}
		}
		if needToSubscribe {
			log.Info(node.flowOpCtx.FlowId+"<TrigNode> Subscribing for service by address :", node.meta.Address)
			node.transport.Subscribe(node.meta.Address)
			*node.activeSubscriptions = append(*node.activeSubscriptions, node.meta.Address)
		}
	}
}

func (node *TriggerNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		log.Error(err)
	}
	return err
}


func (node *TriggerNode) OnInput(msg *model.Message) ([]model.NodeID, error) {
	log.Debug(node.flowOpCtx.FlowId+"<TrigNode> Waiting for event . chan size = ",len(node.msgInStream))
	start := time.Now()
	timeout := node.config.Timeout
	if timeout == 0 {
		timeout = 86400 // 24 hours
	}
	for {
		select {
		case newMsg := <-node.msgInStream:
			log.Debug(node.flowOpCtx.FlowId+"<TrigNode> New message from InStream ")
			if utils.RouteIncludesTopic(node.meta.Address,newMsg.AddressStr) &&
			   (newMsg.Payload.Service == node.meta.Service || node.meta.Service == "*") &&
			   (newMsg.Payload.Type == node.meta.ServiceInterface || node.meta.ServiceInterface == "*") {

				*msg = newMsg
				if !node.config.IsValueFilterEnabled {
					return []model.NodeID{node.meta.SuccessTransition}, nil
				}else if newMsg.Payload.Value == node.config.ValueFilter.Value {
					return []model.NodeID{node.meta.SuccessTransition}, nil
				}
			}
			if node.config.Timeout > 0 {
				elapsed := time.Since(start)
				timeout =  timeout - int64(elapsed.Seconds())
			}
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Not interested .")

		case <-time.After(time.Second * time.Duration(timeout)):
			log.Debug(node.flowOpCtx.FlowId+"<TrigNode> Timeout ")
			return []model.NodeID{node.meta.TimeoutTransition}, nil
		case signal := <-node.flowOpCtx.NodeControlSignalChannel:
			log.Debug(node.flowOpCtx.FlowId+"<TrigNode> Control signal ")
			if signal == model.SIGNAL_STOP {
				return nil, nil
			}
		}
	}
}
