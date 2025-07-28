package server

import (
	"encoding/binary"
	"fmt"

	"github.com/0meet1/zero-framework/global"
)

const (
	CONNECT     = 0b0001
	CONNACK     = 0b0010
	PUBLISH     = 0b0011
	PUBACK      = 0b0100
	PUBREC      = 0b0101
	PUBREL      = 0b0110
	PUBCOMP     = 0b0111
	SUBSCRIBE   = 0b1000
	SUBACK      = 0b1001
	UNSUBSCRIBE = 0b1010
	UNSUBACK    = 0b1011
	PINGREQ     = 0b1100
	PINGRESP    = 0b1101
	DISCONNECT  = 0b1110

	LESS_LENGTH_LIMIT4 = 0x0FFFFFFF
	LESS_LENGTH_LIMIT3 = 0x001FFFFF
	LESS_LENGTH_LIMIT2 = 0x00003FFF
	LESS_LENGTH_LIMIT1 = 0x0000007F
)

type MqttFixedHeader struct {
	header byte
	length []byte
}

func (fixedHeader *MqttFixedHeader) varLTB(length int) []byte {
	if length > LESS_LENGTH_LIMIT4 {
		return make([]byte, 0)
	} else if length > LESS_LENGTH_LIMIT3 {
		bytes := make([]byte, 4)
		bytes[3] = byte(length << 4 & 0xFFFFFFFF >> 25 & 0xFF)
		bytes[2] = byte(length<<11&0xFFFFFFFF>>25&0xFF) + 0b10000000
		bytes[1] = byte(length<<18&0xFFFFFFFF>>25&0xFF) + 0b10000000
		bytes[0] = byte(length<<25&0xFFFFFFFF>>25&0xFF) + 0b10000000
		return bytes
	} else if length > LESS_LENGTH_LIMIT2 {
		bytes := make([]byte, 3)
		bytes[2] = byte(length << 11 & 0xFFFFFFFF >> 25 & 0xFF)
		bytes[1] = byte(length<<18&0xFFFFFFFF>>25&0xFF) + 0b10000000
		bytes[0] = byte(length<<25&0xFFFFFFFF>>25&0xFF) + 0b10000000
		return bytes
	} else if length > LESS_LENGTH_LIMIT1 {
		bytes := make([]byte, 2)
		bytes[1] = byte(length << 18 & 0xFFFFFFFF >> 25 & 0xFF)
		bytes[0] = byte(length<<25&0xFFFFFFFF>>25&0xFF) + 0b10000000
		return bytes
	} else {
		bytes := make([]byte, 1)
		bytes[0] = byte(length << 25 & 0xFFFFFFFF >> 25 & 0xFF)
		return bytes
	}
}

func (fixedHeader *MqttFixedHeader) BTvarL(bytes []byte) int {
	lessLength := 0
	byteLen := len(bytes)
	for i := 0; i < byteLen; i++ {
		lessLength += int(bytes[i]&0x7F) << (7 * i)
	}
	return lessLength
}

func (fixedHeader *MqttFixedHeader) With(header byte, length []byte) {
	fixedHeader.header = header
	fixedHeader.length = length
}

func (fixedHeader *MqttFixedHeader) make(messageType byte, messageTypeExt byte, length int) {
	fixedHeader.header = (messageType<<4)&0xF0 + (messageTypeExt & 0x0F)
	fixedHeader.length = fixedHeader.varLTB(length)
}

func (fixedHeader *MqttFixedHeader) MessageType() byte {
	return (fixedHeader.header >> 4) & 0x0f
}

func (fixedHeader *MqttFixedHeader) MessageTypeString() string {
	switch fixedHeader.MessageType() {
	case CONNECT:
		return "CONNECT"
	case CONNACK:
		return "CONNACK"
	case PUBLISH:
		return "PUBLISH"
	case PUBACK:
		return "PUBACK"
	case PUBREC:
		return "PUBREC"
	case PUBREL:
		return "PUBREL"
	case PUBCOMP:
		return "PUBCOMP"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	case SUBACK:
		return "SUBACK"
	case UNSUBSCRIBE:
		return "UNSUBSCRIBE"
	case UNSUBACK:
		return "UNSUBACK"
	case PINGREQ:
		return "PINGREQ"
	case PINGRESP:
		return "PINGRESP"
	case DISCONNECT:
		return "DISCONNECT"
	default:
		return "UNKNOW"
	}
}

