package rocketmq

import (
	"encoding/json"
	"time"
	"zero-framework/structs"

	"github.com/gofrs/uuid"
)

type MQMessageObserver interface {
	Name() string
	OnMessage(*MQNotifyMessage) error
}

type MQNotifyMessage struct {
	MessageId   string       `json:"messageId,omitempty"`
	Topic       string       `json:"topic,omitempty"`
	CreateTime  structs.Date `json:"createTime,omitempty"`
	MessageType string       `json:"messageType,omitempty"`
	Payload     interface{}  `json:"payload,omitempty"`
}

func (notify *MQNotifyMessage) NewMessage(topic string, messageType string, payload interface{}) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	notify.MessageId = uid.String()
	notify.Topic = topic
	notify.CreateTime = structs.Date(time.Now())
	notify.MessageType = messageType
	notify.Payload = payload

	return nil
}

func (notify *MQNotifyMessage) JSONString() ([]byte, error) {
	return json.Marshal(notify)
}

func (notify *MQNotifyMessage) WithJSONString(jsonbytes []byte) error {
	return json.Unmarshal(jsonbytes, &notify)
}
