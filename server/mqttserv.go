package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/0meet1/zero-framework/global"
)

const (
	CORE_MQTT_SERVER = "X##!CORE_MQTT_SERVER"
)

type MqttMessageListener interface {
	Publish(ZeroConnect, *MqttMessage) error
}

type MqttConnectBuilder struct{}

func (xDefault *MqttConnectBuilder) NewConnect() ZeroConnect {
	mqttconn := &MqttConnect{
		topcis: make(map[string]byte),
	}
	mqttconn.ThisDef(mqttconn)
	mqttconn.AddChecker(&xMqttDataChecker{})
	return mqttconn
}

type xMqttDataChecker struct {
	cachebytes      []byte
	fixedheader     *MqttFixedHeader
	cachebytesMutex sync.Mutex
}

func (checker *xMqttDataChecker) expectedLength(data []byte) []byte {
	xLengthBytes := make([]byte, 0)
	end := 5
	if len(data) < end {
		end = len(data)
	}
	for i := 1; i < end; i++ {
		xLengthBytes = append(xLengthBytes, data[i])
		flag := data[i] >> 7 & 0b00000001
		if flag == 0b0 {
			return xLengthBytes
		}
	}
	return nil
}

func (checker *xMqttDataChecker) unpacking(registerId string, historys ...[]byte) [][]byte {
	comps := make([][]byte, 0)
	comps = append(comps, historys...)
	comps = append(comps, checker.cachebytes[:checker.fixedheader.LessLength()+2])
	checker.cachebytes = checker.cachebytes[checker.fixedheader.LessLength()+2:]
	checker.fixedheader = nil
	expectedLength := checker.expectedLength(checker.cachebytes)
	if expectedLength != nil {
		checker.fixedheader = &MqttFixedHeader{}
		checker.fixedheader.With(checker.cachebytes[0], expectedLength)
	}

	if checker.fixedheader != nil {
		nLen := len(checker.cachebytes) - len(checker.fixedheader.length) - 1
		if nLen == checker.fixedheader.LessLength() {
			bts := make([]byte, len(checker.cachebytes))
			copy(bts, checker.cachebytes)
			checker.cachebytes = nil
			checker.fixedheader = nil
			comps = append(comps, bts)
			return comps
		} else {
			if nLen > checker.fixedheader.LessLength() {
				return checker.unpacking(registerId, comps...)
			}
		}
	}
	return comps
}

func (checker *xMqttDataChecker) CheckPackageData(registerId string, data []byte) [][]byte {
	checker.cachebytesMutex.Lock()
	defer func() {
		checker.cachebytesMutex.Unlock()
		err := recover()
		if err != nil {
			global.Logger().Errorf("mqttconn %s on check package err : %s", registerId, err.(error).Error())
		}
	}()
	// global.Logger().Debugf("mqttconn %s on check %s", registerId, structs.BytesString(data...))

	if checker.cachebytes == nil {
		checker.cachebytes = make([]byte, 0)
	}
	checker.cachebytes = append(checker.cachebytes, data...)
	expectedLength := checker.expectedLength(checker.cachebytes)
	if expectedLength != nil {
		checker.fixedheader = &MqttFixedHeader{}
		checker.fixedheader.With(checker.cachebytes[0], expectedLength)
	}
	if checker.fixedheader != nil {
		nLen := len(checker.cachebytes) - len(checker.fixedheader.length) - 1
		if nLen == checker.fixedheader.LessLength() {
			bts := make([]byte, len(checker.cachebytes))
			copy(bts, checker.cachebytes)
			checker.cachebytes = nil
			checker.fixedheader = nil
			return [][]byte{bts}
		} else {
			if nLen > checker.fixedheader.LessLength() {
				return checker.unpacking(registerId)
			}
		}
	}
	return nil
}

var DefaultMqttChecker = func() *xMqttDataChecker {
	return &xMqttDataChecker{}
}

type MqttConnect struct {
	ZeroSocketConnect

	topcis               map[string]byte
	messageSerialNnumber uint16
	serialNnumberMutex   sync.Mutex

	xListener MqttMessageListener
}

func NewMqttConnect() MqttConnect {
	return MqttConnect{topcis: make(map[string]byte)}
}

func NewMqttConnectPtr() *MqttConnect {
	return &MqttConnect{topcis: make(map[string]byte)}
}

func (mqttconn *MqttConnect) AddListener(xListener MqttMessageListener) {
	mqttconn.xListener = xListener
}

func (mqttconn *MqttConnect) RegisterId() string {
	return mqttconn.This().(ZeroConnect).RemoteAddr()
}

func (mqttconn *MqttConnect) Accept(_ ZeroServ, connect net.Conn) error {
	mqttconn.ZeroSocketConnect.Accept(global.Value(CORE_MQTT_SERVER).(*MqttServer), connect)
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

func NewMqttServer(address string, authWaitSeconds int64, heartbeatSeconds int64, bufferSize int) *MqttServer {
	return &MqttServer{
		TCPServer: *NewTCPServer(address, authWaitSeconds, heartbeatSeconds, bufferSize),
		topicsMap: make(map[string]map[string]*MqttConnect),
	}
}

func (mqttserv *MqttServer) RunServer() {
	if mqttserv.ConnectBuilder == nil {
		mqttserv.ConnectBuilder = &MqttConnectBuilder{}
	}
	global.Key(CORE_MQTT_SERVER, mqttserv)
	mqttserv.TCPServer.RunServer()
}
