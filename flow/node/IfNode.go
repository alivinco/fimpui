package node

import (
	log "github.com/Sirupsen/logrus"
	"github.com/alivinco/fimpui/flow/model"
	"errors"
	"github.com/alivinco/fimpui/flow/utils"
)


type IFExpressions struct {
	Expression      []IFExpression
	TrueTransition  model.NodeID
	FalseTransition model.NodeID
}

type IFExpression struct {
	VariableName    string // variable to compare with . Message is used if the value is empty.
	Operand         string // eq , gr , lt
	Value           interface{}
	ValueType       string
	BooleanOperator string // and , or , not
}



func IfNode(node *model.MetaNode, msg *model.Message) error {
	var configValue , msgValue float64
	var err error
	conf, ok := node.Config.(IFExpressions)
	if ok {
		booleanOperator := ""
		var finalResult bool
		for _, item := range conf.Expression {
			if item.ValueType != msg.Payload.ValueType {
				return errors.New("Incompatible value type")
			}
			var result bool
			log.Info("<Node> Operand = ", item.Operand)

			if item.Operand == "gt" || item.Operand == "lt"  {
				if item.ValueType == "int" || item.ValueType  == "float" {
					configValue , err = utils.ConfigValueToNumber(item.ValueType,item.Value)
					if err != nil {
						return err
					}
					msgValue,err = utils.MsgValueToNumber(msg)
					if err != nil {
						return err
					}
				}else {
					return errors.New("Incompatible value type . gt and lt can be used only with numeric types")
				}
			}

			switch item.Operand {
			case "eq":
				result = item.Value == msg.Payload.Value
			case "gt":
				result = msgValue > configValue
			case "lt":
				result = msgValue < configValue
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
			node.SuccessTransition = conf.TrueTransition
		} else {
			node.SuccessTransition = conf.FalseTransition
		}
		return nil
	} else {
		log.Error(node.Config)
		return errors.New("Incompatible node configuration format")
	}

}