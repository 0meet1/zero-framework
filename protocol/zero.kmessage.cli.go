package protocol

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/server"
)

type kZeroKMessageClientListener struct {
	uniquekey string
}

func (xListener *kZeroKMessageClientListener) OnConnect(conn server.ZeroClientConnect) error {
	cMessage, err := NewKMessage(MESSAGE_TYPE_CONNECT, make([]byte, 0))
	if err != nil {
		return err
	}
	cMessage.AddUniqueKey(xListener.uniquekey)
	err = cMessage.Complete()
	if err != nil {
		return err
	}
	conn.(*kZeroKMessageClient).connectMessage = cMessage

	<-time.After(time.Duration(time.Second * 1))
	err = conn.(*kZeroKMessageClient).Write(cMessage.Bytes())
	if err != nil {
		return err
	}
	global.Logger().Debug(fmt.Sprintf("0protocol/1.0 client connect %s send connect message \n%s", conn.(*kZeroKMessageClient).RemoteAddr(), cMessage.String()))
	return nil
}

func (xListener *kZeroKMessageClientListener) OnHeartbeat(conn server.ZeroClientConnect) error {
	beatMessage, err := NewKMessage(MESSAGE_TYPE_HEARTBEAT, make([]byte, 0))
	if err != nil {
		return err
	}
	beatMessage.AddUniqueKey(xListener.uniquekey)
	err = beatMessage.Complete()
	if err != nil {
		return err
	}

	err = conn.(*kZeroKMessageClient).PushMessage(beatMessage)
	if err != nil {
		return err
	}
	global.Logger().Debug(fmt.Sprintf("0protocol/1.0 client connect %s send heartbeat message \n%s", conn.(*kZeroKMessageClient).RemoteAddr(), beatMessage.String()))
	return nil
}

type kZeroKMessageClient struct {
	server.TCPClient

	uniquekey string
	operator  ZeroKMessageOperator

	request      *ZeroKMessage
	requestMutex sync.Mutex
	responseChan chan *ZeroKMessage

	connectMessage *ZeroKMessage
}

func (client *kZeroKMessageClient) This() interface{} {
	if client.ZeroMeta.This() == nil {
		client.ThisDef(client)
	}
	return client.ZeroMeta.This()
}

func (client *kZeroKMessageClient) Active() bool {
	return client.connectMessage == nil
}

func (client *kZeroKMessageClient) ExecMessage(message *ZeroKMessage, withSecond int) (*ZeroKMessage, error) {
	if client.responseChan != nil {
		return nil, fmt.Errorf("0protocol/1.0 client connect %s is busying", client.RemoteAddr())
	}
	err := client.Write(message.Bytes())
	if err != nil {
		return nil, err
	}

	client.requestMutex.Lock()
	client.request = message
	client.requestMutex.Unlock()
	client.responseChan = make(chan *ZeroKMessage, 1)
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

func (client *kZeroKMessageClient) PushMessage(message *ZeroKMessage) error {
	return client.Write(message.Bytes())
}

func (client *kZeroKMessageClient) OnMessage(datas []byte) error {
	client.TCPClient.OnMessage(datas)
	uMessage := ParseKMessage(datas)

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

func (client *kZeroKMessageClient) Connect() {
	client.AddListener(&kZeroKMessageClientListener{uniquekey: client.uniquekey})
	client.AddChecker(&kZeroKMessageChecker{})
	client.TCPClient.Connect()
}

func RunKMessageClient(addr string, heartbeatTime int, heartbeatCheckInterval int, operator ZeroKMessageOperator, unk ...string) {
	_uniquekey := ""
	if len(unk) > 0 {
		_uniquekey = unk[0]
	}
	kMessageCli := &kZeroKMessageClient{
		TCPClient: *server.NewTCPClient(
			addr,
			xDEFAULT_AUTH_WAIT,
			int64(heartbeatTime),
			int64(heartbeatCheckInterval),
			xDEFAULT_BUFFER_SIZE,
		),
		uniquekey: _uniquekey,
		operator:  operator,
	}
	kMessageCli.ThisDef(kMessageCli)
	global.Key(ZEROKMSG_CLIENT, kMessageCli)
	kMessageCli.Connect()
}
