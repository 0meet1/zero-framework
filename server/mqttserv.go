package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/0meet1/zero-framework/global"
)

const (
	CORE_MQTT_SERVER = "##!CORE_MQTT_SERVER"
)

type MqttMessageListener interface {
	Publish(ZeroConnect, *MqttMessage) error
}

type xMqttConnectBuilder struct{}

func (xDefault *xMqttConnectBuilder) NewConnect() ZeroConnect {
	return &MqttConnect{
		ZeroSocketConnect: ZeroSocketConnect{
			zserv: global.Value(CORE_MQTT_SERVER).(*MqttServer),
		},
		topcis: make(map[string]byte),
	}
}

type MqttConnect struct {
	ZeroSocketConnect

	topcis               map[string]byte
	messageSerialNnumber uint16
	// connectMutex         sync.Mutex
	serialNnumberMutex sync.Mutex

	// active bool

	xListener MqttMessageListener
}

func (mqttconn *MqttConnect) AddListener(xListener MqttMessageListener) {
	mqttconn.xListener = xListener
}

func (mqttconn *MqttConnect) RegisterId() string {
	return mqttconn.This().(ZeroConnect).RemoteAddr()
}

func (mqttconn *MqttConnect) Accept(zserv ZeroServ, connect net.Conn) error {
	mqttconn.ZeroSocketConnect.Accept(zserv, connect)
	mqttconn.topcis = make(map[string]byte)
	mqttconn.messageSerialNnumber = 0
	return nil
}

func (mqttconn *MqttConnect) Close() error {
	err := mqttconn.ZeroSocketConnect.Close()

	mqttserv := mqttconn.zserv.(*MqttServer)
	mqttserv.topicsMapMutex.Lock()
	for topic := range mqttconn.topcis {
		_, ok := mqttserv.topicsMap[topic]
		if ok {
			delete(mqttserv.topicsMap[topic], mqttconn.RemoteAddr())
			if len(mqttserv.topicsMap[topic]) <= 0 {
				delete(mqttserv.topicsMap, topic)
			}
		}
	}
	mqttserv.topicsMapMutex.Unlock()

	return err
}

func (mqttconn *MqttConnect) UpdateSerialNnumber(serialNnumber uint16) {
	mqttconn.serialNnumberMutex.Lock()
	if mqttconn.messageSerialNnumber < serialNnumber || serialNnumber < 10 {
		mqttconn.messageSerialNnumber = serialNnumber
	}
	mqttconn.serialNnumberMutex.Unlock()
}

func (mqttconn *MqttConnect) UseSerialNnumber() uint16 {
	mqttconn.serialNnumberMutex.Lock()
	mqttconn.messageSerialNnumber++
	serialNnumber := mqttconn.messageSerialNnumber
	mqttconn.serialNnumberMutex.Unlock()
	return serialNnumber
}

func (mqttconn *MqttConnect) OnMessage(datas []byte) error {
	mqttMessage := &MqttMessage{}
	err := mqttMessage.build(datas)
	if err != nil {
		global.Logger().Error(fmt.Sprintf("mqtt server connect %s message error %s", mqttconn.RemoteAddr(), err.Error()))
		return err
	}
	global.Logger().Debug(fmt.Sprintf("mqtt connect %s on message type `%s`", mqttconn.RemoteAddr(), mqttMessage.FixedHeader().MessageTypeString()))
	err = mqttconn.onMqttMessage(mqttMessage)
	if err != nil {
		global.Logger().Error(fmt.Sprintf("mqtt server connect %s on message error %s", mqttconn.RemoteAddr(), err.Error()))
	}
	return err
}

func (mqttconn *MqttConnect) onConnect(_ *MqttMessage) error {
	message := &MqttMessage{}
	message.MakeConnackMessage()
	return mqttconn.This().(ZeroConnect).Write(message.Bytes())
}

func (mqttconn *MqttConnect) onPingreq(_ *MqttMessage) error {
	global.Logger().Info(fmt.Sprintf("mqtt connect %s on pingreq", mqttconn.This().(ZeroConnect).RemoteAddr()))

	message := &MqttMessage{}
	message.MakePingrespMessage()
	defer mqttconn.This().(ZeroConnect).Heartbeat()
	return mqttconn.This().(ZeroConnect).Write(message.Bytes())
}

