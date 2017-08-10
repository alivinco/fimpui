package node

import (
"github.com/alivinco/fimpgo"
"github.com/alivinco/fimpui/flow/model"
log "github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"time"
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
	Timeout int64
	ValueFilter model.Variable
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
	log.Info("<Node> ReceiveNode is listening for events . Name = ", node.meta.Label)
	needToSubscribe := true
	for i := range *node.activeSubscriptions {
			if (*node.activeSubscriptions)[i] == node.meta.Address {
				needToSubscribe = false
				break
			}
	}
	if needToSubscribe {
			log.Info("<ReceiveNode> Subscribing for service by address :", node.meta.Address)
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

func (node *ReceiveNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Debug("<ReceiveNode> Waiting for event ")
	start := time.Now()
	timeout := node.config.Timeout
	for {
		select {
		case newMsg := <-node.msgInStream:
			log.Info("<ReceiveNode> New message :")
			if node.config.ValueFilter.ValueType == "" {
				return []model.NodeID{node.meta.SuccessTransition}, nil
			}else if newMsg.Payload.Value == node.config.ValueFilter.Value {
				return []model.NodeID{node.meta.SuccessTransition}, nil
			}else {
				elapsed := time.Since(start)
				timeout =  timeout - int64(elapsed.Seconds())
			}
		case <-time.After(time.Second * time.Duration(timeout)):
			log.Debug("<ReceiveNode> Timeout ")
			return []model.NodeID{node.meta.TimeoutTransition}, nil
		case signal := <-node.flowOpCtx.NodeControlSignalChannel:
			log.Debug("<ReceiveNode> Control signal ")
			if signal == model.SIGNAL_STOP {
				return nil,nil
			}
		}
	}

}

