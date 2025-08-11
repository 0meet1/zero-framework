package server

import (
	"container/list"
	"container/ring"
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
	OnHeartbeat(ZeroConnect) error
	UseConnect(string) (ZeroConnect, error)
}

type ZeroDataChecker interface {
	CheckPackageData(string, []byte) [][]byte
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

	Node() *list.Element
	Clock() *ring.Ring
	FlushNode(*list.Element)
	FlushClock(*ring.Ring)

	Authorized(authMessage ...byte) bool
	OnMessage(datas []byte) error

	AddChecker(ZeroDataChecker)
	CheckPackageData(data []byte) [][]byte
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

	node  *list.Element
	clock *ring.Ring
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
func (zSock *ZeroSocketConnect) Node() *list.Element {
	return zSock.node
}
func (zSock *ZeroSocketConnect) Clock() *ring.Ring {
	return zSock.clock
}
func (zSock *ZeroSocketConnect) FlushNode(node *list.Element) {
	zSock.node = node
}
func (zSock *ZeroSocketConnect) FlushClock(clock *ring.Ring) {
	zSock.clock = clock
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
	zSock.zserv.OnHeartbeat(zSock)
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

func (zSock *ZeroSocketConnect) CheckPackageData(data []byte) [][]byte {
	if zSock.checker != nil {
		return zSock.checker.CheckPackageData(zSock.This().(ZeroConnect).RegisterId(), data)
	}
	return [][]byte{data}
}

type xDefaultConnectBuilder struct{}

func (xDefault *xDefaultConnectBuilder) NewConnect() ZeroConnect {
	return &ZeroSocketConnect{}
}

type ZeroServerWatcher interface {
	WatcherName() string
	OnConnect(ZeroConnect) error
	OnAuthorized(ZeroConnect) error
	OnDisconnect(ZeroConnect) error
	OnHeartbeat(ZeroConnect) error
	OnMessage(ZeroConnect, []byte) error
}

type ZeroSocketServer struct {
	authWaitSeconds  int64
	heartbeatSeconds int64
	bufferSize       int

	connects     map[string]ZeroConnect
	connectMutex sync.RWMutex

	ConnectBuilder ZeroConnectBuilder

	acceptClock    *ring.Ring
	heartbeatClock *ring.Ring

	watchers []ZeroServerWatcher
}

func (sockServer *ZeroSocketServer) AddWatchers(watchers ...ZeroServerWatcher) {
	if sockServer.watchers == nil {
		sockServer.watchers = make([]ZeroServerWatcher, 0)
	}
	sockServer.watchers = append(sockServer.watchers, watchers...)
}

func (sockServer *ZeroSocketServer) notifyOnConnect(conn ZeroConnect) {
	if len(sockServer.watchers) <= 0 {
		return
	}
	call := func(watcher ZeroServerWatcher) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger().Errorf("watcher `%s` on error: %s", watcher.WatcherName(), err.(error))
			}
		}()
		err := watcher.OnConnect(conn)
		if err != nil {
			panic(err)
		}
	}
	for _, watcher := range sockServer.watchers {
		call(watcher)
	}
}

func (sockServer *ZeroSocketServer) OnConnect(conn ZeroConnect) error {
	defer sockServer.notifyOnConnect(conn)
	clock := sockServer.acceptClock.Prev()
	node := clock.Value.(*structs.ZeroLinked).PushBack(conn)
	conn.FlushClock(clock)
	conn.FlushNode(node)
	return nil
}

func (sockServer *ZeroSocketServer) notifyOnDisconnect(conn ZeroConnect) {
	if len(sockServer.watchers) <= 0 {
		return
	}
	call := func(watcher ZeroServerWatcher) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger().Errorf("watcher `%s` on error: %s", watcher.WatcherName(), err.(error))
			}
		}()
		err := watcher.OnDisconnect(conn)
		if err != nil {
			panic(err)
		}
	}
	for _, watcher := range sockServer.watchers {
		call(watcher)
	}
}

func (sockServer *ZeroSocketServer) OnDisconnect(conn ZeroConnect) error {
	defer sockServer.notifyOnDisconnect(conn)
	conn.Clock().Value.(*structs.ZeroLinked).Remove(conn.Node())
	sockServer.connectMutex.Lock()
	_, ok := sockServer.connects[conn.This().(ZeroConnect).RegisterId()]
	if ok {
		delete(sockServer.connects, conn.This().(ZeroConnect).RegisterId())
	}
	sockServer.connectMutex.Unlock()

	return nil
}

func (sockServer *ZeroSocketServer) notifyOnAuthorized(conn ZeroConnect) {
	if len(sockServer.watchers) <= 0 {
		return
	}
	call := func(watcher ZeroServerWatcher) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger().Errorf("watcher `%s` on error: %s", watcher.WatcherName(), err.(error))
			}
		}()
		err := watcher.OnAuthorized(conn)
		if err != nil {
			panic(err)
		}
	}
	for _, watcher := range sockServer.watchers {
		call(watcher)
	}
}

func (sockServer *ZeroSocketServer) OnAuthorized(conn ZeroConnect) error {
	defer sockServer.notifyOnAuthorized(conn)
	conn.Clock().Value.(*structs.ZeroLinked).Remove(conn.Node())
	clock := sockServer.heartbeatClock.Prev()
	node := clock.Value.(*structs.ZeroLinked).PushBack(conn)
	conn.FlushClock(clock)
	conn.FlushNode(node)

	sockServer.connectMutex.Lock()
	sockServer.connects[conn.This().(ZeroConnect).RegisterId()] = conn.This().(ZeroConnect)
	sockServer.connectMutex.Unlock()

	return nil
}

