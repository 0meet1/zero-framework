package server

import (
	"fmt"
	"net"

	"github.com/0meet1/zero-framework/global"
)

type TCPServer struct {
	ZeroSocketServer

	address   string
	tcpServer net.Listener
}

func NewTCPServer(address string, heartbeatSeconds int64, heartbeatCheckInterval int64, bufferSize int) *TCPServer {
	return &TCPServer{
		ZeroSocketServer: ZeroSocketServer{
			accepts:                make(map[string]ZeroConnect),
			connects:               make(map[string]ZeroConnect),
			heartbeatSeconds:       heartbeatSeconds,
			heartbeatCheckInterval: heartbeatCheckInterval,
			bufferSize:             bufferSize,
		},
		address: address,
	}
}

func (tcpserv *TCPServer) RunServer() {
	tcpserv.ZeroSocketServer.RunServer()

	tcpServer, err := net.Listen("tcp", tcpserv.address)
	if err != nil {
		global.Logger().Error(fmt.Sprintf("tcp server start error : %s", err.Error()))
		panic(err)
	}
	tcpserv.tcpServer = tcpServer

	global.Logger().Info(fmt.Sprintf("tcp server start success on tcp://%s", tcpserv.address))

	for {
		conn, err := tcpServer.Accept()
		if err != nil {
			global.Logger().Error(fmt.Sprintf("tcp server accept error : %s", err.Error()))
			continue
		}
		go tcpserv.accept(conn)
	}
}
