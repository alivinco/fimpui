package fhcore

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"math/rand"
	"net/url"
	"time"
	log "github.com/Sirupsen/logrus"
)

type VinculumClient struct {
	host            string
	inboundMsgCh    chan []byte
	client          *websocket.Conn
	isRunning       bool
	runningRequests map[int]chan VinculumMsg
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
			}
			if !vc.isRunning {
				break
			}
		}
	}()
	return nil
}

func (vc *VinculumClient) GetMessage(components []string) (VinculumMsg, error) {
	if !vc.isRunning {
		err := vc.Connect()
		if err != nil {
			return VinculumMsg{}, errors.New("Vinculum is Not connected ")
		}
	}
	reqId := rand.Intn(1000)
	msg := VinculumMsg{Ver: "sevenOfNine", Msg: Msg{Type: "request", Src: "fimpui", Dst: "vinculum", Data: Data{Cmd: "get", RequestID: reqId, Param: Param{Components: components}}}}
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

func (vc *VinculumClient) Stop() {
	vc.isRunning = false
}
