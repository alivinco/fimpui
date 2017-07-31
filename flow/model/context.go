package model

import "github.com/pkg/errors"

type Variable struct {
	Value     interface{}
	ValueType string
}
type Context struct {
	IsFlowRunning bool
	variables map[string]Variable
}

func NewContext() Context {
	ctx := Context{}
	ctx.variables = make(map[string]Variable)
	return ctx
}

func (ctx *Context) SetVariable(name string,valueType string,value interface{}) {
	ctx.variables[name] = Variable{ValueType:valueType,Value:value}
}

func (ctx *Context) GetVariable(name string) (Variable,error) {
	variable , ok := ctx.variables[name]
	if ok {
		return variable,nil
	}
	return Variable{},errors.New("Variable doesn't exist")
}

func (ctx *Context) GetAllVariables() map[string]Variable  {
	return ctx.variables
}
