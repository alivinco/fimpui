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

type ValueMappingRecord struct {
	LValue model.Variable
	RValue model.Variable
}

type TransformNodeConfig struct {
	TargetVariableName string  // Variable
	IsTargetVariableGlobal bool
	TransformType string       // map , calc
	IsRVariableGlobal bool                    // true - update global variable ; false - update local variable
	IsLVariableGlobal bool                    // true - update global variable ; false - update local variable
	Operation string 			// type of transform operation , flip , add , subtract , multiply , divide , to_bool
	RType     string            // var , const
	RValue model.Variable 		// Constant Right variable value .
	RVariableName string 		// Right variable name , if empty , RValue will be used instead
 	LVariableName string  		// Update input message if LVariable is empty
 	ValueMapping []ValueMappingRecord // ["LValue":1,"RValue":"mode-1"]
 	//value mapping
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
	var result model.Variable
	var err error

	if node.nodeConfig.LVariableName == "" {
		// Use input message
		lValue.Value = msg.Payload.Value
		lValue.ValueType = msg.Payload.ValueType
	} else {
		// Use variable
		if node.nodeConfig.IsLVariableGlobal {
			lValue,err = node.ctx.GetVariable(node.nodeConfig.LVariableName,"global")
		}else {
			lValue,err = node.ctx.GetVariable(node.nodeConfig.LVariableName,node.flowOpCtx.FlowId)
		}
	}

	if err != nil {
		return nil , err
	}

    if node.nodeConfig.RType == "var" {
		if node.nodeConfig.RVariableName == "" {
			rValue.Value = msg.Payload.Value
			rValue.ValueType = msg.Payload.ValueType
		}else{
			// Use variable
			if node.nodeConfig.IsRVariableGlobal {
				rValue,err = node.ctx.GetVariable(node.nodeConfig.RVariableName,"global")
			}else {
				rValue,err = node.ctx.GetVariable(node.nodeConfig.RVariableName,node.flowOpCtx.FlowId)
			}
		}
	}else {
		rValue = node.nodeConfig.RValue
	}


	if err != nil {
		return nil , err
	}

    if lValue.ValueType == rValue.ValueType || (lValue.IsNumber() && rValue.IsNumber())  {

    	if node.nodeConfig.TransformType == "calc" {
			switch node.nodeConfig.Operation {
			case "flip":
				if lValue.ValueType == "bool" {
					val,ok := rValue.Value.(bool)
					if ok {
						result.Value = !val
						result.ValueType = rValue.ValueType
					}else {
						log.Error(node.flowOpCtx.FlowId+"<Transf> Value type is not bool. Has to bool")
					}
				}else {
					log.Warn(node.flowOpCtx.FlowId+"<Transf> Only bool variable can be flipped")
				}
			case "to_bool":
				if lValue.IsNumber() {
					val,err := lValue.ToNumber()
					if err == nil {
						if val == 0 {
							result.Value = false
						} else {
							result.Value = true
						}
						result.ValueType = "bool"
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
					var calcResult float64
					if err == nil {
						switch node.nodeConfig.Operation {
						case "add":
							calcResult = lval + rval
						case "subtract":
							calcResult = lval - rval
						case "multiply":
							calcResult = lval * rval
						case "divide":
							calcResult = lval / rval
						default:
							log.Warn(node.flowOpCtx.FlowId+"<Transf> Unknown arithmetic operator")
						}
						if rValue.ValueType == "float" {
							result.Value = calcResult
						}else {
							result.Value = int64(calcResult)
						}
						result.ValueType = lValue.ValueType

					}else {
						log.Error(node.flowOpCtx.FlowId+"<Transf> Value type is not number.")
					}
				}else {
					log.Warn(node.flowOpCtx.FlowId+"<Transf> Only numeric value can be used for arithmetic operations")
				}

			}
		}else if node.nodeConfig.TransformType == "map" {
			for i := range node.nodeConfig.ValueMapping {
				log.Debug(node.flowOpCtx.FlowId+"<Transf> record Value ",node.nodeConfig.ValueMapping[i].LValue.Value)
				log.Debug(node.flowOpCtx.FlowId+"<Transf> record input Value = ",lValue.Value )
				if lValue.ValueType == node.nodeConfig.ValueMapping[i].LValue.ValueType {
					varsAreEqual , err :=  lValue.IsEqual(&node.nodeConfig.ValueMapping[i].LValue)
					if err != nil {
						log.Warn(node.flowOpCtx.FlowId+"<Transf> Error while comparing map vars : ",err)
					}
					if varsAreEqual {
						result = node.nodeConfig.ValueMapping[i].RValue
						log.Debug(node.flowOpCtx.FlowId+"<Transf> Result is set")
						break
					}
				}
			}
		}



	}

	if node.nodeConfig.TargetVariableName == "" {
		// Update input message
		msg.Payload.Value = result.Value
		msg.Payload.ValueType = result.ValueType
	}else {
		// Save value into variable
		// Save default value from node config to variable
		if node.nodeConfig.IsTargetVariableGlobal {
			    node.ctx.SetVariable(node.nodeConfig.TargetVariableName, result.ValueType, result.Value, "", "global", false)
		} else {
				node.ctx.SetVariable(node.nodeConfig.TargetVariableName, result.ValueType, result.Value, "", node.flowOpCtx.FlowId, false)
		}

	}
	return []model.NodeID{node.meta.SuccessTransition},nil
}

func (node *TransformNode) WaitForEvent(responseChannel chan model.ReactorEvent) {

}
