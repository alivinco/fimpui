package fhcore

import "time"

type Data struct {
	Errors    interface{} `json:"errors"`
	Cmd       string      `json:"cmd"`
	Component interface{} `json:"component"`
	Param     Param       `json:"param"`
	RequestID int         `json:"requestId"`
	Success   bool        `json:"success"`
}

type Param struct {
	Components []string `json:"components"`
	Device     []Device `json:"device,omitempty"`
	House      House    `json:"house,omitempty"`
}

type Msg struct {
	Type string `json:"type"`
	Src  string `json:"src"`
	Dst  string `json:"dst"`
	Data Data   `json:"data"`
}

type VinculumMsg struct {
	Ver string `json:"ver"`
	Msg Msg    `json:"msg"`
}

type Fimp struct {
	Adapter string `json:"adapter"`
	Address string `json:"address"`
	Group   string `json:"group"`
}

type Client struct {
	Name string `json:"name"`
}

type Device struct {
	Fimp   Fimp `json:"_fimp"`
	Client Client `json:"client"`
	Functionality string `json:"functionality"`
	ID            int    `json:"id"`
	Lrn           bool   `json:"lrn"`
	Model         string `json:"model"`
	Param         interface {} `json:"param"`
	Problem bool        `json:"problem"`
	Room    interface{} `json:"room"`
}

type House struct {
	Learning interface{} `json:"learning"`
	Mode     string      `json:"mode"`
	Time     time.Time   `json:"time"`
	Uptime   int         `json:"uptime"`
}