func (fixedHeader *MqttFixedHeader) B3() byte {
	return fixedHeader.header << 4 & 0xFF >> 7 & 0xFF
}

func (fixedHeader *MqttFixedHeader) B2() byte {
	return fixedHeader.header << 5 & 0xFF >> 7 & 0xFF
}

func (fixedHeader *MqttFixedHeader) B1() byte {
	return fixedHeader.header << 6 & 0xFF >> 7 & 0xFF
}

func (fixedHeader *MqttFixedHeader) B0() byte {
	return fixedHeader.header << 7 & 0xFF >> 7 & 0xFF
}

func (fixedHeader *MqttFixedHeader) LessLength() int {
	return fixedHeader.BTvarL(fixedHeader.length)
}

func (fixedHeader *MqttFixedHeader) Length() []byte {
	return fixedHeader.length
}

const (
	MQTT_HEADER      = "MQTT"
	MQTT_LEVEL_3_1_1 = 0x04

	CONNECT_VARIABLE_HEADER_LEN = 10
)

type MqttCoreVariableHeader interface {
	VariableHeader() []byte
}

type MqttVariableHeader struct {
	variableHeader []byte
}

func (variableHeader *MqttVariableHeader) build(data []byte) error {
	variableHeader.variableHeader = data
	return nil
}

func (variableHeader *MqttVariableHeader) VariableHeader() []byte {
	return variableHeader.variableHeader
}

type MqttCorePayload interface {
	Payload() []byte
}

type MqttPayload struct {
	payload []byte
}

func (payload *MqttPayload) build(data []byte) error {
	payload.payload = data
	return nil
}

func (payload *MqttPayload) Payload() []byte {
	return payload.payload
}

type MqttMessage struct {
	fixedHeader    *MqttFixedHeader
	variableHeader MqttCoreVariableHeader
	payload        MqttCorePayload
}

func ParseMqttMessage(data []byte) (*MqttMessage, error) {
	m := &MqttMessage{}
	err := m.build(data)
	return m, err
}

func (message *MqttMessage) build(data []byte) error {
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("mqttcore build message err : %s", err))
		}
	}()

	lengthBytes := make([]byte, 0)
	i := 1
	for {
		if i >= len(data) {
			break
		}

		lengthBytes = append(lengthBytes, data[i])
		flag := data[i] & 0xFF >> 7 & 0xFF
		if flag == 0b0 || i >= 4 {
			break
		} else {
			i++
		}
	}

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.With(data[0], lengthBytes)

	if len(data)-len(message.fixedHeader.length)-1 != message.fixedHeader.LessLength() {
		return fmt.Errorf("message less length inconsistent real %d record %d", len(data), message.fixedHeader.LessLength())
	}

	fixedHeaderLen := len(message.fixedHeader.length) + 1

	switch message.fixedHeader.MessageType() {
	case CONNECT:
		variableHeader := &MqttConnectVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+10])
		message.variableHeader = variableHeader

		payload := &MqttParamsPayload{}
		payload.build(data[fixedHeaderLen+10:])
		message.payload = payload
	case CONNACK:
		variableHeader := &MqttConnackVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		message.payload = &MqttPayload{}
	case PUBLISH:
		variableHeader := &MqttPublishVariableHeader{}
		variableHeader.build(data[fixedHeaderLen:])
		message.variableHeader = variableHeader

		payload := &MqttPayload{}
		payload.build(data[fixedHeaderLen+len(variableHeader.variableHeader):])
		message.payload = payload
	case PUBACK:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		message.payload = &MqttPayload{}
	case PUBREC:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		message.payload = &MqttPayload{}
	case PUBREL:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		message.payload = &MqttPayload{}
	case PUBCOMP:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		message.payload = &MqttPayload{}
	case SUBSCRIBE:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		payload := &MqttSubscribePayload{}
		payload.build(data[fixedHeaderLen+len(variableHeader.variableHeader):])
		message.payload = payload
	case SUBACK:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		payload := &MqttPayload{}
		payload.build(data[fixedHeaderLen+2:])
		message.payload = payload
	case UNSUBSCRIBE:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		payload := &MqttParamsPayload{}
		payload.build(data[fixedHeaderLen+2:])
		message.payload = payload
	case UNSUBACK:
		variableHeader := &MqttIdentifierVariableHeader{}
		variableHeader.build(data[fixedHeaderLen : fixedHeaderLen+2])
		message.variableHeader = variableHeader

		message.payload = &MqttPayload{}
	case PINGREQ:
		message.variableHeader = &MqttVariableHeader{}
		message.payload = &MqttPayload{}
	case PINGRESP:
		message.variableHeader = &MqttVariableHeader{}
		message.payload = &MqttPayload{}
	case DISCONNECT:
		message.variableHeader = &MqttVariableHeader{}
		message.payload = &MqttPayload{}
	}

	return nil
}

