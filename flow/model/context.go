package model


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

func (ctx *Context) GetVariable(name string) (string , interface{}) {
	return ctx.variables[name].ValueType,ctx.variables[name].Value
}

func (ctx *Context) GetAllVariables() map[string]Variable  {
	return ctx.variables
}
