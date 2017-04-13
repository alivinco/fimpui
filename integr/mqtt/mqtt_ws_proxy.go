package mqtt
import (
	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"fmt"
	"net/http"
	"net"
	//"encoding/hex"
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

func (wu * WsUpgrader) Upgrade(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}



	fmt.Println("Upgraded ")
	session := MqttWsProxySession{wsConn:ws}
	session.Connect(wu.BrokerAddress)
	session.WsReader()
	return nil
}

type MqttWsProxySession struct {
	wsConn *websocket.Conn
	brokerConn net.Conn
}

func (mp *MqttWsProxySession) Connect(address string ) error {
	var err error
	mp.brokerConn, err = net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Can't connect to broker error :",err)
		return err
	}
	go mp.brokerReader()
	return nil
}

func (mp *MqttWsProxySession) WsReader(){

	defer mp.wsConn.Close()
	for {
		msgType, msg, err := mp.wsConn.ReadMessage()
		if err != nil {
			fmt.Println("Read error :",err)
			break
		}else
		if msgType == websocket.BinaryMessage {
			//fmt.Println("Sending packet WS -> broker")
			//fmt.Printf("%s", hex.Dump(msg))
			mp.brokerConn.Write(msg)
		}else {
			fmt.Println(" Message with type = ",msgType)
		}

	}
	fmt.Println("Loop Kaput")
}

func (mp *MqttWsProxySession) brokerReader(){

	for {
		packet := make([]byte,1)
		// reading header byte
		_ ,err := io.ReadFull(mp.brokerConn, packet)
		if err != nil {
			fmt.Println("Can't read packets from broker error =",err)
			break
		}
		// reading length bytes
		packetLen , lenBytes := decodeLength2(mp.brokerConn)
		packet = append(packet,lenBytes...)
		// reading payload
		if packetLen > 0 {
			payload := make([]byte,packetLen)
			io.ReadFull(mp.brokerConn, payload)
			packet = append(packet,payload...)
		}else {
			fmt.Println("Empty payload")
		}
		err = mp.wsConn.WriteMessage(websocket.BinaryMessage,packet)
		if err != nil {
			fmt.Println("Write error :",err)
			mp.wsConn.Close()
			break
		}
	}

}

func decodeLength2(r io.Reader) (int,[]byte) {
	var rLength uint32
	var multiplier uint32
	b := make([]byte, 1)
	var bytes []byte
	for {
		io.ReadFull(r, b)
		bytes = append(bytes,b...)
		digit := b[0]
		rLength |= uint32(digit&127) << multiplier
		if (digit & 128) == 0 {
			break
		}
		multiplier += 7
	}
	return int(rLength),bytes
}