func (message *MqttMessage) FixedHeader() *MqttFixedHeader {
	return message.fixedHeader
}

func (message *MqttMessage) VariableHeader() MqttCoreVariableHeader {
	return message.variableHeader
}

func (message *MqttMessage) Payload() MqttCorePayload {
	return message.payload
}

func (message *MqttMessage) MakeConnackMessage() {

	connack := &MqttConnackVariableHeader{}
	connack.make(0x00, 0x00)
	message.variableHeader = connack

	payload := &MqttPayload{}
	payload.build(make([]byte, 0))
	message.payload = payload

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.make(CONNACK, 0b0000, len(connack.variableHeader)+len(payload.payload))
}

func (message *MqttMessage) MakePingrespMessage() {

	pingresp := &MqttVariableHeader{}
	pingresp.build(make([]byte, 0))
	message.variableHeader = pingresp

	payload := &MqttPayload{}
	payload.build(make([]byte, 0))
	message.payload = payload

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.make(PINGRESP, 0b0000, len(pingresp.variableHeader)+len(payload.payload))
}

func (message *MqttMessage) MakeSubackMessage(identifier uint16, results []byte) {

	suback := &MqttIdentifierVariableHeader{}
	suback.make(identifier)
	message.variableHeader = suback

	payload := &MqttPayload{}
	payload.build(results)
	message.payload = payload

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.make(SUBACK, 0b0000, len(suback.variableHeader)+len(payload.payload))
}

func (message *MqttMessage) MakePubackMessage(identifier uint16) {

	puback := &MqttIdentifierVariableHeader{}
	puback.make(identifier)
	message.variableHeader = puback

	payload := &MqttPayload{}
	payload.build(make([]byte, 0))
	message.payload = payload

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.make(PUBACK, 0b0000, len(puback.variableHeader)+len(payload.payload))
}

func (message *MqttMessage) MakePubrelMessage(identifier uint16) {

	pubrel := &MqttIdentifierVariableHeader{}
	pubrel.make(identifier)
	message.variableHeader = pubrel

	payload := &MqttPayload{}
	payload.build(make([]byte, 0))
	message.payload = payload

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.make(PUBREL, 0b0000, len(pubrel.variableHeader)+len(payload.payload))
}

func (message *MqttMessage) MakePublistMessage(topic string, identifier uint16, data []byte) {

	publish := &MqttPublishVariableHeader{}
	publish.make(topic, int(identifier))
	message.variableHeader = publish

	payload := &MqttPayload{}
	payload.build(data)
	message.payload = payload

	message.fixedHeader = &MqttFixedHeader{}
	message.fixedHeader.make(PUBLISH, 0b0100, len(publish.variableHeader)+len(payload.payload))
}

func (message *MqttMessage) Bytes() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, message.fixedHeader.header)
	bytes = append(bytes, message.fixedHeader.length...)
	bytes = append(bytes, message.variableHeader.VariableHeader()...)
	bytes = append(bytes, message.payload.Payload()...)
	return bytes
}

type MqttIdentifierVariableHeader struct {
	MqttVariableHeader
}

func (identifierHeader *MqttIdentifierVariableHeader) build(data []byte) error {
	identifierHeader.MqttVariableHeader.build(data)
	if len(identifierHeader.MqttVariableHeader.variableHeader) != 2 {
		return fmt.Errorf("invalid connect variable header length : %d", len(identifierHeader.MqttVariableHeader.variableHeader))
	}
	return nil
}

func (identifierHeader *MqttIdentifierVariableHeader) make(identifier uint16) {
	identifierHeader.MqttVariableHeader.variableHeader = make([]byte, 2)
	binary.BigEndian.PutUint16(identifierHeader.MqttVariableHeader.variableHeader, uint16(identifier))
}

func (identifierHeader *MqttIdentifierVariableHeader) Identifier() uint16 {
	return binary.BigEndian.Uint16(identifierHeader.MqttVariableHeader.variableHeader)
}

type MqttConnectVariableHeader struct {
	MqttVariableHeader
}

