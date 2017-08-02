package model

import (
	"github.com/pkg/errors"
	"time"
)

type Variable struct {
	Value     interface{}
	ValueType string
}

type ContextRecord struct {
	Name string
	UpdatedAt time.Time
	Variable Variable
}

type Context struct {
	IsFlowRunning bool
	IsGlobal bool
	records map[string]ContextRecord
	parentContext *Context // Normally it's global or shared context
}

func NewContext(parentContext *Context) *Context {
	ctx := Context{}
	ctx.records = make(map[string]ContextRecord)
	ctx.parentContext = parentContext
	return &ctx
}

func (ctx *Context) SetVariable(name string,valueType string,value interface{}) {
	ctx.records[name] = ContextRecord{Name:name,UpdatedAt:time.Now(),Variable: Variable{ValueType:valueType,Value:value}}
}

func (ctx *Context) GetParentContext() *Context {
	return ctx.parentContext
}

func (ctx *Context) GetVariable(name string) (Variable,error) {
	rec , ok := ctx.records[name]
	if ok {
		return rec.Variable,nil
	}else {
		// global context lookup
		if ctx.parentContext != nil {
			rec , ok = ctx.parentContext.records[name]
			if ok {
				return rec.Variable, nil
			}
		}
	}
	return Variable{},errors.New("Variable doesn't exist")
}

func (ctx *Context) GetRecord(name string) (*ContextRecord,error) {
	rec , ok := ctx.records[name]
	if ok {
		return &rec,nil
	}
	return &ContextRecord{},errors.New("Variable doesn't exist")
}

func (ctx *Context) GetRecords() map[string]ContextRecord  {
	return ctx.records
}
