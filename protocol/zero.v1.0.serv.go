package protocol

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/server"
)

const (
	xDEFAULT_AUTH_WAIT   = 10
	xDEFAULT_BUFFER_SIZE = 8 * 1024 * 1024
)

type xZeroV1ConnectBuilder struct{}

func (xDefault *xZeroV1ConnectBuilder) NewConnect() server.ZeroConnect {
	tcpconn := &xZeroV1Connect{
		keeper: global.Value(ZEROV1SERV_KEEPER).(*xZeroV1ServKeeper),
	}
	tcpconn.ThisDef(tcpconn)
	tcpconn.AddChecker(&xZeroV1DataChecker{})
	return tcpconn
}

type xZeroV1DataChecker struct {
	cachebytes      []byte
	cachebytesMutex sync.Mutex
}

func (checker *xZeroV1DataChecker) CheckPackageData(data []byte) []byte {
	checker.cachebytesMutex.Lock()
	defer func() {
		checker.cachebytesMutex.Unlock()
		err := recover()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("zerov1 on check package data err : %s", err))
		}
	}()

	if len(data) < 4 {
		if checker.cachebytes != nil {
			checker.cachebytes = append(checker.cachebytes, data...)
		}
	} else if reflect.DeepEqual(data[:4], xZERO_MESSAGE_HEAD) {
		checker.cachebytes = make([]byte, 0)
		checker.cachebytes = append(checker.cachebytes, data...)
	} else if checker.cachebytes != nil {
		checker.cachebytes = append(checker.cachebytes, data...)
	}

	if len(checker.cachebytes) > 0 && reflect.DeepEqual(checker.cachebytes[len(checker.cachebytes)-2:], xZERO_MESSAGE_END) {
		bts := make([]byte, len(checker.cachebytes))
		copy(bts, checker.cachebytes)
		v1msg := ParseV1Message(bts)
		err := v1msg.Check()
		if err != nil {
			global.Logger().Debug(err.Error())
			return nil
		}
		checker.cachebytes = nil
		return bts
	}
	return nil
}

type xZeroV1Connect struct {
	server.ZeroSocketConnect

	keeper *xZeroV1ServKeeper

	request      *ZeroV1Message
	response     *ZeroV1Message
	requestMutex sync.Mutex
	responseChan chan *ZeroV1Message
}

func (v1conn *xZeroV1Connect) Authorized(datas ...byte) bool {
	authMessage := ParseV1Message(datas)

	global.Logger().Info(fmt.Sprintf("zerov1 connect %s authorized", v1conn.RemoteAddr()))
	v1conn.Heartbeat()
	v1conn.ZeroSocketConnect.Authorized()

	ackMessage := NewV1AckMessage(MESSAGE_TYPE_CONNACK, authMessage.MessageId(), make([]byte, 0))

	err := ackMessage.Complete()
	if err != nil {
		global.Logger().Error(fmt.Sprintf("%s", err.Error()))
		return false
	}

	err = v1conn.pushMessage(ackMessage)
	if err != nil {
		global.Logger().Error(fmt.Sprintf("%s", err.Error()))
		return false
	}

	conn, err := v1conn.keeper.UseConnect(v1conn.RegisterId())
	if err == nil && conn != nil {
		conn.Close()
	}

	return true
}

func (v1conn *xZeroV1Connect) Close() error {
	return v1conn.ZeroSocketConnect.Close()
}

func (v1conn *xZeroV1Connect) RegisterId() string {
	return v1conn.RemoteAddr()
}

func (v1conn *xZeroV1Connect) execMessage(message *ZeroV1Message, withSecond int) (*ZeroV1Message, error) {
	if v1conn.responseChan != nil {
		return nil, errors.New(fmt.Sprintf("zerov1 connect %s is busying", v1conn.RemoteAddr()))
	}
	err := v1conn.Write(message.Bytes())
	if err != nil {
		return nil, err
	}
	v1conn.requestMutex.Lock()
	v1conn.request = message
	v1conn.requestMutex.Unlock()
	v1conn.responseChan = make(chan *ZeroV1Message, 1)
	select {
	case resp := <-v1conn.responseChan:
		v1conn.requestMutex.Lock()
		v1conn.request = nil
		v1conn.responseChan = nil
		v1conn.requestMutex.Unlock()
		return resp, nil
	case <-time.After(time.Second * time.Duration(withSecond)):
		v1conn.requestMutex.Lock()
		v1conn.request = nil
		v1conn.responseChan = nil
		v1conn.requestMutex.Unlock()
		return nil, errors.New(" request timeout ")
	}
}

func (v1conn *xZeroV1Connect) pushMessage(message *ZeroV1Message) error {
	return v1conn.Write(message.Bytes())
}

func (v1conn *xZeroV1Connect) OnMessage(datas []byte) error {
	uMessage := ParseV1Message(datas)

	global.Logger().Debug(fmt.Sprintf("zerov1 connect %s on message \n%s", v1conn.RegisterId(), uMessage.String()))
	if uMessage.MessageType() == MESSAGE_TYPE_CONNECT {
		if !v1conn.Authorized(datas...) {
			v1conn.Close()
		}
	} else if uMessage.MessageType() == MESSAGE_TYPE_HEARTBEAT {
		v1conn.Heartbeat()
		beatack := NewV1AckMessage(MESSAGE_TYPE_BEATACK, uMessage.MessageId(), make([]byte, 0))

		err := beatack.Complete()
		if err != nil {
			return err
		}

		err = v1conn.pushMessage(beatack)
		if err != nil {
			return err
		}
	} else {
		callback := func() {
			if v1conn.responseChan != nil {
				v1conn.responseChan <- uMessage
			} else {
				global.Logger().Debug(fmt.Sprintf("zerov1 connect %s ignore message \n%s", v1conn.RegisterId(), uMessage.String()))
			}
		}
		if v1conn.keeper.operator == nil {
			callback()
		} else {
			ok, err := v1conn.keeper.operator.Operation(uMessage)
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

type xZeroV1ServKeeper struct {
	server.TCPServer

	operator ZeroV1MessageOperator
}

func (keeper *xZeroV1ServKeeper) ExecMessage(registerId string, message *ZeroV1Message, withSecond int) (*ZeroV1Message, error) {
	conn, err := keeper.UseConnect(registerId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("use connect `%s` error: %s", registerId, err.Error()))
	}
	return conn.(*xZeroV1Connect).execMessage(message, withSecond)
}

func (keeper *xZeroV1ServKeeper) PushMessage(registerId string, message *ZeroV1Message) error {
	conn, err := keeper.UseConnect(registerId)
	if err != nil {
		return errors.New(fmt.Sprintf("use connect `%s` error: %s", registerId, err.Error()))
	}
	return conn.(*xZeroV1Connect).pushMessage(message)
}

func (keeper *xZeroV1ServKeeper) RunServer() {
	keeper.ConnectBuilder = &xZeroV1ConnectBuilder{}
	keeper.TCPServer.RunServer()
}

func RunZeroV1Server(addr string, heartbeatTime int, heartbeatCheckInterval int, operator ZeroV1MessageOperator) {
	zerov1serv := &xZeroV1ServKeeper{
		TCPServer: *server.NewTCPServer(
			addr,
			xDEFAULT_AUTH_WAIT,
			int64(heartbeatTime),
			int64(heartbeatCheckInterval),
			xDEFAULT_BUFFER_SIZE,
		),
		operator: operator,
	}
	global.Key(ZEROV1SERV_KEEPER, zerov1serv)
	zerov1serv.RunServer()
}
