package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/structs"
)

type TCPClient struct {
	structs.ZeroMeta

	connAddr     string
	connect      net.Conn
	connectMutex sync.Mutex

	authWaitSeconds        int64
	heartbeatSeconds       int64
	heartbeatCheckInterval int64
	bufferSize             int

	heartbeatTimer *time.Timer

	connectTime    int64
	heartbeatTime  int64
	heartbeatMutex sync.Mutex

	checker   ZeroDataChecker
	xListener ZeroClientListener
}

func (client *TCPClient) initHeartbeatTimer() {
	client.heartbeatTimer = time.NewTimer(time.Second * time.Duration(client.heartbeatCheckInterval))
	for {
		select {
		case <-client.heartbeatTimer.C:
			if !client.HeartbeatCheck(client.heartbeatSeconds) {
				client.connect.Close()
			} else if client.xListener != nil {
				err := client.xListener.OnHeartbeat(client.This().(ZeroClientConnect))
				if err != nil {
					global.Logger().Error(err.Error())
				}
			}
		}
	}
}

func (client *TCPClient) This() interface{} {
	if client.ZeroMeta.This() == nil {
		client.ThisDef(client)
	}
	return client.ZeroMeta.This()
}

func (client *TCPClient) Connect() {
	client.startingLoop()
}

func (client *TCPClient) RemoteAddr() string {
	if client.connect != nil {
		return client.connect.RemoteAddr().String()
	} else {
		return "disconnect"
	}
}

func (client *TCPClient) HeartbeatCheck(heartbeatSeconds int64) bool {
	if time.Now().Unix()-client.heartbeatTime > heartbeatSeconds {
		global.Logger().Info(fmt.Sprintf("tcp client connect %s exceeding heartbeat time, acceptTime %s ,heartbeatTime %s ,now %s ,heartbeat interval %ds",
			client.connect.RemoteAddr().String(),
			time.Unix(client.connectTime, 0).Format("2006-01-02 15:04:05"),
			time.Unix(client.heartbeatTime, 0).Format("2006-01-02 15:04:05"),
			time.Now().Format("2006-01-02 15:04:05"),
			heartbeatSeconds))
		return false
	}
	return true
}

func (client *TCPClient) Active() bool {
	return client.connect != nil
}

func (client *TCPClient) Heartbeat() {
	client.heartbeatMutex.Lock()
	client.heartbeatTime = time.Now().Unix()
	client.heartbeatMutex.Unlock()
	global.Logger().Info(fmt.Sprintf("tcp client connect %s on heartbeat", client.connect.RemoteAddr()))
}

func (client *TCPClient) CheckPackageData(data []byte) []byte {
	if client.checker != nil {
		return client.checker.CheckPackageData(data)
	}
	return data
}

func (client *TCPClient) AddListener(xListener ZeroClientListener) {
	client.xListener = xListener
}

func (client *TCPClient) Close() error {
	err := client.connect.Close()
	client.connect = nil
	return err
}

func (client *TCPClient) Write(datas []byte) error {
	client.connectMutex.Lock()
	_, err := client.connect.Write(datas)
	client.connectMutex.Unlock()
	return err
}

func (client *TCPClient) receive() {
	time.AfterFunc(time.Duration(client.authWaitSeconds)*time.Second, func() {
		if !client.This().(ZeroClientConnect).Active() {
			client.connect.Close()
			global.Logger().Info(fmt.Sprintf("tcp client connect auth time out -> %s", client.RemoteAddr()))
		} else {
			global.Logger().Info(fmt.Sprintf("tcp client connect auth checked -> %s", client.RemoteAddr()))
		}
	})

	defer func() {
		client.connect.Close()
		global.Logger().Info(fmt.Sprintf("tcp client connect close -> %s", client.RemoteAddr()))

		client.connect = nil
		client.startingLoop()
	}()

	dataBuf := make([]byte, client.bufferSize)
	for {
		dataLen, err := client.connect.Read(dataBuf[:])
		if err != nil {
			global.Logger().Error(fmt.Sprintf("tcp client connect %s on message error %s", client.connect.RemoteAddr().String(), err.Error()))
			break
		}

		data := dataBuf[:dataLen]
		messageDatas := client.CheckPackageData(data)
		if messageDatas != nil {
			err = client.This().(ZeroClientConnect).OnMessage(messageDatas)
			if err != nil {
				global.Logger().Error(fmt.Sprintf("tcp client connect %s on message error %s", client.connect.RemoteAddr().String(), err.Error()))
			}
		}
	}
}

func (client *TCPClient) startingLoop() {
	for {
		<-time.After(time.Duration(time.Second * 5))
		global.Logger().Info(fmt.Sprintf("tcp client starting -> %s", client.connAddr))
		err := client.start()
		if err != nil {
			global.Logger().Error(err.Error())
			global.Logger().Info(fmt.Sprintf("tcp client will restart after 5s -> %s", client.connAddr))
		} else {
			global.Logger().Info(fmt.Sprintf("tcp client start success -> %s", client.connAddr))
			break
		}
	}
}

func (client *TCPClient) start() error {
	if client.heartbeatTimer != nil {
		client.heartbeatTimer.Stop()
		client.heartbeatTimer = nil
	}

	conn, err := net.DialTimeout("tcp", client.connAddr, time.Second*time.Duration(30))
	if err != nil {
		return err
	}
	client.connect = conn
	go client.receive()

	if client.xListener != nil {
		err := client.xListener.OnConnect(client.This().(ZeroClientConnect))
		if err != nil {
			return err
		}
	}
	client.connectTime = time.Now().Unix()
	return nil
}

func (client *TCPClient) OnMessage(datas []byte) error {
	client.This().(ZeroClientConnect).Heartbeat()
	return nil
}

func (client *TCPClient) AddChecker(checker ZeroDataChecker) {
	client.checker = checker
}

func NewTCPClient(address string, authWaitSeconds int64, heartbeatSeconds int64, heartbeatCheckInterval int64, bufferSize int) *TCPClient {
	return &TCPClient{
		connAddr:               address,
		authWaitSeconds:        authWaitSeconds,
		heartbeatSeconds:       heartbeatSeconds,
		heartbeatCheckInterval: heartbeatCheckInterval,
		bufferSize:             bufferSize,
	}
}
