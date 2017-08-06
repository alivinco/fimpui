package model

import "github.com/alivinco/fimpgo"

type MsgPipeline chan Message

type Message struct {
	AddressStr string
	Address    fimpgo.Address
	Payload    fimpgo.FimpMessage
	Header     map[string]string
}

type FlowMeta struct {
	Id          string
	Name        string
	Description string
	Nodes       []MetaNode
}

type FlowOperationalContext struct {
	FlowId string
	IsFlowRunning bool
	State string
}