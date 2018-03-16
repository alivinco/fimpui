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
)

type ExecNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
	config ExecNodeConfig
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
		scripFullPath := filepath.Join(node.flowOpCtx.StoragePath,node.flowOpCtx.FlowId+"_"+string(node.meta.Id)+".py")
		err = ioutil.WriteFile(scripFullPath, []byte(node.config.ScriptBody), 0644)

	}
	return err
}

// is invoked when node flow is stopped
func (node *ExecNode) Cleanup() error {
	scripFullPath := filepath.Join(node.flowOpCtx.StoragePath,node.flowOpCtx.FlowId+"_"+string(node.meta.Id)+".py")
	os.Remove(scripFullPath)
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
	}
	output , err := cmd.CombinedOutput()
	msg.Payload.Value = string(output)
	msg.Payload.ValueType = "string"
	if err != nil {
		return []model.NodeID{node.meta.ErrorTransition},err
	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

