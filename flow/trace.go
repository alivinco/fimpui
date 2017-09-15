package flow

import (
	"github.com/alivinco/fimpui/flow/model"
	"time"
)

type FlowTrace struct {
	records []TraceRecord
}

func NewFlowTrace (size int) *FlowTrace {
	trace := FlowTrace{}
	trace.records = make([]TraceRecord,0,size)
	return &trace
}

type TraceRecord struct {
    NodeId model.NodeID
    TimeStamp time.Time

}

func (tr *FlowTrace) AddRecord(NodeId model.NodeID) {

}