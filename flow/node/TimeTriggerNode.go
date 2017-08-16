package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/robfig/cron"
	"github.com/mitchellh/mapstructure"
	//"github.com/cpucycle/astrotime"
	//"time"
)

type TimeTriggerNode struct {
	BaseNode
	ctx                 *model.Context
	config              TimeTriggerConfig
	cron 	*cron.Cron
	cronMessageCh model.MsgPipeline
}

type TimeTriggerConfig struct {
	DefaultMsg model.Variable
	Expressions []TimeExpression
	GenerateAstroTimeEvents bool
	Latitude float64
	Longitude float64
}

type TimeExpression struct {
	Name string
	Expression string   //https://godoc.org/github.com/robfig/cron#Job
	Comment string
}

func NewTimeTriggerNode(flowOpCtx *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node {
	node := TimeTriggerNode{ctx: ctx}
	node.isStartNode = true
	node.flowOpCtx = flowOpCtx
	node.meta = meta
	node.config = TimeTriggerConfig{}
	node.cron = cron.New()
	node.cronMessageCh = make(model.MsgPipeline)
	return &node
}

func (node *TimeTriggerNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		log.Error(err)
	}
	return err
}

// is invoked when node is started
func (node *TimeTriggerNode) Init() error {
	if node.config.GenerateAstroTimeEvents {
		//t := astrotime.NextSunrise(time.Now(), node.config.Latitude, node.config.Longitude)


	}else {
		for i := range node.config.Expressions {
			node.cron.AddFunc(node.config.Expressions[i].Expression,func() {
				log.Debug("<TimeTrigNode> New time event")
				msg := model.Message{Payload:fimpgo.FimpMessage{Value:node.config.DefaultMsg.Value,ValueType:node.config.DefaultMsg.ValueType},
					Header:map[string]string{"name":node.config.Expressions[i].Name}}
				node.cronMessageCh <- msg
			})

		}
		node.cron.Start()
	}

	return nil
}

// is invoked when node flow is stopped
func (node *TimeTriggerNode) Cleanup() error {
	node.cron.Stop()
	return nil
}

func (node *TimeTriggerNode) OnInput(msg *model.Message) ([]model.NodeID, error) {
	newMsg :=<- node.cronMessageCh
	msg = &newMsg
	return []model.NodeID{node.meta.SuccessTransition}, nil
}
