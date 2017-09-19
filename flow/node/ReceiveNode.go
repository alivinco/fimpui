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

func (node *ReceiveNode) SignalNodeIsReadyForNextMessage(){
	//select {
	//case node.flowOpCtx.NodeIsReady <- true:
	//default:
	//}
}

func (node *ReceiveNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Waiting for event .chan size = ",len(node.msgInStream))
	start := time.Now()
	timeout := node.config.Timeout
	if timeout == 0 {
		timeout = 86400 // 24 hours
	}
	for {
		select {
		case newMsg := <-node.msgInStream:
			log.Info(node.flowOpCtx.FlowId+"<ReceiveNode> New message :")
			if utils.RouteIncludesTopic(node.meta.Address,newMsg.AddressStr) &&
				(newMsg.Payload.Service == node.meta.Service || node.meta.Service == "*") &&
				(newMsg.Payload.Type == node.meta.ServiceInterface || node.meta.ServiceInterface == "*") {
				*msg = newMsg
				if !node.config.IsValueFilterEnabled {
					return []model.NodeID{node.meta.SuccessTransition}, nil
				} else if newMsg.Payload.Value == node.config.ValueFilter.Value {
					return []model.NodeID{node.meta.SuccessTransition}, nil
				}
			}
			if node.config.Timeout > 0 {
				elapsed := time.Since(start)
				timeout = timeout - int64(elapsed.Seconds())
			}
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Not interested .")

		case <-time.After(time.Second * time.Duration(timeout)):
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Timeout ")
			return []model.NodeID{node.meta.TimeoutTransition}, nil
		case signal := <-node.flowOpCtx.NodeControlSignalChannel:
			log.Debug(node.flowOpCtx.FlowId+"<ReceiveNode> Control signal ")
			if signal == model.SIGNAL_STOP {
				return nil,nil
			}
		}
	}

}