func (connectHeader *MqttConnectVariableHeader) build(data []byte) error {
	connectHeader.MqttVariableHeader.build(data)
	if len(connectHeader.MqttVariableHeader.variableHeader) != 10 {
		return fmt.Errorf("invalid connect variable header length : %d", len(connectHeader.MqttVariableHeader.variableHeader))
	}

	if connectHeader.ProtocolLength() != 4 {
		return fmt.Errorf("invalid connect variable protocol length : %d", connectHeader.ProtocolLength())
	}

	if connectHeader.Protocol() != MQTT_HEADER {
		return fmt.Errorf("invalid connect variable protocol : %s", connectHeader.Protocol())
	}

	return nil
}

// func (connectHeader *MqttConnectVariableHeader) make(
// 	level byte,
// 	UserNameFlag byte,
// 	PasswordFlag byte,
// 	WillRetain byte,
// 	WillQos byte,
// 	WillFlag byte,
// 	CleanSession byte,
// 	KeepAlive int,
// ) {
// 	protocolLength := make([]byte, 8)
// 	protocolLength[0] = 0x00
// 	protocolLength[1] = 0x04
// 	protocolLength[2] = 'M'
// 	protocolLength[3] = 'Q'
// 	protocolLength[4] = 'T'
// 	protocolLength[5] = 'T'
// 	protocolLength[6] = level
// 	protocolLength[7] = 0x00

// 	keepAliveBuf := make([]byte, 2)
// 	binary.BigEndian.PutUint16(keepAliveBuf, uint16(KeepAlive))
// 	protocolLength = append(protocolLength, keepAliveBuf...)

// 	protocolLength[7] += UserNameFlag
// 	protocolLength[7] <<= 1
// 	protocolLength[7] += PasswordFlag
// 	protocolLength[7] <<= 1
// 	protocolLength[7] += WillRetain
// 	protocolLength[7] <<= 2
// 	protocolLength[7] += WillQos
// 	protocolLength[7] <<= 1
// 	protocolLength[7] += WillFlag
// 	protocolLength[7] <<= 1
// 	protocolLength[7] += CleanSession
// 	protocolLength[7] <<= 1
// 	protocolLength[7] += 0b0
// }

func (connectHeader *MqttConnectVariableHeader) ProtocolLength() int {
	return int(binary.BigEndian.Uint16(connectHeader.MqttVariableHeader.variableHeader[:2]))
}

func (connectHeader *MqttConnectVariableHeader) Protocol() string {
	return string(connectHeader.MqttVariableHeader.variableHeader[2:6])
}

func (connectHeader *MqttConnectVariableHeader) Level() byte {
	return connectHeader.MqttVariableHeader.variableHeader[6]
}

func (connectHeader *MqttConnectVariableHeader) UserNameFlag() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] >> 7 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) PasswordFlag() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] << 1 & 0xFF >> 7 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) WillRetain() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] << 2 & 0xFF >> 7 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) WillQos() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] << 3 & 0xFF >> 6 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) WillFlag() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] << 5 & 0xFF >> 7 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) CleanSession() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] << 6 & 0xFF >> 7 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) Reserved() byte {
	return connectHeader.MqttVariableHeader.variableHeader[7] << 7 & 0xFF >> 7 & 0xFF
}

func (connectHeader *MqttConnectVariableHeader) KeepAlive() int {
	return int(binary.BigEndian.Uint16(connectHeader.MqttVariableHeader.variableHeader[8:]))
}

type MqttConnackVariableHeader struct {
	MqttVariableHeader
}

func (connackHeader *MqttConnackVariableHeader) build(data []byte) error {
	connackHeader.MqttVariableHeader.build(data)
	if len(connackHeader.MqttVariableHeader.variableHeader) != 2 {
		return fmt.Errorf("invalid connack variable header length : %d", len(connackHeader.MqttVariableHeader.variableHeader))
	}
	return nil
}

func (connackHeader *MqttConnackVariableHeader) make(sessionPresent byte, returnCode byte) {
	connackHeader.variableHeader = make([]byte, 2)
	connackHeader.variableHeader[0] = sessionPresent
	connackHeader.variableHeader[1] = returnCode
}

func (connackHeader *MqttConnackVariableHeader) SessionPresent() byte {
	return connackHeader.variableHeader[0]
}

func (connackHeader *MqttConnackVariableHeader) ReturnCode() byte {
	return connackHeader.variableHeader[1]
}

type MqttPublishVariableHeader struct {
	MqttVariableHeader
}

