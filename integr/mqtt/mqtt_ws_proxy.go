package mqtt

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"net"
	"net/http"
	//"encoding/hex"
	log "github.com/Sirupsen/logrus"
	"io"
)

var (
	upgrader = websocket.Upgrader{
		Subprotocols: []string{"mqtt"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type WsUpgrader struct {
	BrokerAddress string
}

func (wu *WsUpgrader) Upgrade(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error("<MqWsProxy> Can't upgrade . Error:", err)
		return err
	}

	log.Info("<MqWsProxy> Upgraded ")
	session := MqttWsProxySession{wsConn: ws}
	session.Connect(wu.BrokerAddress)
	session.WsReader()
	return nil
}

type MqttWsProxySession struct {
	wsConn     *websocket.Conn
	brokerConn net.Conn
}

func (mp *MqttWsProxySession) Connect(address string) error {
	var err error
	mp.brokerConn, err = net.Dial("tcp", address)
	if err != nil {
		log.Error("<MqWsProxy> Can't connect to broker . Error :", err)
		return err
	}
	go mp.brokerReader()
	return nil
}

func (mp *MqttWsProxySession) WsReader() {

	defer mp.wsConn.Close()
	for {
		msgType, msg, err := mp.wsConn.ReadMessage()
		if err != nil {
			log.Error("<MqWsProxy> Read error :", err)
			break
		} else if msgType == websocket.BinaryMessage {
			//fmt.Println("Sending packet WS -> broker")
			//fmt.Printf("%s", hex.Dump(msg))
			mp.brokerConn.Write(msg)
		} else {
			log.Debug(" Message with type = ", msgType)
		}

	}
	log.Info("<MqWsProxy> Quit from WsReader loop")
}

func (mp *MqttWsProxySession) brokerReader() {

	for {
		packet := make([]byte, 1)
		// reading header byte
		_, err := io.ReadFull(mp.brokerConn, packet)
		if err != nil {
			log.Error("<MqWsProxy> Can't read packets from broker error =", err)
			break
		}
		// reading length bytes
		packetLen, lenBytes := decodeLength2(mp.brokerConn)
		packet = append(packet, lenBytes...)
		// reading payload
		if packetLen > 0 {
			payload := make([]byte, packetLen)
			io.ReadFull(mp.brokerConn, payload)
			packet = append(packet, payload...)
		} else {
			log.Debug("<MqWsProxy> Empty payload")
		}
		err = mp.wsConn.WriteMessage(websocket.BinaryMessage, packet)
		if err != nil {
			log.Error("<MqWsProxy> Write error :", err)
			mp.wsConn.Close()
			break
		}
	}

}

func decodeLength2(r io.Reader) (int, []byte) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	var bytes []byte
	for {
		io.ReadFull(r, b)
		bytes = append(bytes, b...)
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength), bytes
}
