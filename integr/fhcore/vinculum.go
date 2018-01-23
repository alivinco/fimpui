package fhcore

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"math/rand"
	"net/url"
	"time"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type VinculumClient struct {
	host            string
	inboundMsgCh    chan []byte
	client          *websocket.Conn
	isRunning       bool
	runningRequests map[int]chan VinculumMsg
	subscribers     []chan VinculumMsg
}

func NewVinculumClient(host string) *VinculumClient {
	vc := VinculumClient{host: host, isRunning: true}
	return &vc
}

func (vc *VinculumClient) Connect() error {
	defer func() {
		if r := recover(); r != nil {
			log.Error("<VincClient> Process CRASHED with error : ",r)
		}
	}()
	vc.runningRequests = make(map[int]chan VinculumMsg)
	u := url.URL{Scheme: "ws", Host: vc.host, Path: "/ws"}
	log.Infof("<VincClient> Connecting to %s", u.String())
	var err error
	vc.client, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error("<VincClient> Dial error",err)
		vc.isRunning = false
		return err
	}

	go func() {
		defer vc.client.Close()
		defer close(vc.inboundMsgCh)
		for {
			vincMsg := VinculumMsg{}
			err := vc.client.ReadJSON(&vincMsg)

			if err != nil {
				//if vincMsg.Msg.Data. != "notify" {
				//	fmt.Println("read:", err)
				//}
				continue
			}
			if vincMsg.Msg.Type == "response" {
				for k, vchan := range vc.runningRequests {
					if k == vincMsg.Msg.Data.RequestID {
						//log.Debugf("Connecting to %s", u.String())
						vchan <- vincMsg
					}
				}
			}else {
				for i := range vc.subscribers{
					select {
					case vc.subscribers[i] <- vincMsg:
					default:
						log.Infof("<VincClient> No listeners on the channel")
					}
				}
			}

			if !vc.isRunning {
				break
			}
		}
	}()
	return nil
}

func (vc *VinculumClient) RegisterSubscriber() chan VinculumMsg {
	newChan := make(chan VinculumMsg)
	vc.subscribers = append(vc.subscribers,newChan)
	return newChan

}

func (vc *VinculumClient) GetMessage(components []string) (VinculumMsg, error) {
	if !vc.isRunning {
		err := vc.Connect()
		if err != nil {
			return VinculumMsg{}, errors.New("Vinculum is Not connected ")
		}
	}
	reqId := rand.Intn(1000)
	msg := VinculumMsg{Ver: "sevenOfNine", Msg: Msg{Type: "request", Src: "fimpui", Dst: "vinculum", Data: Data{Cmd: "get",RequestID: reqId, Param: Param{Components: components}}}}
	vc.client.WriteJSON(msg)
	vc.runningRequests[reqId] = make(chan VinculumMsg)
	select {
	case msg := <-vc.runningRequests[reqId]:
		delete(vc.runningRequests, reqId)
		return msg, nil
	case <-time.After(time.Second * 5):
		log.Info("<VincClient> imeout 5 sec")
	}
	delete(vc.runningRequests, reqId)

	return VinculumMsg{}, errors.New("Timeout")

}

func (vc *VinculumClient) SetMode(mode string) error {
	reqId := rand.Intn(1000)
	msg := VinculumMsg{Ver: "sevenOfNine", Msg: Msg{Type: "request", Src: "fimpui", Dst: "vinculum", Data: Data{Cmd: "set", RequestID: reqId,Component:"mode",Id:mode}}}
	return vc.client.WriteJSON(msg)
}

func (vc *VinculumClient) SetShortcut(shortcutId string) error {
	reqId := rand.Intn(1000)
	numId,err := strconv.ParseInt(shortcutId,10,64)
	if err == nil {
		msg := VinculumMsg{Ver: "sevenOfNine", Msg: Msg{Type: "request", Src: "fimpui", Dst: "vinculum", Data: Data{Cmd: "set", RequestID: reqId,Component:"shortcut",Id:numId}}}
		return vc.client.WriteJSON(msg)
	}
	return err

}

func (vc *VinculumClient) GetShortcuts() ([]Shortcut,error) {
	vincMsg , err := vc.GetMessage([]string{"shortcut"})
	return vincMsg.Msg.Data.Param.Shortcut,err
}

func (vc *VinculumClient) Stop() {
	vc.isRunning = false
}

/*
**** Mode change event *****

{
    "msg": {
        "data": {
            "cmd": "set",
            "component": "house",
            "id": null,
            "param": {
                "learning": null,
                "mode": "home",
                "time": "2018-01-11T21:53:09Z",
                "uptime": 1510545049
            }
        },
        "dst": "clients",
        "type": "notify"
    },
    "ver": "sevenOfNine"
}



**** Mode set command ******

 {
    "ver": "sevenOfNine",
    "msg": {
        "type": "request",
        "src": "chat",
        "dst": "vinculum",
        "data": {
            "cmd": "set",
            "component": "mode",
            "id": "home",
            "param": null,
            "requestId": 1
        }
    }
}

****** GEt list of shortcuts *****************

{
  "ver": "sevenOfNine",
  "msg": {
    "type": "request",
    "src": "chat",
    "data": {
      "cmd": "get",
      "component": null,
      "requestId": 1,
      "param": {
        "components": [
          "shortcut"
        ]
      }
    }
  }
}


******* Trigger shortcut ************************

{
  "ver": "sevenOfNine",
  "msg": {
    "type": "request",
    "src": "chat",
    "data": {
      "cmd": "set",
      "component": "shortcut",
      "id": 1,
      "requestId": 4,
      "param": null
    }
  }
}

***********************************************


 */