func (publishVariableHeader *MqttPublishVariableHeader) build(data []byte) error {
	topicLen := int(binary.BigEndian.Uint16(data[:2]))
	publishVariableHeader.MqttVariableHeader.build(data[:2+topicLen+2])
	return nil
}

func (publishVariableHeader *MqttPublishVariableHeader) make(topic string, identifier int) {
	topicLen := len(topic)
	topicLenbytes := make([]byte, 2)
	binary.BigEndian.PutUint16(topicLenbytes, uint16(topicLen))
	publishVariableHeader.variableHeader = make([]byte, 0)
	publishVariableHeader.variableHeader = append(publishVariableHeader.variableHeader, topicLenbytes...)
	publishVariableHeader.variableHeader = append(publishVariableHeader.variableHeader, []byte(topic)...)

	identifierLenbytes := make([]byte, 2)
	binary.BigEndian.PutUint16(identifierLenbytes, uint16(identifier))
	publishVariableHeader.variableHeader = append(publishVariableHeader.variableHeader, identifierLenbytes...)
}

func (publishVariableHeader *MqttPublishVariableHeader) Topic() string {
	return string(publishVariableHeader.variableHeader[2 : len(publishVariableHeader.variableHeader)-2])
}

func (publishVariableHeader *MqttPublishVariableHeader) Identifier() uint16 {
	return binary.BigEndian.Uint16(publishVariableHeader.variableHeader[len(publishVariableHeader.variableHeader)-2:])
}

type MqttParamsPayload struct {
	MqttPayload
	params []string
}

func (connectPayload *MqttParamsPayload) build(data []byte) error {
	connectPayload.MqttPayload.build(data)
	connectPayload.params = make([]string, 0)
	index := 0
	for {
		if index >= len(data) {
			break
		}
		bodyLenBytes := connectPayload.payload[index : index+2]
		bodyLen := int(binary.BigEndian.Uint16(bodyLenBytes))
		index += 2
		connectPayload.params = append(connectPayload.params, string(connectPayload.payload[index:index+bodyLen]))
		index += bodyLen
	}
	return nil
}

// func (connectPayload *MqttParamsPayload) make(params ...string) error {
// 	connectPayload.params = params
// 	connectPayload.payload = make([]byte, 0)

// 	for i := 0; i < len(params); i++ {
// 		paramLenBuf := make([]byte, 2)
// 		binary.BigEndian.PutUint16(paramLenBuf, uint16(len(params[i])))
// 		connectPayload.payload = append(connectPayload.payload, paramLenBuf...)
// 		connectPayload.payload = append(connectPayload.payload, string(params[i])...)
// 	}

// 	return nil
// }

func (connectPayload *MqttParamsPayload) Params() []string {
	return connectPayload.params
}

type MqttTopic struct {
	TopicName string
	Qos       byte
}

type MqttSubscribePayload struct {
	MqttPayload
	topics []*MqttTopic
}

func (subscribePayload *MqttSubscribePayload) build(data []byte) error {
	subscribePayload.MqttPayload.build(data)
	subscribePayload.topics = make([]*MqttTopic, 0)
	index := 0
	for {
		if index >= len(data) {
			break
		}
		bodyLenBytes := subscribePayload.payload[index:2]
		bodyLen := int(binary.BigEndian.Uint16(bodyLenBytes))
		index += 2
		topic := string(subscribePayload.payload[index : index+bodyLen])
		index += bodyLen
		qos := subscribePayload.payload[index]
		index += 1
		subscribePayload.topics = append(subscribePayload.topics, &MqttTopic{
			TopicName: topic,
			Qos:       qos,
		})
	}
	return nil
}

// func (subscribePayload *MqttSubscribePayload) make(topics []*MqttTopic) error {
// 	subscribePayload.topics = topics
// 	subscribePayload.payload = make([]byte, 0)

// 	for i := 0; i < len(topics); i++ {
// 		paramLenBuf := make([]byte, 2)
// 		binary.BigEndian.PutUint16(paramLenBuf, uint16(len(topics[i].TopicName)))
// 		subscribePayload.payload = append(subscribePayload.payload, paramLenBuf...)
// 		subscribePayload.payload = append(subscribePayload.payload, string(topics[i].TopicName)...)
// 		subscribePayload.payload = append(subscribePayload.payload, topics[i].Qos)
// 	}

// 	return nil
// }

func (subscribePayload *MqttSubscribePayload) Topics() []*MqttTopic {
	return subscribePayload.topics
}
