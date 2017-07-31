package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpui/flow/model"
	"errors"
	"github.com/alivinco/fimpui/flow/utils"
	"github.com/alivinco/fimpgo"
	//"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/mapstructure"
)


type IFExpressions struct {
	Expression      []IFExpression
	TrueTransition  model.NodeID
	FalseTransition model.NodeID
}

type IFExpression struct {
	LeftVariableName    string  // Left variable of expression . If empty , Message value will be used .
	LeftVariable    model.Variable  `json:"-"` // Right variable of expression . Have to be defined , empty value will generate error .
	RightVariable   model.Variable // Right variable of expression . Have to be defined , empty value will generate error .
	Operand         string // eq , gr , lt
	BooleanOperator string // and , or , not
}

type IfNode struct {
	BaseNode
	ctx *model.Context
	transport *fimpgo.MqttTransport
}

func NewIfNode(meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node {
	node := IfNode{ctx:ctx,transport:transport}
	node.meta = meta
	return &node
}

func (node *IfNode) LoadNodeConfig() error {
	exp := IFExpressions{}
	err := mapstructure.Decode(node.meta.Config,&exp)
	if err != nil{
		log.Error(err)
	}else {
		node.meta.Config = exp
	}
	return nil
}

func (node *IfNode) OnInput( msg *model.Message) ([]model.NodeID,error) {
	var leftNumericValue , rightNumericValue float64
	var err error
	conf, ok := node.meta.Config.(IFExpressions)
	if ok {
		booleanOperator := ""
		var finalResult bool
		for _, item := range conf.Expression {
			if item.RightVariable.ValueType == "" {
				return nil,errors.New("Right variable is not defined. IfNode is skipped.")
			}
			if item.LeftVariableName == ""{
				item.LeftVariable = model.Variable{ValueType:msg.Payload.ValueType,Value:msg.Payload.Value}
			}else {
				item.LeftVariable ,err = node.ctx.GetVariable(item.LeftVariableName)
			}
			if item.LeftVariable.ValueType != item.RightVariable.ValueType {
				return nil,errors.New(" Right and left of expression have different types ")
			}


			var result bool
			log.Debug("<IfNode> Operand = ", item.Operand)

			if item.Operand == "gt" || item.Operand == "lt"  {
				if item.LeftVariable.ValueType == "int" || item.LeftVariable.ValueType  == "float" {
					leftNumericValue , err = utils.ConfigValueToNumber(item.LeftVariable.ValueType,item.LeftVariable.Value)
					if err != nil {
						log.Error("<IfNode> Error while converting left variable to number.Error : ",err)
						return nil,err
					}
					rightNumericValue , err = utils.ConfigValueToNumber(item.RightVariable.ValueType,item.RightVariable.Value)
					if err != nil {
						log.Error("<IfNode> Error while converting right variable to number.Error : ",err)
						return nil,err
					}

				}else {
					return nil,errors.New("Incompatible value type . gt and lt can be used only with numeric types")
				}
			}
			log.Debug("<IfNode> Left numeric value = ", leftNumericValue)
			log.Debug("<IfNode> Right numeric value = ", rightNumericValue)
			switch item.Operand {
			case "eq":
				result = item.LeftVariable.Value == item.RightVariable.Value
			case "gt":
				result = leftNumericValue > rightNumericValue
			case "lt":
				result = leftNumericValue < rightNumericValue
			}
			
			switch booleanOperator {
			case "":
				finalResult = result
			case "and":
				finalResult = finalResult && result
				booleanOperator = item.BooleanOperator
			case "or":
				finalResult = finalResult || result
				booleanOperator = item.BooleanOperator
			case "not":
				finalResult = !result
			}
		}
		if finalResult {
			return []model.NodeID{conf.TrueTransition},nil
		} else {
			return []model.NodeID{conf.FalseTransition},nil
		}
		return nil,nil
	} else {
		log.Error(node.meta.Config)
		return nil, errors.New("Incompatible node configuration format")
	}

	return []model.NodeID{node.meta.SuccessTransition},nil
}


//func IfNodeF(ctx *model.Context,node *model.MetaNode, msg *model.Message) error {
//	var leftNumericValue , rightNumericValue float64
//	var err error
//	conf, ok := node.Config.(IFExpressions)
//	if ok {
//		booleanOperator := ""
//		var finalResult bool
//		for _, item := range conf.Expression {
//			if item.RightVariable.ValueType == "" {
//				return errors.New("Right variable is not defined. IfNode is skipped.")
//			}
//			if item.LeftVariableName == ""{
//				item.LeftVariable = model.Variable{ValueType:msg.Payload.ValueType,Value:msg.Payload.Value}
//			}else {
//				item.LeftVariable ,err = ctx.GetVariable(item.LeftVariableName)
//			}
//			if item.LeftVariable.ValueType != item.RightVariable.ValueType {
//				return errors.New(" Right and left of expression have different types ")
//			}
//
//
//			var result bool
//			log.Info("<IfNode> Operand = ", item.Operand)
//
//			if item.Operand == "gt" || item.Operand == "lt"  {
//				if item.LeftVariable.ValueType == "int" || item.LeftVariable.ValueType  == "float" {
//					leftNumericValue , err = utils.ConfigValueToNumber(item.LeftVariable.ValueType,item.LeftVariable.Value)
//					if err != nil {
//						return err
//					}
//					rightNumericValue , err = utils.ConfigValueToNumber(item.RightVariable.ValueType,item.RightVariable.Value)
//					if err != nil {
//						return err
//					}
//
//				}else {
//					return errors.New("Incompatible value type . gt and lt can be used only with numeric types")
//				}
//			}
//
//			switch item.Operand {
//			case "eq":
//				result = item.LeftVariable.Value == item.RightVariable.Value
//			case "gt":
//				result = leftNumericValue > rightNumericValue
//			case "lt":
//				result = leftNumericValue < rightNumericValue
//			}
//			switch booleanOperator {
//			case "":
//				finalResult = result
//			case "and":
//				finalResult = finalResult && result
//				booleanOperator = item.BooleanOperator
//			case "or":
//				finalResult = finalResult || result
//				booleanOperator = item.BooleanOperator
//			case "not":
//				finalResult = !result
//			}
//		}
//		if finalResult {
//			node.SuccessTransition = conf.TrueTransition
//		} else {
//			node.SuccessTransition = conf.FalseTransition
//		}
//		return nil
//	} else {
//		log.Error(node.Config)
//		return errors.New("Incompatible node configuration format")
//	}
//
//}