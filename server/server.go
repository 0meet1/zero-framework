package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/structs"
	"github.com/gofrs/uuid"
)

type ZeroServ interface {
	OnConnect(ZeroConnect) error
	OnDisconnect(ZeroConnect) error
	OnAuthorized(ZeroConnect) error
	UseConnect(string) (ZeroConnect, error)
}

type ZeroDataChecker interface {
	CheckPackageData(data []byte) []byte
}

type ZeroConnectBuilder interface {
	NewConnect() ZeroConnect
}

type ZeroConnect interface {
	structs.ZeroMetaDef

	AcceptTime() int64
	HeartbeatTime() int64

	Accept(ZeroServ, net.Conn) error
	RegisterId() string
	ConnectId() string
	RemoteAddr() string
	HeartbeatCheck(int64) bool
	Active() bool
	Heartbeat()
	Close() error
	Write([]byte) error

	Authorized(authMessage ...byte) bool
	OnMessage(datas []byte) error

	AddChecker(ZeroDataChecker)
	CheckPackageData(data []byte) []byte
}

type ZeroSocketConnect struct {
	structs.ZeroMeta

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

func (zSock *ZeroSocketConnect) This() interface{} {
	if zSock.ZeroMeta.This() == nil {
		zSock.ThisDef(zSock)
	}
	return zSock.ZeroMeta.This()
}

func (zSock *ZeroSocketConnect) AcceptTime() int64 {
	return zSock.acceptTime
}
func (zSock *ZeroSocketConnect) HeartbeatTime() int64 {
	return zSock.heartbeatTime
}

func (zSock *ZeroSocketConnect) Accept(zserv ZeroServ, connect net.Conn) error {
	zSock.connect = connect
	if zSock.zserv == nil {
		zSock.zserv = zserv
	}
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
	global.Logger().Info(fmt.Sprintf("sock connect %s on heartbeat", zSock.This().(ZeroConnect).RegisterId()))
}

func (zSock *ZeroSocketConnect) HeartbeatCheck(heartbeatSeconds int64) bool {

	if time.Now().Unix()-zSock.heartbeatTime > heartbeatSeconds {
		global.Logger().Info(fmt.Sprintf("sock connect %s exceeding heartbeat time, acceptTime %s ,heartbeatTime %s ,now %s ,heartbeat interval %ds",
			zSock.This().(ZeroConnect).RegisterId(),
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

func (zSock *ZeroSocketConnect) AddChecker(checker ZeroDataChecker) {
	zSock.checker = checker
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
	authWaitSeconds        int64
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
	sockServer.accepts[conn.This().(ZeroConnect).RegisterId()] = conn.This().(ZeroConnect)
	sockServer.acceptMutex.Unlock()
	return nil
}

func (sockServer *ZeroSocketServer) OnDisconnect(conn ZeroConnect) error {
	sockServer.acceptMutex.Lock()
	_, ok := sockServer.accepts[conn.This().(ZeroConnect).RegisterId()]
	if ok {
		delete(sockServer.accepts, conn.This().(ZeroConnect).RegisterId())
	}
	sockServer.acceptMutex.Unlock()

	sockServer.connectMutex.Lock()
	_, ok = sockServer.connects[conn.This().(ZeroConnect).RegisterId()]
	if ok {
		delete(sockServer.connects, conn.This().(ZeroConnect).RegisterId())
	}
	sockServer.connectMutex.Unlock()

	return nil
}

func (sockServer *ZeroSocketServer) OnAuthorized(conn ZeroConnect) error {
	sockServer.acceptMutex.Lock()
	_acceptConnect, ok := sockServer.accepts[conn.This().(ZeroConnect).RegisterId()]
	if ok {
		delete(sockServer.accepts, _acceptConnect.RegisterId())
	}
	sockServer.acceptMutex.Unlock()

	sockServer.connectMutex.Lock()
	sockServer.connects[conn.This().(ZeroConnect).RegisterId()] = conn.This().(ZeroConnect)
	sockServer.connectMutex.Unlock()

	return nil
}

func (sockServer *ZeroSocketServer) UseConnect(registerId string) (ZeroConnect, error) {
	sockServer.connectMutex.RLock()
	connect, ok := sockServer.connects[registerId]
	sockServer.connectMutex.RUnlock()
	if ok {
		return connect, nil
	}
	return nil, fmt.Errorf("connect %s not found", registerId)
}

func (sockServer *ZeroSocketServer) initHeartbeatTimer() {
	sockServer.heartbeatTimer = time.NewTimer(time.Second * time.Duration(sockServer.heartbeatCheckInterval))
	for {
		<-sockServer.heartbeatTimer.C
		global.Logger().Info("sock heartbeat check starting")
		removes := make([]ZeroConnect, 0)
		sockServer.connectMutex.RLock()
		for _, connect := range sockServer.connects {
			if !connect.HeartbeatCheck(sockServer.heartbeatSeconds) {
				removes = append(removes, connect)
			}
		}
		sockServer.connectMutex.RUnlock()

		for _, conn := range removes {
			global.Logger().Info(fmt.Sprintf("sock connect %s heartbeat timeout", conn.This().(ZeroConnect).RegisterId()))
			err := conn.This().(ZeroConnect).Close()
			if err != nil {
				global.Logger().Error(fmt.Sprintf("sock connect check %s closing error : %s", conn.This().(ZeroConnect).RegisterId(), err.Error()))
			}
		}
		global.Logger().Info("sock heartbeat check finished")
		sockServer.heartbeatTimer = time.NewTimer(time.Second * time.Duration(sockServer.heartbeatCheckInterval))
	}
}

func (sockServer *ZeroSocketServer) accept(conn net.Conn) {
	connect := sockServer.ConnectBuilder.NewConnect()
	connect.Accept(sockServer, conn)

	global.Logger().Info(fmt.Sprintf("sock server accept connect -> %s", connect.This().(ZeroConnect).RegisterId()))

	time.AfterFunc(time.Duration(sockServer.authWaitSeconds)*time.Second, func() {
		sockServer.connectMutex.RLock()
		_, ok := sockServer.connects[connect.This().(ZeroConnect).RegisterId()]
		sockServer.connectMutex.RUnlock()
		if !ok {
			connect.Close()
			global.Logger().Info(fmt.Sprintf("sock server connect auth time out -> %s", connect.This().(ZeroConnect).RegisterId()))
		} else {
			global.Logger().Info(fmt.Sprintf("sock server connect auth checked -> %s", connect.This().(ZeroConnect).RegisterId()))
		}
	})

	defer func() {
		connect.Close()
		global.Logger().Info(fmt.Sprintf("sock server connect close -> %s", connect.This().(ZeroConnect).RegisterId()))
	}()
	dataBuf := make([]byte, sockServer.bufferSize)
	for {
		if !connect.Active() {
			global.Logger().Error(fmt.Sprintf("sock server connect %s is already closed", connect.This().(ZeroConnect).RegisterId()))
			break
		}

		dataLen, err := conn.Read(dataBuf[:])
		if err != nil {
			global.Logger().Error(fmt.Sprintf("sock server connect %s on message error %s", connect.This().(ZeroConnect).RegisterId(), err.Error()))
			break
		}

		data := dataBuf[:dataLen]
		messageDatas := connect.CheckPackageData(data)
		if messageDatas != nil {
			err = connect.OnMessage(messageDatas)
			if err != nil {
				global.Logger().Error(fmt.Sprintf("sock server connect %s on message error %s", connect.This().(ZeroConnect).RegisterId(), err.Error()))
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

type ZeroClientListener interface {
	OnConnect(ZeroClientConnect) error
	OnHeartbeat(ZeroClientConnect) error
}

type ZeroClientConnect interface {
	structs.ZeroMetaDef

	Connect()
	RemoteAddr() string
	HeartbeatCheck(int64) bool
	Active() bool
	Heartbeat()
	Close() error
	Write([]byte) error

	OnMessage([]byte) error

	AddChecker(ZeroDataChecker)
	CheckPackageData([]byte) []byte

	AddListener(ZeroClientListener)
}
