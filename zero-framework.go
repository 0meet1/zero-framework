package zeroframework

import (
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/rocketmq"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

type Date = structs.Date
type ZeroCoreStructs = structs.ZeroCoreStructs
type ZeroRequest = structs.ZeroRequest
type ZeroResponse = structs.ZeroResponse

type ZeroCoreProcessor = processors.ZeroCoreProcessor
type ZeroQuery = processors.ZeroQuery
type ZeroCondition = processors.ZeroCondition
type ZeroOrderBy = processors.ZeroOrderBy
type ZeroLimit = processors.ZeroLimit
type ZeroQueryOperation = processors.ZeroQueryOperation

type EQueryRequest = database.EQueryRequest
type EQueryResponse = database.EQueryResponse
type EQuerySearch = database.EQuerySearch
type DataSource = database.DataSource

type MQNotifyMessage = rocketmq.MQNotifyMessage
type MQMessageObserver = rocketmq.MQMessageObserver

type ZeroServ = server.ZeroServ
type ZeroDataChecker = server.ZeroDataChecker
type ZeroConnectBuilder = server.ZeroConnectBuilder
type ZeroConnect = server.ZeroConnect

type ZeroSocketConnect = server.ZeroSocketConnect
type ZeroSocketServer = server.ZeroSocketServer

type UDPMessageProcesser = server.UDPMessageProcesser

type IPCServer = server.IPCServer
type TCPServer = server.TCPServer
type UDPServer = server.UDPServer

type MqttFixedHeader = server.MqttFixedHeader
type MqttCoreVariableHeader = server.MqttCoreVariableHeader
type MqttVariableHeader = server.MqttVariableHeader
type MqttCorePayload = server.MqttCorePayload
type MqttPayload = server.MqttPayload
type MqttMessage = server.MqttMessage

type MqttIdentifierVariableHeader = server.MqttIdentifierVariableHeader
type MqttConnectVariableHeader = server.MqttConnectVariableHeader
type MqttConnackVariableHeader = server.MqttConnackVariableHeader
type MqttPublishVariableHeader = server.MqttPublishVariableHeader

type MqttParamsPayload = server.MqttParamsPayload
type MqttTopic = server.MqttTopic

type MqttMessageProcessor = server.MqttMessageProcessor
type MqttConnect = server.MqttConnect
type MqttServer = server.MqttServer
