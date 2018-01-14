package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
	"github.com/mitchellh/mapstructure"
)

type TransformNode struct {
	BaseNode
	ctx *model.Context
	nodeConfig TransformNodeConfig
	transport *fimpgo.MqttTransport
}

type TransformNodeConfig struct {
	Description string
	IsVariableGlobal bool                    // true - update global variable ; false - update local variable
	Operation string // type of transform operation , flip , add , subtract , multiply , divide , map , to_bool
	RValue model.Variable // Constant Right variable value .
	RVariableName string // Right variable name , if empty , RValue will be used instead
 	LVariableName string  // Update input message if LVariable is empty

}

func NewTransformNode(flowOpCtx *model.FlowOperationalContext,meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := TransformNode{ctx:ctx,transport:transport}
	node.meta = meta
	node.flowOpCtx = flowOpCtx
	return &node
}

func (node *TransformNode) LoadNodeConfig() error {
	defValue := TransformNodeConfig{}
	err := mapstructure.Decode(node.meta.Config,&defValue)
	if err != nil{
		log.Error(node.flowOpCtx.FlowId+"<Transf> Can't decode configuration",err)
	}else {
		node.nodeConfig = defValue
		node.meta.Config = defValue
	}
	return nil
}

func (node *TransformNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	log.Info(node.flowOpCtx.FlowId+"<Transf> Executing TransformNode . Name = ", node.meta.Label)

	// There are 3 possible sources for RVariable : default value , inputMessage , variable from context
	// There are 2 possible destinations for LVariable : inputMessage , variable from context
	var lValue model.Variable
	var rValue model.Variable
	var err error

	if node.nodeConfig.LVariableName == "" {
		// Use input message
		lValue.Value = msg.Payload.Value
		lValue.ValueType = msg.Payload.ValueType
	} else {
		// Use variable
		if node.nodeConfig.IsVariableGlobal {
			lValue,err = node.ctx.GetVariable(node.nodeConfig.LVariableName,"global")
		}else {
			lValue,err = node.ctx.GetVariable(node.nodeConfig.LVariableName,node.flowOpCtx.FlowId)
		}
	}

	if err != nil {
		return nil , err
	}


	if node.nodeConfig.RVariableName != "" {
		// Use variable
		if node.nodeConfig.IsVariableGlobal {
			rValue,err = node.ctx.GetVariable(node.nodeConfig.RVariableName,"global")
		}else {
			rValue,err = node.ctx.GetVariable(node.nodeConfig.RVariableName,node.flowOpCtx.FlowId)
		}

	}else if node.nodeConfig.RValue.ValueType != "" {
		rValue = node.nodeConfig.RValue
	}else {
		rValue.Value = msg.Payload.Value
		rValue.ValueType = msg.Payload.ValueType
	}

	if err != nil {
		return nil , err
	}

    if lValue.ValueType == rValue.ValueType || (lValue.IsNumber() && rValue.IsNumber())  {
		switch node.nodeConfig.Operation {
		case "flip":
			if lValue.ValueType == "bool" {
				val,ok := rValue.Value.(bool)
				if ok {
					lValue.Value = !val
				}else {
					log.Error(node.flowOpCtx.FlowId+"<Transf> Value type is not bool. Has to bool")
				}
			}else {
				log.Warn(node.flowOpCtx.FlowId+"<Transf> Only bool variable can be flipped")
			}
		case "to_bool":
			if lValue.IsNumber() {
				val,err := rValue.ToNumber()
				if err == nil {
					if val == 0 {
						lValue.Value = false
					} else {
						lValue.Value = true
					}
					lValue.ValueType = "bool"
				}else {
					log.Error(node.flowOpCtx.FlowId+"<Transf> Value type is not number.")
				}
			}else {
				log.Warn(node.flowOpCtx.FlowId+"<Transf> Only numeric value can be converted into bool")
			}
		case "add","subtract","multiply","divide":
			if lValue.IsNumber(){
				rval,err := rValue.ToNumber()
				lval,err := lValue.ToNumber()
				var result float64
				if err == nil {
					switch node.nodeConfig.Operation {
					case "add":
						result = lval + rval
					case "subtract":
						result = lval - rval
					case "multiply":
						result = lval * rval
					case "divide":
						result = lval / rval
					default:
						log.Warn(node.flowOpCtx.FlowId+"<Transf> Unknown arithmetic operator")
					}
					if lValue.ValueType == "float" {
						lValue.Value = result
					}else {
						lValue.Value = int64(result)
					}

				}else {
					log.Error(node.flowOpCtx.FlowId+"<Transf> Value type is not number.")
				}
			}else {
				log.Warn(node.flowOpCtx.FlowId+"<Transf> Only numeric value can be used for arithmetic operations")
			}

		}
	}

	if node.nodeConfig.LVariableName == "" {
		// Update input message
		msg.Payload.Value = lValue.Value
		msg.Payload.ValueType = lValue.ValueType
	}else {
		// Save value into variable
		// Save default value from node config to variable
		if node.nodeConfig.IsVariableGlobal {
				node.ctx.SetVariable(node.nodeConfig.LVariableName, lValue.ValueType, lValue.ValueType, node.nodeConfig.Description, "global", false)
		} else {
				node.ctx.SetVariable(node.nodeConfig.LVariableName, lValue.ValueType, lValue.ValueType, node.nodeConfig.Description, node.flowOpCtx.FlowId, false)
		}

	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

func (node *TransformNode) WaitForEvent(responseChannel chan model.ReactorEvent) {

}
