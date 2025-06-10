package server

import (
	"fmt"
	"net"

	"github.com/0meet1/zero-framework/global"
)

type UDPMessageProcesser interface {
	OnMessage([]byte) error
}

type UDPServer struct {
	bufferSize int

	port    int
	udpconn *net.UDPConn

	checker   ZeroDataChecker
	processer UDPMessageProcesser
}

func NewUDPServer(port int, bufferSize int, checker ZeroDataChecker, processer UDPMessageProcesser) *UDPServer {
	return &UDPServer{
		bufferSize: bufferSize,
		port:       port,
		checker:    checker,
		processer:  processer,
	}
}

func (udpserv *UDPServer) Write(datas []byte, addr *net.UDPAddr) error {
	_, err := udpserv.udpconn.WriteToUDP(datas, addr)
	return err
}

func (udpserv *UDPServer) checkPackageData(data []byte) []byte {
	if udpserv.checker != nil {
		return udpserv.checker.CheckPackageData(fmt.Sprintf(":%d", udpserv.port), data)
	}
	return data
}

func (udpserv *UDPServer) read() {
	defer udpserv.udpconn.Close()
	dataBuf := make([]byte, udpserv.bufferSize)
	for {
		dataLen, addr, err := udpserv.udpconn.ReadFromUDP(dataBuf[:])
		if err != nil {
			global.Logger().Error(fmt.Sprintf("udp:%d read failed, err: %s", udpserv.port, err.Error()))
			continue
		}

		global.Logger().Debug(fmt.Sprintf("udp:%d from: %s:%d on message, data length: %d", udpserv.port, addr.IP, addr.Port, dataLen))

		if udpserv.processer != nil {
			data := dataBuf[:dataLen]
			messageDatas := udpserv.checkPackageData(data)
			if messageDatas != nil {
				err = udpserv.processer.OnMessage(messageDatas)
				if err != nil {
					global.Logger().Error(fmt.Sprintf("udp:%d on message error %s", udpserv.port, err.Error()))
				}
			}
		}
	}
}

func (udpserv *UDPServer) RunServer() {
	udpconn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: udpserv.port,
	})
	udpserv.udpconn = udpconn
	if err != nil {
		panic(fmt.Sprintf("udp Listen port: %d failed, reason :%s", udpserv.port, err.Error()))
	}
	go udpserv.read()
}
