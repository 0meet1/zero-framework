package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/0meet1/zero-framework/structs"
	"github.com/gofrs/uuid"
)

var (
	xV1_VERSION = []byte{0x00, 0x01}
)

const (
	MESSAGE_TYPE_CONNECT   = 0x01
	MESSAGE_TYPE_HEARTBEAT = 0x02

	MESSAGE_TYPE_CONNACK = 0x11
	MESSAGE_TYPE_BEATACK = 0x12
)

type ZeroV1Message struct {
	head        []byte
	version     []byte
	dataLength  []byte
	messageId   []byte
	messageType byte
	bodyLength  []byte
	messageBody []byte
	checkSum    []byte
	end         []byte
}

func NewV1Message(messageType byte, xBody []byte) (*ZeroV1Message, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &ZeroV1Message{
		head:        xZERO_MESSAGE_HEAD,
		version:     xV1_VERSION,
		dataLength:  []byte{0x00, 0x00, 0x00, 0x00},
		messageId:   []byte(strings.ReplaceAll(uid.String(), "-", "")),
		messageType: messageType,
		bodyLength:  []byte{0x00, 0x00, 0x00, 0x00},
		messageBody: xBody,
		checkSum:    []byte{0x00, 0x00},
		end:         xZERO_MESSAGE_END,
	}, nil
}

func NewV1AckMessage(messageType byte, messageId string, xBody []byte) *ZeroV1Message {
	return &ZeroV1Message{
		head:        xZERO_MESSAGE_HEAD,
		version:     xV1_VERSION,
		dataLength:  []byte{0x00, 0x00, 0x00, 0x00},
		messageId:   []byte(messageId),
		messageType: messageType,
		bodyLength:  []byte{0x00, 0x00, 0x00, 0x00},
		messageBody: xBody,
		checkSum:    []byte{0x00, 0x00},
		end:         xZERO_MESSAGE_END,
	}
}

func ParseV1Message(datas []byte) *ZeroV1Message {
	xDatasLength := len(datas)
	return &ZeroV1Message{
		head:        datas[0:4],
		version:     datas[4:6],
		dataLength:  datas[6:10],
		messageId:   datas[10:42],
		messageType: datas[42],
		bodyLength:  datas[43:47],
		messageBody: datas[47 : xDatasLength-6],
		checkSum:    datas[xDatasLength-6 : xDatasLength-4],
		end:         datas[xDatasLength-4:],
	}
}

func (v1msg *ZeroV1Message) Head() []byte {
	return v1msg.head
}

func (v1msg *ZeroV1Message) HeadString() string {
	return string(v1msg.head)
}

func (v1msg *ZeroV1Message) End() []byte {
	return v1msg.end
}

func (v1msg *ZeroV1Message) EndString() string {
	return string(v1msg.end)
}

func (v1msg *ZeroV1Message) MessageId() string {
	return string(v1msg.messageId)
}

func (v1msg *ZeroV1Message) Version() int {
	return int(binary.BigEndian.Uint16(v1msg.version))
}

func (v1msg *ZeroV1Message) DataLength() int {
	return int(binary.BigEndian.Uint32(v1msg.dataLength))
}

func (v1msg *ZeroV1Message) MessageType() int {
	return int(v1msg.messageType)
}

func (v1msg *ZeroV1Message) BodyLength() int {
	return int(binary.BigEndian.Uint32(v1msg.bodyLength))
}

func (v1msg *ZeroV1Message) MessageBody() []byte {
	return v1msg.messageBody
}

func (v1msg *ZeroV1Message) Complete() error {
	binary.BigEndian.PutUint32(v1msg.bodyLength, uint32(len(v1msg.messageBody)))
	binary.BigEndian.PutUint32(v1msg.dataLength, uint32(21+len(v1msg.messageBody)))

	bodys := make([]byte, 0)
	bodys = append(bodys, v1msg.head...)
	bodys = append(bodys, v1msg.version...)
	bodys = append(bodys, v1msg.dataLength...)
	bodys = append(bodys, v1msg.messageId...)
	bodys = append(bodys, v1msg.messageType)
	bodys = append(bodys, v1msg.bodyLength...)
	bodys = append(bodys, v1msg.messageBody...)
	bodys = append(bodys, 0x00, 0x00)
	bodys = append(bodys, v1msg.end...)

	crc16code := structs.NewCRC16Table(structs.CRC16_AUG_CCITT).Complete(bodys)
	binary.BigEndian.PutUint16(v1msg.checkSum, crc16code)

	return nil
}

func (v1msg *ZeroV1Message) Check() error {
	bodys := make([]byte, 0)
	bodys = append(bodys, v1msg.head...)
	bodys = append(bodys, v1msg.version...)
	bodys = append(bodys, v1msg.dataLength...)
	bodys = append(bodys, v1msg.messageId...)
	bodys = append(bodys, v1msg.messageType)
	bodys = append(bodys, v1msg.bodyLength...)
	bodys = append(bodys, v1msg.messageBody...)
	bodys = append(bodys, 0x00, 0x00)
	bodys = append(bodys, v1msg.end...)

	if !reflect.DeepEqual(xZERO_MESSAGE_HEAD, v1msg.head) {
		return errors.New(fmt.Sprintf("\n### err message head %s ### message datas \n%s", structs.BytesString(v1msg.head...), structs.BytesString(bodys...)))
	}

	if !reflect.DeepEqual(xZERO_MESSAGE_END, v1msg.end) {
		return errors.New(fmt.Sprintf("\n### err message end %s ### message datas \n%s", structs.BytesString(v1msg.end...), structs.BytesString(bodys...)))
	}

	if v1msg.dataLength == nil || v1msg.DataLength() != len(bodys) {
		return errors.New(fmt.Sprintf("\n### err message data length %d reality %d ### message datas \n%s",
			v1msg.DataLength(),
			len(bodys),
			structs.BytesString(bodys...)))
	}

	if v1msg.bodyLength == nil || v1msg.BodyLength() != len(v1msg.messageBody) {
		return errors.New(fmt.Sprintf("\n### err message body length %d reality %d ### message datas \n%s",
			v1msg.BodyLength(),
			len(v1msg.messageBody),
			structs.BytesString(v1msg.messageBody...)))
	}

	crc16code := structs.NewCRC16Table(structs.CRC16_AUG_CCITT).Complete(bodys)
	crc16bin := make([]byte, 2)
	binary.BigEndian.PutUint16(crc16bin, crc16code)

	if !reflect.DeepEqual(crc16bin, v1msg.checkSum) {
		return errors.New(fmt.Sprintf("\n### err message verify %s ### message datas \n%s", structs.BytesString(v1msg.checkSum...), structs.BytesString(bodys...)))
	}

	return nil
}

func (v1msg *ZeroV1Message) Bytes() []byte {
	bodys := make([]byte, 0)
	bodys = append(bodys, v1msg.head...)
	bodys = append(bodys, v1msg.version...)
	bodys = append(bodys, v1msg.dataLength...)
	bodys = append(bodys, v1msg.messageId...)
	bodys = append(bodys, v1msg.messageType)
	bodys = append(bodys, v1msg.bodyLength...)
	bodys = append(bodys, v1msg.messageBody...)
	bodys = append(bodys, v1msg.checkSum...)
	bodys = append(bodys, v1msg.end...)
	return bodys
}

func (v1msg *ZeroV1Message) String() string {
	return structs.BytesString(v1msg.Bytes()...)
}
