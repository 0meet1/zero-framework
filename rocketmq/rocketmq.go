package rocketmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"zero-framework/global"
	"zero-framework/structs"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gofrs/uuid"
)

const (
	ROCKETMQ_KEEPER = "zero.rocketmq.keeper"
)

var (
	notifyPushConsumer rocketmq.PushConsumer
	notifyProducer     rocketmq.Producer

	nameserv    string
	groupName   string
	topics      []string
	testMessage string

	observerMutex sync.RWMutex
	observers     map[string]MQMessageObserver
)

const (
	ENABLE  = "enable"
	DISABLE = "disable"

	MESSAGE_TYPE_TEST = "notify.test"

	NOTIFY_ONLINE  = "device.online"
	NOTIFY_OFFLINE = "device.offline"
	NOTIFY_NORMAL  = "message.normal"
)

func makeTestMessage(topic string) *MQNotifyMessage {
	uid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return &MQNotifyMessage{
		MessageId:   uid.String(),
		Topic:       topic,
		CreateTime:  structs.Date(time.Now()),
		MessageType: MESSAGE_TYPE_TEST,
		Payload:     MESSAGE_TYPE_TEST,
	}
}

func InitRocketMQ(newObservers ...MQMessageObserver) {
	testMessage = ENABLE
	nameserv = global.StringValue("zero.rocketmq.nameserv")
	groupName = global.StringValue("zero.rocketmq.groupname")
	topics = global.SliceStringValue("zero.rocketmq.topics")
	testMessage = global.StringValue("zero.rocketmq.testmessage")

	observers = make(map[string]MQMessageObserver)
	observerMutex.Lock()
	defer observerMutex.Unlock()
	for _, obs := range newObservers {
		_, ok := observers[obs.Name()]
		if ok {
			panic(fmt.Sprintf("mqobserver '%s' is already exists", obs.Name()))
		}
		observers[obs.Name()] = obs
	}

	uid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	addr, err := primitive.NewNamesrvAddr(nameserv)
	if err != nil {
		panic(err)
	}
	notifyProducer, err = rocketmq.NewProducer(
		producer.WithGroupName(uid.String()),
		producer.WithNameServer(addr),
		producer.WithRetry(1),
	)
	if err != nil {
		panic(err)
	}

	err = notifyProducer.Start()
	if err != nil {
		panic(err)
	}

	global.Key(ROCKETMQ_KEEPER, &RocketmqKeeper{})
	if testMessage == ENABLE {
		<-time.After(time.Second * time.Duration(2))
		for _, topic := range topics {
			global.Value(ROCKETMQ_KEEPER).(*RocketmqKeeper).SendMessage(makeTestMessage(topic))
		}
	}

	<-time.After(time.Second * time.Duration(1))
	initRocketMQConsumer()
}

func initRocketMQConsumer() {
	notifyPushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{nameserv})),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset))
	if err != nil {
		panic(err)
	}

	for _, topic := range topics {
		err = notifyPushConsumer.Subscribe(topic, consumer.MessageSelector{}, onMessage)
		if err != nil {
			panic(err)
		}
		err = notifyPushConsumer.Start()
		if err != nil {
			panic(err)
		}
	}
}

func ShutdownNotifyConsumer() {
	err := notifyPushConsumer.Shutdown()
	if err != nil {
		panic(err)
	}
}

func onMessage(_ context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, message := range messages {
		notify := MQNotifyMessage{}
		notify.WithJSONString(message.Body)
		switch notify.MessageType {
		case MESSAGE_TYPE_TEST:
			global.Logger().Info(fmt.Sprintf(" topic : %s test success", message.Topic))
		default:
			observerMutex.RLock()
			if observers != nil {
				for _, obs := range observers {
					err := obs.OnMessage(&notify)
					if err != nil {
						global.Logger().Error(fmt.Sprintf(" message observer error %s", err.Error()))
					}
				}
			}
			observerMutex.RUnlock()
			global.Logger().Debug(fmt.Sprintf(" topic : %s on message %s", message.Topic, string(message.Body)))
		}
	}
	return consumer.ConsumeSuccess, nil
}

type RocketmqKeeper struct{}

func (keeper *RocketmqKeeper) SendMessage(message *MQNotifyMessage) error {
	jsonBytes, err := message.JSONString()
	if err != nil {
		return err
	}
	_, err = notifyProducer.SendSync(context.Background(), &primitive.Message{
		Topic: message.Topic,
		Body:  jsonBytes,
	})
	if err != nil {
		return err
	}
	return nil
}

func (keeper *RocketmqKeeper) NewMessage(topic string, messageType string, payload interface{}) (*MQNotifyMessage, error) {
	message := &MQNotifyMessage{}
	err := message.NewMessage(topic, messageType, payload)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (keeper *RocketmqKeeper) AddObservers(newObservers ...MQMessageObserver) error {
	observerMutex.Lock()
	defer observerMutex.Unlock()
	for _, obs := range newObservers {
		_, ok := observers[obs.Name()]
		if ok {
			return errors.New(fmt.Sprintf("mqobserver '%s' is already exists", obs.Name()))
		}
	}
	for _, obs := range newObservers {
		observers[obs.Name()] = obs
	}
	return nil
}

func (keeper *RocketmqKeeper) RemoveObservers(observerNames ...string) error {
	observerMutex.Lock()
	defer observerMutex.Unlock()
	for _, obsName := range observerNames {
		delete(observers, obsName)
	}
	return nil
}