func (mqttconn *MqttConnect) onSubscribe(mqttMessage *MqttMessage) error {
	results := make([]byte, 0)
	for _, topic := range mqttMessage.Payload().(*MqttSubscribePayload).topics {
		mqttconn.topcis[topic.TopicName] = topic.Qos
		results = append(results, topic.Qos)
	}
	mqttconn.This().(ZeroConnect).Authorized()

	mqttserv := mqttconn.zserv.(*MqttServer)
	mqttserv.topicsMapMutex.Lock()
	for topic := range mqttconn.topcis {
		_, ok := mqttserv.topicsMap[topic]
		if !ok {
			mqttserv.topicsMap[topic] = make(map[string]*MqttConnect)
		}
		mqttserv.topicsMap[topic][mqttconn.This().(ZeroConnect).RemoteAddr()] = mqttconn
	}
	mqttserv.topicsMapMutex.Unlock()

	message := &MqttMessage{}
	message.MakeSubackMessage(mqttMessage.VariableHeader().(*MqttIdentifierVariableHeader).Identifier(), results)
	return mqttconn.This().(ZeroConnect).Write(message.Bytes())
}

func (mqttconn *MqttConnect) onPublish(mqttMessage *MqttMessage) error {
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("mqttserv on publish err : %s", err))
		}
	}()

	if mqttconn.xListener != nil {
		err := mqttconn.xListener.Publish(mqttconn.This().(ZeroConnect), mqttMessage)
		if err != nil {
			global.Logger().Error(fmt.Sprintf("mqttserv process publish err : %s", err))
		}
	}

	mqttconn.UpdateSerialNnumber(mqttMessage.VariableHeader().(*MqttPublishVariableHeader).Identifier())
	message := &MqttMessage{}
	message.MakePubackMessage(mqttMessage.VariableHeader().(*MqttPublishVariableHeader).Identifier())
	return mqttconn.This().(ZeroConnect).Write(message.Bytes())
}

func (mqttconn *MqttConnect) onPubrec(mqttMessage *MqttMessage) error {
	message := &MqttMessage{}
	message.MakePubrelMessage(mqttMessage.VariableHeader().(*MqttIdentifierVariableHeader).Identifier())
	return mqttconn.This().(ZeroConnect).Write(message.Bytes())
}

func (mqttconn *MqttConnect) onMqttMessage(mqttMessage *MqttMessage) error {
	defer mqttconn.Heartbeat()

	switch mqttMessage.FixedHeader().MessageType() {
	case CONNECT:
		return mqttconn.onConnect(mqttMessage)
	case CONNACK:
	case PUBLISH:
		return mqttconn.onPublish(mqttMessage)
	case PUBACK:
	case PUBREC:
		return mqttconn.onPubrec(mqttMessage)
	case PUBREL:
	case PUBCOMP:
	case SUBSCRIBE:
		return mqttconn.onSubscribe(mqttMessage)
	case SUBACK:
	case UNSUBSCRIBE:
	case UNSUBACK:
	case PINGREQ:
		return mqttconn.onPingreq(mqttMessage)
	case PINGRESP:
	case DISCONNECT:
	default:
	}
	return nil
}

type MqttServer struct {
	TCPServer

	topicsMap      map[string]map[string]*MqttConnect
	topicsMapMutex sync.RWMutex
}

func NewMqttServer(address string, authWaitSeconds int64, heartbeatSeconds int64, heartbeatCheckInterval int64, bufferSize int) *MqttServer {
	return &MqttServer{
		TCPServer: *NewTCPServer(address, authWaitSeconds, heartbeatSeconds, heartbeatCheckInterval, bufferSize),
		topicsMap: make(map[string]map[string]*MqttConnect),
	}
}

func (mqttserv *MqttServer) RunServer() {
	if mqttserv.ConnectBuilder == nil {
		mqttserv.ConnectBuilder = &xMqttConnectBuilder{}
	}
	mqttserv.TCPServer.RunServer()
	global.Key(CORE_MQTT_SERVER, mqttserv)
}
