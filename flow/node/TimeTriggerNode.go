package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/robfig/cron"
)

type TimeTriggerNode struct {
	BaseNode
	ctx                 *model.Context
	config              TimeTriggerConfig
	cron 	*cron.Cron
}

type TimeTriggerConfig struct {
	ValueFilter model.Variable
}

func NewTimeTriggerNode(flowOpCtx *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node {
	node := TimeTriggerNode{ctx: ctx}
	node.isStartNode = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.config = TimeTriggerConfig{}
	node.cron = cron.New()
	return &node
}

func (node *TimeTriggerNode) LoadNodeConfig() error {
	
	return nil
}

// is invoked when node is started
func (node *TimeTriggerNode) Init() error {
	node.cron.Start()
	return nil
}

// is invoked when node flow is stopped
func (node *TimeTriggerNode) Cleanup() error {
	node.cron.Stop()
	return nil
}

func (node *TimeTriggerNode) OnInput(msg *model.Message) ([]model.NodeID, error) {
	log.Debug("<TrigNode> Waiting for event ")
	return []model.NodeID{node.meta.SuccessTransition}, nil
}
