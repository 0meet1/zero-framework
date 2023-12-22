package server

import (
	"fmt"
	"net"
	"os"
	"path"

	"github.com/0meet1/zero-framework/global"
)

type IPCServer struct {
	ZeroSocketServer

	ipcsock   string
	ipcServer *net.UnixListener
}

func NewIPCServer(ipcsock string, heartbeatSeconds int64, heartbeatCheckInterval int64, bufferSize int) *IPCServer {
	return &IPCServer{
		ZeroSocketServer: ZeroSocketServer{
			heartbeatSeconds:       heartbeatSeconds,
			heartbeatCheckInterval: heartbeatCheckInterval,
			bufferSize:             bufferSize,
		},
		ipcsock: ipcsock,
	}
}

func (ipcserv *IPCServer) RunServer() {
	ipcserv.ZeroSocketServer.RunServer()

	_, err := os.Stat(path.Dir(ipcserv.ipcsock))
	if err != nil {
		os.MkdirAll(ipcserv.ipcsock, os.ModePerm)
	}

	os.RemoveAll(ipcserv.ipcsock)
	ipcaddr, err := net.ResolveUnixAddr("unix", ipcserv.ipcsock)
	if err != nil {
		panic(err)
	}
	ipcserv.ipcServer, err = net.ListenUnix("unix", ipcaddr)
	if err != nil {
		panic(err)
	}

	global.Logger().Info(fmt.Sprintf("ipc server start success on ipc://%s", ipcserv.ipcsock))

	go func() {
		for {
			conn, err := ipcserv.ipcServer.Accept()
			if err != nil {
				global.Logger().Error(fmt.Sprintf("ipc server accept error : %s", err.Error()))
				continue
			}
			go ipcserv.accept(conn)
		}
	}()
}
