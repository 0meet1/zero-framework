package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/gofrs/uuid"
)

type ZeroServ interface {
	OnConnect(ZeroConnect) error
	OnDisconnect(ZeroConnect) error
	OnAuthorized(ZeroConnect) error
}

type ZeroDataChecker interface {
	CheckPackageData(data []byte) []byte
}

type ZeroConnectBuilder interface {
	NewConnect() ZeroConnect
}

type ZeroConnect interface {
	Accept(ZeroServ, net.Conn) error
	RegisterId() string
	ConnectId() string
	RemoteAddr() string
	HeartbeatCheck(heartbeatSeconds int64) bool
	Active() bool
	Heartbeat()
	Close() error
	Write([]byte) error

	Authorized(authMessage ...byte) bool
	OnMessage(datas []byte) error
	CheckPackageData(data []byte) []byte
}

type ZeroSocketConnect struct {
	connectId  string
	acceptTime int64

	connect      net.Conn
	connectMutex sync.Mutex

	heartbeatTime  int64
	heartbeatMutex sync.Mutex

	active  bool
	zserv   ZeroServ
	checker ZeroDataChecker
}

func (zSock *ZeroSocketConnect) Accept(zserv ZeroServ, connect net.Conn) error {
	zSock.connect = connect
	zSock.zserv = zserv
	zSock.acceptTime = time.Now().Unix()
	zSock.heartbeatTime = 0

	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	zSock.connectId = uid.String()
	err = zSock.zserv.OnConnect(zSock)
	if err != nil {
		return err
	}

	zSock.active = true
	return nil
}

func (zSock *ZeroSocketConnect) Close() error {
	if !zSock.active {
		return nil
	}

	err := zSock.zserv.OnDisconnect(zSock)
	if err != nil {
		return err
	}
	zSock.active = false

	return zSock.connect.Close()
}

func (zSock *ZeroSocketConnect) RegisterId() string {
	return zSock.ConnectId()
}

func (zSock *ZeroSocketConnect) ConnectId() string {
	return zSock.connectId
}

func (zSock *ZeroSocketConnect) RemoteAddr() string {
	return zSock.connect.RemoteAddr().String()
}

func (zSock *ZeroSocketConnect) Active() bool {
	return zSock.active
}

func (zSock *ZeroSocketConnect) Authorized(authMessage ...byte) bool {
	zSock.zserv.OnAuthorized(zSock)
	return true
}

func (zSock *ZeroSocketConnect) Heartbeat() {
	zSock.heartbeatMutex.Lock()
	zSock.heartbeatTime = time.Now().Unix()
	zSock.heartbeatMutex.Unlock()
	global.Logger().Info(fmt.Sprintf("sock connect %s on heartbeat", zSock.ConnectId()))
}

