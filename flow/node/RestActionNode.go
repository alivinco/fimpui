package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
	"github.com/ChrisTrenkamp/goxpath"
	"github.com/ChrisTrenkamp/goxpath/tree/xmltree"
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
	urlTemplate *template.Template
	httpClient  *http.Client
}

type ResponseToVariableMap struct {
	Name string
	Path string
	PathType string // xml , json
	TargetVariableName string
	IsVariableGlobal bool
	TargetVariableType string
	UpdateTriggerMessage bool
}

type Header struct {
	Name string
	Value string
}

type RestActionNodeConfig struct {
	Url string
	Method string // GET,POST,PUT,DELETE, etc.
	TemplateVariableName string
	IsVariableGlobal bool
	RequestPayloadType string // json,xml,string
	RequestTemplate string
	Headers []Header
	ResponseMapping []ResponseToVariableMap
	LogResponse bool
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
			log.Error(node.flowOpCtx.FlowId+"<RestActionNode> Failed while parsing request template.Error:",err)
		}
		node.urlTemplate,err = template.New("url").Parse(node.config.Url)
		if err != nil {
			log.Error(node.flowOpCtx.FlowId+"<RestActionNode> Failed while parsing url template.Error:",err)
		}
	}
	return err
}

func (node *RestActionNode) WaitForEvent(responseChannel chan model.ReactorEvent) {

}

func (node *RestActionNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info(node.flowOpCtx.FlowId+"<RestActionNode> Executing RestActionNode . Name = ", node.meta.Label)

	var templateBuffer bytes.Buffer
	var urlTemplateBuffer bytes.Buffer
	templateParams := RestActionNodeTemplateParams{}
    templateParams.Variable = msg.Payload.Value
    templateParams.Message = msg

	node.reqTemplate.Execute(&templateBuffer,templateParams)
	node.urlTemplate.Execute(&urlTemplateBuffer,templateParams)

	log.Debug("<RestActionNode> Url:",urlTemplateBuffer.String())
	log.Debug("<RestActionNode> Request:",templateBuffer.String())
	req, err := http.NewRequest(node.config.Method, urlTemplateBuffer.String(), &templateBuffer)
	for i := range node.config.Headers{
		req.Header.Add(node.config.Headers[i].Name,node.config.Headers[i].Value)
	}

	if err != nil {
		return []model.NodeID{},err
	}

	resp, err := node.httpClient.Do(req)
	if err != nil {
		return []model.NodeID{},err
	}

	for i := range node.config.ResponseMapping {
		if node.config.ResponseMapping[i].PathType == "xml" {
			xTree, err := xmltree.ParseXML(resp.Body)
			if err == nil {
				var varValue interface{}
				var xpExec = goxpath.MustParse(node.config.ResponseMapping[i].Path)

				switch node.config.ResponseMapping[i].TargetVariableType {
				case "string":
					result, err := xpExec.Exec(xTree)
					if err == nil {
						log.Info("<RestActionNode> Xpath result :",result.String())
						varValue = result.String()
					}
				case "bool":
					result, err := xpExec.ExecBool(xTree)
					if err == nil {
						log.Info("<RestActionNode> Xpath result :",result)
						varValue = result
					}
				case "int":
					result, err := xpExec.ExecNum(xTree)
					if err == nil {
						log.Info("<RestActionNode> Xpath result :",int(result))
						varValue = int(result)
					}
				case "float":
					result, err := xpExec.ExecNum(xTree)
					if err == nil {
						log.Info("<RestActionNode> Xpath result :",result)
						varValue = result
					}
				}

				if err != nil {
					log.Error("<RestActionNode> Can't find result :",err)
				}else {
					flowId := node.flowOpCtx.FlowId
					if node.config.ResponseMapping[i].IsVariableGlobal {
						flowId = "global"
					}
					node.ctx.SetVariable(node.config.ResponseMapping[i].TargetVariableName,node.config.ResponseMapping[i].TargetVariableType,varValue,"",flowId,false )
				}


			}else {
				log.Error("<RestActionNode> Can't parse XML :",err)
			}
			//fmt.Println(res)
		}

	}

	if node.config.LogResponse {
		var respBuff bytes.Buffer
		respBuff.ReadFrom(resp.Body)
		log.Info("<RestActionNode> Response:",respBuff.String())
	}

	log.Infof(node.flowOpCtx.FlowId+"<RestActionNode> Done . Name = %s,Status = %s", node.meta.Label,resp.Status)
	return []model.NodeID{node.meta.SuccessTransition},nil
}

