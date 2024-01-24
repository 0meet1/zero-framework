package protocol

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/server"
)

type xZeroV1ClientListener struct{}

func (xListener *xZeroV1ClientListener) OnConnect(conn server.ZeroClientConnect) error {
	connectMessage, err := NewV1Message(MESSAGE_TYPE_CONNECT, make([]byte, 0))
	if err != nil {
		return err
	}
	err = connectMessage.Complete()
	if err != nil {
		return err
	}
	conn.(*xZeroV1Client).connectMessage = connectMessage

	<-time.After(time.Duration(time.Second * 1))
	err = conn.(*xZeroV1Client).Write(connectMessage.Bytes())
	if err != nil {
		return err
	}
	global.Logger().Debug(fmt.Sprintf("0protocol/1.0 client connect %s send connect message \n%s", conn.(*xZeroV1Client).RemoteAddr(), connectMessage.String()))
	return nil
}

func (xListener *xZeroV1ClientListener) OnHeartbeat(conn server.ZeroClientConnect) error {
	beatMessage, err := NewV1Message(MESSAGE_TYPE_HEARTBEAT, make([]byte, 0))
	if err != nil {
		return err
	}
	err = beatMessage.Complete()
	if err != nil {
		return err
	}

	err = conn.(*xZeroV1Client).PushMessage(beatMessage)
	if err != nil {
		return err
	}
	global.Logger().Debug(fmt.Sprintf("0protocol/1.0 client connect %s send heartbeat message \n%s", conn.(*xZeroV1Client).RemoteAddr(), beatMessage.String()))
	return nil
}

type xZeroV1Client struct {
	server.TCPClient

	operator ZeroV1MessageOperator

	request      *ZeroV1Message
	response     *ZeroV1Message
	requestMutex sync.Mutex
	responseChan chan *ZeroV1Message

	connectMessage *ZeroV1Message
}

func (client *xZeroV1Client) This() interface{} {
	if client.ZeroMeta.This() == nil {
		client.ThisDef(client)
	}
	return client.ZeroMeta.This()
}

func (client *xZeroV1Client) Active() bool {
	return client.connectMessage == nil
}

func (client *xZeroV1Client) ExecMessage(message *ZeroV1Message, withSecond int) (*ZeroV1Message, error) {
	if client.responseChan != nil {
		return nil, errors.New(fmt.Sprintf("0protocol/1.0 client connect %s is busying", client.RemoteAddr()))
	}
	err := client.Write(message.Bytes())
	if err != nil {
		return nil, err
	}

	client.requestMutex.Lock()
	client.request = message
	client.requestMutex.Unlock()
	client.responseChan = make(chan *ZeroV1Message, 1)
	select {
	case resp := <-client.responseChan:
		client.requestMutex.Lock()
		client.request = nil
		client.responseChan = nil
		client.requestMutex.Unlock()
		return resp, nil
	case <-time.After(time.Second * time.Duration(withSecond)):
		client.requestMutex.Lock()
		client.request = nil
		client.responseChan = nil
		client.requestMutex.Unlock()
		return nil, errors.New(" request timeout ")
	}
}

func (client *xZeroV1Client) PushMessage(message *ZeroV1Message) error {
	return client.Write(message.Bytes())
}

func (client *xZeroV1Client) OnMessage(datas []byte) error {
	client.TCPClient.OnMessage(datas)
	uMessage := ParseV1Message(datas)

	global.Logger().Debug(fmt.Sprintf("0protocol/1.0 client connect %s on message \n%s", client.RemoteAddr(), uMessage.String()))
	if uMessage.MessageType() == MESSAGE_TYPE_CONNACK {
		if uMessage.MessageId() == client.connectMessage.MessageId() {
			client.connectMessage = nil
		}
	} else if uMessage.MessageType() == MESSAGE_TYPE_BEATACK {
		client.Heartbeat()
	} else {
		callback := func() {
			if client.responseChan != nil {
				client.responseChan <- uMessage
			} else {
				global.Logger().Debug(fmt.Sprintf("0protocol/1.0 client connect %s ignore message \n%s", client.RemoteAddr(), uMessage.String()))
			}
		}
		if client.operator == nil {
			callback()
		} else {
			ok, err := client.operator.Operation(nil, uMessage)
			if err != nil {
				return err
			}
			if !ok {
				callback()
			}
		}
	}
	return nil
}

func (client *xZeroV1Client) Connect() {
	client.AddListener(&xZeroV1ClientListener{})
	client.AddChecker(&xZeroV1DataChecker{})
	client.TCPClient.Connect()
}

func RunZeroV1Client(addr string, heartbeatTime int, heartbeatCheckInterval int, operator ZeroV1MessageOperator) {
	zerov1cli := &xZeroV1Client{
		TCPClient: *server.NewTCPClient(
			addr,
			xDEFAULT_AUTH_WAIT,
			int64(heartbeatTime),
			int64(heartbeatCheckInterval),
			xDEFAULT_BUFFER_SIZE,
		),
		operator: operator,
	}
	zerov1cli.ThisDef(zerov1cli)
	global.Key(ZEROV1SERV_CLIENT, zerov1cli)
	zerov1cli.Connect()
}
