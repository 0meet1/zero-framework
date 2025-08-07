package protocol

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

const (
	xDEFAULT_AUTH_WAIT   = 10
	xDEFAULT_BUFFER_SIZE = 8 * 1024 * 1024
)

type xZeroKMessageConnectBuilder struct{}

func (xDefault *xZeroKMessageConnectBuilder) NewConnect() server.ZeroConnect {
	tcpconn := &kZeroKMessageConnect{
		keeper: global.Value(ZEROKMSG_SERVER).(*kZeroKMessageKeeper),
	}
	tcpconn.ThisDef(tcpconn)
	tcpconn.AddChecker(&kZeroKMessageChecker{})
	return tcpconn
}

type kZeroKMessageChecker struct {
	cachebytes      []byte
	cachebytesMutex sync.Mutex
}

func (checker *kZeroKMessageChecker) unpacking(registerId string, historys ...[]byte) [][]byte {
	comps := make([][]byte, 0)
	comps = append(comps, historys...)

	dataLength := ParseKMessageLength(checker.cachebytes)
	if dataLength < 0 {
		return comps
	}

	if len(checker.cachebytes) >= dataLength {
		comps = append(comps, checker.cachebytes[:dataLength])
		checker.cachebytes = checker.cachebytes[dataLength:]

		if len(checker.cachebytes) <= 0 {
			checker.cachebytes = nil
		} else if len(checker.cachebytes) >= 4 && !reflect.DeepEqual(checker.cachebytes[:4], kZERO_MESSAGE_HEAD) {
			global.Logger().ErrorS(fmt.Errorf("\n### err message \n%s", structs.BytesString(checker.cachebytes...)))
			checker.cachebytes = nil
		}

		if checker.cachebytes == nil {
			return comps
		} else {
			return checker.unpacking(registerId, historys...)
		}
	}
	return comps
}

func (checker *kZeroKMessageChecker) CheckPackageData(registerId string, data []byte) [][]byte {
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
	} else if reflect.DeepEqual(data[:4], kZERO_MESSAGE_HEAD) {
		checker.cachebytes = make([]byte, 0)
		checker.cachebytes = append(checker.cachebytes, data...)
	} else if checker.cachebytes != nil {
		checker.cachebytes = append(checker.cachebytes, data...)
	}

	if len(checker.cachebytes) > 0 && reflect.DeepEqual(checker.cachebytes[len(checker.cachebytes)-4:], kZERO_MESSAGE_END) {
		pkgs := checker.unpacking(registerId)
		if len(pkgs) > 0 {
			checks := make([][]byte, 0)
			for _, _pkg := range pkgs {
				v1msg := ParseKMessage(_pkg)
				err := v1msg.Check()
				if err != nil {
					global.Logger().Debug(err.Error())
				} else {
					checks = append(checks, _pkg)
				}
			}
			return checks
		}
	}
	return nil
}

type kZeroKMessageConnect struct {
	server.ZeroSocketConnect

	keeper *kZeroKMessageKeeper

	request *ZeroKMessage
	// response     *ZeroKMessage
	requestMutex sync.Mutex
	responseChan chan *ZeroKMessage
}

func (v1conn *kZeroKMessageConnect) Authorized(datas ...byte) bool {
	authMessage := ParseKMessage(datas)

	ackMessage := NewAckKMessage(MESSAGE_TYPE_CONNACK, authMessage.MessageId(), make([]byte, 0))

	err := ackMessage.Complete()
	if err != nil {
		global.Logger().ErrorS(err)
		return false
	}

	err = v1conn.pushMessage(ackMessage)
	if err != nil {
		global.Logger().ErrorS(err)
		return false
	}

	conn, err := v1conn.keeper.UseConnect(v1conn.RegisterId())
	if err == nil && conn != nil {
		conn.Close()
	}

	global.Logger().Info(fmt.Sprintf("zerov1 connect %s authorized", v1conn.RemoteAddr()))
	v1conn.Heartbeat()
	v1conn.ZeroSocketConnect.Authorized()

	return true
}

func (v1conn *kZeroKMessageConnect) Close() error {
	return v1conn.ZeroSocketConnect.Close()
}

func (v1conn *kZeroKMessageConnect) RegisterId() string {
	return v1conn.RemoteAddr()
}

func (v1conn *kZeroKMessageConnect) execMessage(message *ZeroKMessage, withSecond int) (*ZeroKMessage, error) {
	if v1conn.responseChan != nil {
		return nil, fmt.Errorf("zerov1 connect %s is busying", v1conn.RemoteAddr())
	}
	err := v1conn.Write(message.Bytes())
	if err != nil {
		return nil, err
	}
	v1conn.requestMutex.Lock()
	v1conn.request = message
	v1conn.requestMutex.Unlock()
	v1conn.responseChan = make(chan *ZeroKMessage, 1)
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

func (v1conn *kZeroKMessageConnect) pushMessage(message *ZeroKMessage) error {
	return v1conn.Write(message.Bytes())
}

func (v1conn *kZeroKMessageConnect) OnMessage(datas []byte) error {
	uMessage := ParseKMessage(datas)

	global.Logger().Debug(fmt.Sprintf("zerov1 connect %s on message \n%s", v1conn.RegisterId(), uMessage.String()))
	if uMessage.MessageType() == MESSAGE_TYPE_CONNECT {
		if !v1conn.Authorized(datas...) {
			v1conn.Close()
		}
	} else if uMessage.MessageType() == MESSAGE_TYPE_HEARTBEAT {
		v1conn.Heartbeat()
		beatack := NewAckKMessage(MESSAGE_TYPE_BEATACK, uMessage.MessageId(), make([]byte, 0))

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
			ok, err := v1conn.keeper.operator.Operation(v1conn, uMessage)
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

type kZeroKMessageKeeper struct {
	server.TCPServer

	operator ZeroKMessageOperator
}

func (keeper *kZeroKMessageKeeper) ExecMessage(registerId string, message *ZeroKMessage, withSecond int) (*ZeroKMessage, error) {
	conn, err := keeper.UseConnect(registerId)
	if err != nil {
		return nil, fmt.Errorf("use connect `%s` error: %s", registerId, err.Error())
	}
	return conn.(*kZeroKMessageConnect).execMessage(message, withSecond)
}

func (keeper *kZeroKMessageKeeper) PushMessage(registerId string, message *ZeroKMessage) error {
	conn, err := keeper.UseConnect(registerId)
	if err != nil {
		return fmt.Errorf("use connect `%s` error: %s", registerId, err.Error())
	}
	return conn.(*kZeroKMessageConnect).pushMessage(message)
}

func (keeper *kZeroKMessageKeeper) RunServer() {
	keeper.ConnectBuilder = &xZeroKMessageConnectBuilder{}
	keeper.TCPServer.RunServer()
}

func RunKMessageServer(addr string, heartbeatTime int, operator ZeroKMessageOperator) {
	zerov1serv := &kZeroKMessageKeeper{
		TCPServer: *server.NewTCPServer(
			addr,
			xDEFAULT_AUTH_WAIT,
			int64(heartbeatTime),
			xDEFAULT_BUFFER_SIZE,
		),
		operator: operator,
	}
	global.Key(ZEROKMSG_SERVER, zerov1serv)
	zerov1serv.RunServer()
}
