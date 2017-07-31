package node

import (
	"github.com/alivinco/fimpgo"
	"github.com/alivinco/fimpui/flow/model"
)

type Constructor func(meta model.MetaNode,ctx *model.Context,transport *fimpgo.MqttTransport) model.Node

var Registry = map[string]Constructor {
	"if":NewIfNode,
	"action":NewActionNode,
	"wait":NewWaitNode,
}
