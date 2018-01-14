package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
)

type Constructor func(context *model.FlowOperationalContext, meta model.MetaNode, ctx *model.Context, transport *fimpgo.MqttTransport) model.Node

var Registry = map[string]Constructor{
	"trigger":      NewTriggerNode,
	"receive":      NewReceiveNode,
	"if":           NewIfNode,
	"action":       NewActionNode,
	"wait":         NewWaitNode,
	"set_variable": NewSetVariableNode,
	"loop":         NewLoopNode,
	"time_trigger": NewTimeTriggerNode,
	"transform":    NewTransformNode,
}
