package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
	"os/exec"
	"io/ioutil"
	"path/filepath"
	"os"
	"encoding/json"
)

type ExecNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	config ExecNodeConfig
	scriptFullPath string
}

type ExecNodeConfig struct {
	ExecType string // cmd , sh-cmd , python , script
	Command string
	ScriptBody string
	InputVariableName string
	IsInputVariableGlobal bool
	OutputVariableName string
	IsOutputVariableGlobal bool
	IsOutputJson bool
	IsInputJson bool
}

func NewExecNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := ExecNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	node.config = ExecNodeConfig{}
	return &node
}

func (node *ExecNode) LoadNodeConfig() error {
	err := mapstructure.Decode(node.meta.Config,&node.config)
	if err != nil{
		log.Error(node.flowOpCtx.FlowId+"<ExecNode> err")
	}
	if node.config.ExecType == "python" {
		node.scriptFullPath = filepath.Join(node.flowOpCtx.StoragePath,node.flowOpCtx.FlowId+"_"+string(node.meta.Id)+".py")
		err = ioutil.WriteFile(node.scriptFullPath, []byte(node.config.ScriptBody), 0644)
	}
	return err
}

// is invoked when node flow is stopped
func (node *ExecNode) Cleanup() error {
	if node.scriptFullPath != "" {
		os.Remove(node.scriptFullPath)
	}
	return nil
}

func (node *ExecNode) WaitForEvent(responseChannel chan model.ReactorEvent) {

}

func (node *ExecNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info(node.flowOpCtx.FlowId+"<ExecNode> Executing ExecNode . Name = ", node.meta.Label)
    var cmd * exec.Cmd
	switch node.config.ExecType {
	case "cmd":
		cmd = exec.Command(node.config.Command)
	case "sh-cmd":
		cmd = exec.Command("bash", "-c", node.config.Command)
	case "python":
		if node.config.IsInputJson {
			strMsg,err := json.Marshal(msg)
			if err != nil {
				return []model.NodeID{node.meta.ErrorTransition},err
			}
			cmd = exec.Command("python3",node.scriptFullPath,string(strMsg))
		}else {
			cmd = exec.Command("python3",node.scriptFullPath)
		}
	}
	output , err := cmd.CombinedOutput()
	log.Debug(node.flowOpCtx.FlowId+"<ExecNode> Normal Output : ", string(output))
	if err != nil {
		log.Debug(node.flowOpCtx.FlowId+"<ExecNode> Err Output : ", err.Error())
	}


	flowId := node.flowOpCtx.FlowId
	outputJson := make(map[string]interface{})
	if node.config.IsOutputJson {
		err = json.Unmarshal(output,&outputJson)

	}
	if err != nil {
		log.Debug(node.flowOpCtx.FlowId+"<ExecNode> Script output can't be unmarshaled to JSON : ", err.Error())
		return []model.NodeID{node.meta.ErrorTransition},err
	}

	if node.config.OutputVariableName != "" {
		if node.config.IsOutputVariableGlobal {
			flowId = "global"
		}
		if node.config.IsOutputJson {
			log.Debug(node.flowOpCtx.FlowId+"<ExecNode> JSON : ", outputJson["ab"])
			err = node.ctx.SetVariable(node.config.OutputVariableName,"object",outputJson,"",flowId,false )
		}else {
			err = node.ctx.SetVariable(node.config.OutputVariableName,"string",string(output),"",flowId,false )
		}

	}else {
		if node.config.IsOutputJson {
			msg.Payload.Value = outputJson
			msg.Payload.ValueType = "object"
		}else {
			msg.Payload.Value = string(output)
			msg.Payload.ValueType = "string"
		}
	}

	if err != nil {
		log.Debug(node.flowOpCtx.FlowId+"<ExecNode> Failed to save variable : ", err.Error())
		return []model.NodeID{node.meta.ErrorTransition},err
	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

