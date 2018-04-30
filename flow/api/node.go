package api

import (
"github.com/labstack/echo"
	"github.com/alivinco/fimpui/flow"
)

type NodeApi struct {
	flowManager *flow.Manager
	echo *echo.Echo
}

func NewNodeApi(flowManager *flow.Manager,echo *echo.Echo) *NodeApi {
	ctxApi := NodeApi{flowManager:flowManager,echo:echo}
	ctxApi.RegisterRestApi()
	return &ctxApi
}

func (ctx * NodeApi) RegisterRestApi() {

}

func (ctx * NodeApi) RegisterMqttApi() {

}