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

const (
	SIGNAL_STOP = 1

)

type FlowOperationalContext struct {
	FlowId string
	IsFlowRunning bool
	State string
	NodeControlSignalChannel chan int // the channel should be used to stop all waiting nodes .
	NodeIsReady chan bool // Flow should notify message router when next node is ready to process new message .
}