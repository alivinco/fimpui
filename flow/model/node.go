package model

type NodeID string

type MetaNode struct {
	Id                NodeID
	Type              string
	Label             string
	SuccessTransition NodeID
	TimeoutTransition NodeID
	ErrorTransition   NodeID
	Address           string
	Service           string
	ServiceInterface  string
	Config            interface{}
}

type Node interface {
	OnInput( msg *Message) ([]NodeID,error)
	GetMetaNode()*MetaNode
	GetNextSuccessNodes()[]NodeID
	GetNextErrorNode()NodeID
	GetNextTimeoutNode()NodeID
	LoadNodeConfig() error
	IsStartNode() bool
	IsMsgReactorNode() bool
    ConfigureInStream(activeSubscriptions *[]string,msgInStream MsgPipeline)
	Init() error
	Cleanup() error
}







