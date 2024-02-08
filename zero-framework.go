package zeroframework

import (
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/mfgrc"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/protocol"
	"github.com/0meet1/zero-framework/rocketmq"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
)

type Date = structs.Date
type ZeroMetaDef = structs.ZeroMetaDef
type ZeroMetaPtr = structs.ZeroMetaPtr
type ZeroMeta = structs.ZeroMeta
type ZeroCoreStructs = structs.ZeroCoreStructs
type ZeroRequest = structs.ZeroRequest
type ZeroResponse = structs.ZeroResponse

type ZeroCoreProcessor = processors.ZeroCoreProcessor
type ZeroQuery = processors.ZeroQuery
type ZeroCondition = processors.ZeroCondition
type ZeroOrderBy = processors.ZeroOrderBy
type ZeroLimit = processors.ZeroLimit
type ZeroQueryOperation = processors.ZeroQueryOperation

const DATABASE_MYSQL = database.DATABASE_MYSQL
const DATABASE_POSTGRES = database.DATABASE_POSTGRES
const DATABASE_REDIS = database.DATABASE_REDIS

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
type ZeroClientListener = server.ZeroClientListener
type ZeroClientConnect = server.ZeroClientConnect

type UDPMessageProcesser = server.UDPMessageProcesser

type IPCServer = server.IPCServer
type TCPServer = server.TCPServer
type UDPServer = server.UDPServer

type TCPClient = server.TCPClient

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

const ZEROV1SERV_KEEPER = protocol.ZEROV1SERV_KEEPER
const ZEROV1SERV_CLIENT = protocol.ZEROV1SERV_CLIENT

type ZeroV1ServKeeper = protocol.ZeroV1ServKeeper
type ZeroV1MessageOperator = protocol.ZeroV1MessageOperator

type ZeroV1Message = protocol.ZeroV1Message

const WORKER_MONO_STATUS_READY = mfgrc.WORKER_MONO_STATUS_READY
const WORKER_MONO_STATUS_PENDING = mfgrc.WORKER_MONO_STATUS_PENDING
const WORKER_MONO_STATUS_EXECUTING = mfgrc.WORKER_MONO_STATUS_EXECUTING
const WORKER_MONO_STATUS_RETRYING = mfgrc.WORKER_MONO_STATUS_RETRYING
const WORKER_MONO_STATUS_COMPLETE = mfgrc.WORKER_MONO_STATUS_COMPLETE
const WORKER_MONO_STATUS_FAILED = mfgrc.WORKER_MONO_STATUS_FAILED
const WORKER_MONO_STATUS_REVOKE = mfgrc.WORKER_MONO_STATUS_REVOKE
const WORKER_MONO_STATUS_TIMEOUT = mfgrc.WORKER_MONO_STATUS_TIMEOUT

const WORKER_MONOGROUP_STATUS_READY = mfgrc.WORKER_MONOGROUP_STATUS_READY
const WORKER_MONOGROUP_STATUS_PENDING = mfgrc.WORKER_MONOGROUP_STATUS_PENDING
const WORKER_MONOGROUP_STATUS_EXECUTING = mfgrc.WORKER_MONOGROUP_STATUS_EXECUTING
const WORKER_MONOGROUP_STATUS_COMPLETE = mfgrc.WORKER_MONOGROUP_STATUS_COMPLETE
const WORKER_MONOGROUP_STATUS_FAILED = mfgrc.WORKER_MONOGROUP_STATUS_FAILED

type MfgrcMono = mfgrc.MfgrcMono
type MfgrcGroup = mfgrc.MfgrcGroup
type ZeroMfgrcMonoStore = mfgrc.ZeroMfgrcMonoStore
type ZeroMfgrcMonoEventListener = mfgrc.ZeroMfgrcMonoEventListener
type ZeroMfgrcKeeperOpts = mfgrc.ZeroMfgrcKeeperOpts

type ZeroMfgrcMono = mfgrc.ZeroMfgrcMono
type ZeroMfgrcFlux = mfgrc.ZeroMfgrcFlux
type ZeroMfgrcWorker = mfgrc.ZeroMfgrcWorker
type ZeroMfgrcKeeper = mfgrc.ZeroMfgrcKeeper
type ZeroMfgrcGroupKeeperOpts = mfgrc.ZeroMfgrcGroupKeeperOpts
type ZeroMfgrcGroupStore = mfgrc.ZeroMfgrcGroupStore
type ZeroMfgrcGroup = mfgrc.ZeroMfgrcGroup
type ZeroMfgrcGroupWorker = mfgrc.ZeroMfgrcGroupWorker
type ZeroMfgrcGroupKeeper = mfgrc.ZeroMfgrcGroupKeeper
type ZeroMfgrcGroupEventListener = mfgrc.ZeroMfgrcGroupEventListener

type ZeroMfgrcMonoActuator = mfgrc.ZeroMfgrcMonoActuator
type ZeroMfgrcGroupActuator = mfgrc.ZeroMfgrcGroupActuator
type ZeroMfgrcMonoQueueActuator = mfgrc.ZeroMfgrcMonoQueueActuator