func (sockServer *ZeroSocketServer) notifyOnHeartbeat(conn ZeroConnect) {
	if len(sockServer.watchers) <= 0 {
		return
	}
	call := func(watcher ZeroServerWatcher) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger().Errorf("watcher `%s` on error: %s", watcher.WatcherName(), err.(error))
			}
		}()
		err := watcher.OnHeartbeat(conn)
		if err != nil {
			panic(err)
		}
	}
	for _, watcher := range sockServer.watchers {
		call(watcher)
	}
}

func (sockServer *ZeroSocketServer) OnHeartbeat(conn ZeroConnect) error {
	defer sockServer.notifyOnHeartbeat(conn)
	conn.Clock().Value.(*structs.ZeroLinked).Remove(conn.Node())
	clock := sockServer.heartbeatClock.Prev()
	node := clock.Value.(*structs.ZeroLinked).PushBack(conn)
	conn.FlushClock(clock)
	conn.FlushNode(node)
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

func (sockServer *ZeroSocketServer) acceptTimeoutConnect(nodes *structs.ZeroLinked) {
	_nodes := nodes.Front()
	for _nodes != nil {
		conn := _nodes.Value.(ZeroConnect)
		global.Logger().Info(fmt.Sprintf("sock server connect auth time out -> %s", conn.This().(ZeroConnect).RegisterId()))
		err := conn.This().(ZeroConnect).Close()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("sock connect auth %s closing error : %s", conn.This().(ZeroConnect).RegisterId(), err.Error()))
		}
		_nodes = _nodes.Next()
	}
}

func (sockServer *ZeroSocketServer) runAcceptLoop() {
	const DEFAULT_ACCEPT_CHECK = 10
	sockServer.acceptClock = ring.New(DEFAULT_ACCEPT_CHECK)
	for i := 0; i < DEFAULT_ACCEPT_CHECK; i++ {
		sockServer.acceptClock.Value = structs.NewLinked()
		sockServer.acceptClock = sockServer.acceptClock.Next()
	}
	for {
		<-time.After(time.Duration(1) * time.Second)
		v := sockServer.acceptClock.Value.(*structs.ZeroLinked)
		go sockServer.acceptTimeoutConnect(v)
		sockServer.acceptClock.Value = structs.NewLinked()
		sockServer.acceptClock = sockServer.acceptClock.Next()
	}
}

func (sockServer *ZeroSocketServer) cleanTimeoutConnect(nodes *structs.ZeroLinked) {
	_nodes := nodes.Front()
	for _nodes != nil {
		conn := _nodes.Value.(ZeroConnect)
		global.Logger().Info(fmt.Sprintf("sock connect %s heartbeat timeout", conn.This().(ZeroConnect).RegisterId()))
		err := conn.This().(ZeroConnect).Close()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("sock connect check %s closing error : %s", conn.This().(ZeroConnect).RegisterId(), err.Error()))
		}
		_nodes = _nodes.Next()
	}
}

func (sockServer *ZeroSocketServer) runHeartbeatLoop() {
	sockServer.heartbeatClock = ring.New(int(sockServer.heartbeatSeconds))
	for i := 0; i < int(sockServer.heartbeatSeconds); i++ {
		sockServer.heartbeatClock.Value = structs.NewLinked()
		sockServer.heartbeatClock = sockServer.heartbeatClock.Next()
	}
	for {
		<-time.After(time.Duration(1) * time.Second)
		v := sockServer.heartbeatClock.Value.(*structs.ZeroLinked)
		go sockServer.cleanTimeoutConnect(v)
		sockServer.heartbeatClock.Value = structs.NewLinked()
		sockServer.heartbeatClock = sockServer.heartbeatClock.Next()
	}
}

func (sockServer *ZeroSocketServer) notifyOnMessage(conn ZeroConnect, datas []byte) {
	if len(sockServer.watchers) <= 0 {
		return
	}
	call := func(watcher ZeroServerWatcher) {
		defer func() {
			err := recover()
			if err != nil {
				global.Logger().Errorf("watcher `%s` on error: %s", watcher.WatcherName(), err.(error))
			}
		}()
		err := watcher.OnMessage(conn, datas)
		if err != nil {
			panic(err)
		}
	}
	for _, watcher := range sockServer.watchers {
		call(watcher)
	}
}

func (sockServer *ZeroSocketServer) accept(conn net.Conn) {
	connect := sockServer.ConnectBuilder.NewConnect()
	connect.Accept(sockServer, conn)

	global.Logger().Info(fmt.Sprintf("sock server accept connect -> %s", connect.This().(ZeroConnect).RegisterId()))

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
		if len(messageDatas) > 0 {
			for _, messageData := range messageDatas {
				err = connect.OnMessage(messageData)
				if err != nil {
					global.Logger().Error(fmt.Sprintf("sock server connect %s on message error %s", connect.This().(ZeroConnect).RegisterId(), err.Error()))
				} else {
					sockServer.notifyOnMessage(connect, messageData)
				}
			}
		}
	}
}

func (sockServer *ZeroSocketServer) RunServer() {
	if sockServer.ConnectBuilder == nil {
		sockServer.ConnectBuilder = &xDefaultConnectBuilder{}
	}
	go sockServer.runAcceptLoop()
	go sockServer.runHeartbeatLoop()
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
	CheckPackageData([]byte) [][]byte

	AddListener(ZeroClientListener)
}
