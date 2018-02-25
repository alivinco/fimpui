package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"text/template"
	"bytes"
)

type RestActionNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	config RestActionNodeConfig
	reqTemplate *template.Template
	httpClient  *http.Client
}

type RestActionNodeConfig struct {
	VariableName string
	IsVariableGlobal bool
	Method string // GET,POST,PUT,DELETE, etc.
	PayloadType string // json,xml,string
	RequestTemplate string
	Url string
	Headers map[string]string
}

type RestActionNodeTemplateParams struct {
	Variable interface{}
	Message *model.Message
}

func NewRestActionNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := RestActionNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	node.config = RestActionNodeConfig{}
	node.httpClient = &http.Client{}
	return &node
}

func (node *RestActionNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		log.Error(node.flowOpCtx.FlowId+"<RestActionNode> Failed while loading configurations.Error:",err)

	}else {
		node.reqTemplate,err = template.New("request").Parse(node.config.RequestTemplate)
		if err != nil {
			log.Error(node.flowOpCtx.FlowId+"<RestActionNode> Failed while parsing template.Error:",err)
		}
	}
	return err
}

func (node *RestActionNode) WaitForEvent(responseChannel chan model.ReactorEvent) {

}

func (node *RestActionNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info(node.flowOpCtx.FlowId+"<RestActionNode> Executing RestActionNode . Name = ", node.meta.Label)

	var templateBuffer bytes.Buffer
	templateParams := RestActionNodeTemplateParams{}
    templateParams.Variable = msg.Payload.Value
    templateParams.Message = msg
	node.reqTemplate.Execute(&templateBuffer,templateParams)
	log.Debug("<RestActionNode> Request:",templateBuffer.String())
	req, err := http.NewRequest(node.config.Method, node.config.Url, &templateBuffer)
	if err != nil {
		return []model.NodeID{},err
	}

	resp, err := node.httpClient.Do(req)
	if err != nil {
		return []model.NodeID{},err
	}
	var respBuff bytes.Buffer
	respBuff.ReadFrom(resp.Body)
	log.Debug("<RestActionNode> Response:",respBuff.String())
	log.Infof(node.flowOpCtx.FlowId+"<RestActionNode> Done . Name = %s,Status = %s", node.meta.Label,resp.Status)
	return []model.NodeID{node.meta.SuccessTransition},nil
}