func (zSock *ZeroSocketConnect) HeartbeatCheck(heartbeatSeconds int64) bool {

	if time.Now().Unix()-zSock.heartbeatTime > heartbeatSeconds {
		global.Logger().Info(fmt.Sprintf("ipc connect %s exceeding heartbeat time, acceptTime %s ,heartbeatTime %s ,now %s ,heartbeat interval %ds",
			zSock.ConnectId(),
			time.Unix(zSock.acceptTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(zSock.heartbeatTime, 0).Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
			heartbeatSeconds))
		return false
	}
	return true
}

func (zSock *ZeroSocketConnect) OnMessage(datas []byte) error {
	zSock.Heartbeat()
	return nil
}

func (zSock *ZeroSocketConnect) Write(datas []byte) error {
	zSock.connectMutex.Lock()
	_, err := zSock.connect.Write(datas)
	zSock.connectMutex.Unlock()
	return err
}

func (zSock *ZeroSocketConnect) CheckPackageData(data []byte) []byte {
	if zSock.checker != nil {
		return zSock.checker.CheckPackageData(data)
	}
	return data
}

type xDefaultConnectBuilder struct{}

func (xDefault *xDefaultConnectBuilder) NewConnect() ZeroConnect {
	return &ZeroSocketConnect{}
}

type ZeroSocketServer struct {
	heartbeatSeconds       int64
	heartbeatCheckInterval int64
	bufferSize             int

	accepts  map[string]ZeroConnect
	connects map[string]ZeroConnect

	acceptMutex  sync.RWMutex
	connectMutex sync.RWMutex

	heartbeatTimer *time.Timer

	ConnectBuilder ZeroConnectBuilder
}

func (sockServer *ZeroSocketServer) OnConnect(conn ZeroConnect) error {
	sockServer.acceptMutex.Lock()
	sockServer.accepts[conn.RegisterId()] = conn
	sockServer.acceptMutex.Unlock()
	return nil
}

func (sockServer *ZeroSocketServer) OnDisconnect(conn ZeroConnect) error {
	sockServer.acceptMutex.Lock()
	_, ok := sockServer.accepts[conn.RegisterId()]
	if ok {
		delete(sockServer.accepts, conn.RegisterId())
	}
	sockServer.acceptMutex.Unlock()

	sockServer.connectMutex.Lock()
	_, ok = sockServer.connects[conn.RegisterId()]
	if ok {
		delete(sockServer.connects, conn.RegisterId())
	}
	sockServer.connectMutex.Unlock()

	return nil
}

func (sockServer *ZeroSocketServer) OnAuthorized(conn ZeroConnect) error {
	sockServer.acceptMutex.Lock()
	_acceptConnect, ok := sockServer.accepts[conn.RegisterId()]
	if ok {
		delete(sockServer.accepts, _acceptConnect.RegisterId())
	}
	sockServer.acceptMutex.Unlock()

	sockServer.connectMutex.Lock()
	sockServer.connects[conn.RegisterId()] = conn
	sockServer.connectMutex.Unlock()

	return nil
}

func (sockServer *ZeroSocketServer) initHeartbeatTimer() {
	sockServer.heartbeatTimer = time.NewTimer(time.Second * time.Duration(sockServer.heartbeatCheckInterval))
	for {
		select {
		case <-sockServer.heartbeatTimer.C:
			global.Logger().Info(fmt.Sprintf("sock heartbeat check starting"))
			removes := make([]ZeroConnect, 0)
			sockServer.connectMutex.RLock()
			for _, connect := range sockServer.connects {
				if !connect.HeartbeatCheck(sockServer.heartbeatSeconds) {
					removes = append(removes, connect)
				}
			}
			sockServer.connectMutex.RUnlock()

			for _, conn := range removes {
				global.Logger().Info(fmt.Sprintf("sock connect %s heartbeat timeout", conn.ConnectId()))
				err := conn.Close()
				if err != nil {
					global.Logger().Error(fmt.Sprintf("sock connect check %s closing error : %s", conn.ConnectId(), err.Error()))
				}
			}
			global.Logger().Info(fmt.Sprintf("sock heartbeat check finished"))
		}
	}
}

func (sockServer *ZeroSocketServer) accept(conn net.Conn) {
	connect := sockServer.ConnectBuilder.NewConnect()
	connect.Accept(sockServer, conn)

	global.Logger().Info(fmt.Sprintf("sock server accept connect -> %s", connect.ConnectId()))

	time.AfterFunc(5*time.Second, func() {
		sockServer.connectMutex.RLock()
		_, ok := sockServer.connects[connect.ConnectId()]
		sockServer.connectMutex.RUnlock()
		if !ok {
			connect.Close()
			global.Logger().Info(fmt.Sprintf("scok server connect auth time out -> %s", connect.ConnectId()))
		} else {
			global.Logger().Info(fmt.Sprintf("sock server connect auth checked -> %s", connect.ConnectId()))
		}
	})

	defer func() {
		connect.Close()
		global.Logger().Info(fmt.Sprintf("ipc server connect close -> %s", connect.ConnectId()))
	}()
	dataBuf := make([]byte, sockServer.bufferSize)
	for {
		if !connect.Active() {
			global.Logger().Error(fmt.Sprintf("ipc server connect %s is already closed", connect.ConnectId()))
			break
		}

		dataLen, err := conn.Read(dataBuf[:])
		if err != nil {
			global.Logger().Error(fmt.Sprintf("ipc server connect %s on message error %s", connect.ConnectId(), err.Error()))
			break
		}

		data := dataBuf[:dataLen]
		messageDatas := connect.CheckPackageData(data)
		if messageDatas != nil {
			err = connect.OnMessage(messageDatas)
			if err != nil {
				global.Logger().Error(fmt.Sprintf("sock server connect %s on message error %s", connect.ConnectId(), err.Error()))
			}
		}
	}
}

func (sockServer *ZeroSocketServer) RunServer() {
	if sockServer.ConnectBuilder == nil {
		sockServer.ConnectBuilder = &xDefaultConnectBuilder{}
	}
	go sockServer.initHeartbeatTimer()
}
