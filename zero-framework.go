package zeroframework

import (
	"github.com/0meet1/zero-framework/autohttpconf"
	"github.com/0meet1/zero-framework/autosqlconf"
	"github.com/0meet1/zero-framework/consul"
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/errdef"
	"github.com/0meet1/zero-framework/mfgrc"
	"github.com/0meet1/zero-framework/ossminiv2"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/protocol"
	"github.com/0meet1/zero-framework/rocketmq"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/signatures"
	"github.com/0meet1/zero-framework/structs"
)

type Time = structs.Time
type ZeroMetaDef = structs.ZeroMetaDef
type ZeroMetaPtr = structs.ZeroMetaPtr
type ZeroMeta = structs.ZeroMeta
type ZeroCoreStructs = structs.ZeroCoreStructs
type ZeroRequest = structs.ZeroRequest
type ZeroResponse = structs.ZeroResponse
type ZeroXsacAutoParser = structs.ZeroXsacAutoParser

var XreflectMeta = structs.XreflectMeta
var Xalloc = structs.Xalloc
var Md5 = structs.Md5
var Md5Bytes = structs.Md5Bytes
var HmacSha256 = structs.HmacSha256
var GenerateSUID = structs.GenerateSUID

type ZeroCoreProcessor = processors.ZeroCoreProcessor
type ZeroQuery = processors.ZeroQuery
type ZeroCondition = processors.ZeroCondition
type ZeroOrderBy = processors.ZeroOrderBy
type ZeroLimit = processors.ZeroLimit
type ZeroQueryOperation = processors.ZeroQueryOperation
type ZeroPostgresQueryOperation = processors.ZeroPostgresQueryOperation
type ZeroMysqlQueryOperation = processors.ZeroMysqlQueryOperation

const DATABASE_MYSQL = database.DATABASE_MYSQL
const DATABASE_POSTGRES = database.DATABASE_POSTGRES
const DATABASE_REDIS = database.DATABASE_REDIS
const DATABASE_SQLITE = database.DATABASE_SQLITE

type RedisKeeper = database.RedisKeeper
type EQueryRequest = database.EQueryRequest
type EQueryResponse = database.EQueryResponse
type EQuerySearch = database.EQuerySearch
type DataSource = database.DataSource
type SecureDataSource = database.SecureDataSource

var ElasticDatabase = database.ElasticDatabase
var RedisDatabase = database.RedisDatabase
var MySQLDatabase = database.MySQLDatabase
var CustomMySQLDatabase = database.CustomMySQLDatabase
var PostgresDatabase = database.PostgresDatabase
var CustomPostgresDatabase = database.CustomPostgresDatabase
var SQLiteDatabase = database.SQLiteDatabase
var CustomSQLiteDatabase = database.CustomSQLiteDatabase
var NewSQLiteTable = database.NewSQLiteTable

const ROCKETMQ_KEEPER = rocketmq.ROCKETMQ_KEEPER

type RocketmqKeeper = rocketmq.RocketmqKeeper
type MQNotifyMessage = rocketmq.MQNotifyMessage
type MQMessageObserver = rocketmq.MQMessageObserver

var XhttpResponseMaps = server.XhttpResponseMaps
var XhttpResponseDatas = server.XhttpResponseDatas
var XhttpResponseMessages = server.XhttpResponseMessages
var XhttpZeroRequest = server.XhttpZeroRequest
var XhttpZeroQuery = server.XhttpZeroQuery
var XhttpEachQuerys = server.XhttpEachQuerys
var XhttpUnitaryQuerys = server.XhttpUnitaryQuerys
var XhttpMysqlQueryOperation = server.XhttpMysqlQueryOperation
var XhttpPostgresQueryOperation = server.XhttpPostgresQueryOperation
var XhttpCompleteQueryOperation = server.XhttpCompleteQueryOperation
var XhttpQueryOptions = server.XhttpQueryOptions
var XhttpContainsOptions = server.XhttpContainsOptions
var XhttpEQuery = server.XhttpEQuery
var XhttpEQueryRequest = server.XhttpEQueryRequest
var XhttpURIParams = server.XhttpURIParams
var XhttpFromFileRequest = server.XhttpFromFileRequest
var XhttpKeyValueRequest = server.XhttpKeyValueRequest
var XhttpFuncHandle = server.XhttpFuncHandle
var XhttpHandle = server.XhttpHandle
var XhttpPerform = server.XhttpPerform

type XhttpFromFile = server.XhttpFromFile
type XhttpExecutor = server.XhttpExecutor

type XhttpInterceptor = server.XhttpInterceptor

var RunHttpServer = server.RunHttpServer
var RunInterceptor = server.RunInterceptor

type ZeroServ = server.ZeroServ
type ZeroDataChecker = server.ZeroDataChecker
type ZeroConnectBuilder = server.ZeroConnectBuilder
type ZeroConnect = server.ZeroConnect

type ZeroSocketConnect = server.ZeroSocketConnect
type ZeroSocketServer = server.ZeroSocketServer
type ZeroServerWatcher = server.ZeroServerWatcher
type ZeroClientListener = server.ZeroClientListener
type ZeroClientConnect = server.ZeroClientConnect

type UDPMessageProcesser = server.UDPMessageProcesser

type IPCServer = server.IPCServer
type TCPServer = server.TCPServer
type UDPServer = server.UDPServer

type TCPClient = server.TCPClient

const (
	Qos0 = server.Qos0
	Qos1 = server.Qos1
	Qos2 = server.Qos2

	FIXED_FLAG_NONE  = server.FIXED_FLAG_NONE
	FIXED_FLAG_Qos0s = server.FIXED_FLAG_Qos0s
	FIXED_FLAG_Qos1s = server.FIXED_FLAG_Qos1s
	FIXED_FLAG_Qos2s = server.FIXED_FLAG_Qos2s
)

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

type MqttMessageListener = server.MqttMessageListener
type MqttConnectBuilder = server.MqttConnectBuilder
type MqttConnect = server.MqttConnect
type MqttServer = server.MqttServer

var DefaultMqttChecker = server.DefaultMqttChecker

const ZEROKMSG_SERVER = protocol.ZEROKMSG_SERVER
const ZEROKMSG_CLIENT = protocol.ZEROKMSG_CLIENT

var RunKMessageServer = protocol.RunKMessageServer
var RunKMessageClient = protocol.RunKMessageClient

type ZeroKMessageServer = protocol.ZeroKMessageServer
type ZeroKMessageClient = protocol.ZeroKMessageClient
type ZeroKMessageOperator = protocol.ZeroKMessageOperator

type ZeroKMessage = protocol.ZeroKMessage
type ZeroKMessageConnect = protocol.ZeroKMessageConnect

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
type MfgrcXhttpExecutor = mfgrc.MfgrcXhttpExecutor

const (
	DSC_LOCK_TRUNK = consul.DSC_LOCK_TRUNK
)

type ZeroDCSMutex = consul.ZeroDCSMutex
type ZeroDCSMutexTrunk = consul.ZeroDCSMutexTrunk

const (
	ZEOR_XSAC_ENTRY_TYPE_TABLE0S          = structs.ZEOR_XSAC_ENTRY_TYPE_TABLE0S
	ZEOR_XSAC_ENTRY_TYPE_TABLE0FS         = structs.ZEOR_XSAC_ENTRY_TYPE_TABLE0FS
	ZEOR_XSAC_ENTRY_TYPE_COLUMN           = structs.ZEOR_XSAC_ENTRY_TYPE_COLUMN
	ZEOR_XSAC_ENTRY_TYPE_DROPCOLUMN       = structs.ZEOR_XSAC_ENTRY_TYPE_DROPCOLUMN
	ZEOR_XSAC_ENTRY_TYPE_KEY              = structs.ZEOR_XSAC_ENTRY_TYPE_KEY
	ZEOR_XSAC_ENTRY_TYPE_DROPKEY          = structs.ZEOR_XSAC_ENTRY_TYPE_DROPKEY
	ZEOR_XSAC_ENTRY_TYPE_PRIMARY_KEY      = structs.ZEOR_XSAC_ENTRY_TYPE_PRIMARY_KEY
	ZEOR_XSAC_ENTRY_TYPE_DROP_PRIMARY_KEY = structs.ZEOR_XSAC_ENTRY_TYPE_DROP_PRIMARY_KEY
	ZEOR_XSAC_ENTRY_TYPE_UNIQUE_KEY       = structs.ZEOR_XSAC_ENTRY_TYPE_UNIQUE_KEY
	ZEOR_XSAC_ENTRY_TYPE_DROP_UNIQUE_KEY  = structs.ZEOR_XSAC_ENTRY_TYPE_DROP_UNIQUE_KEY
	ZEOR_XSAC_ENTRY_TYPE_FOREIGN_KEY      = structs.ZEOR_XSAC_ENTRY_TYPE_FOREIGN_KEY
	ZEOR_XSAC_ENTRY_TYPE_DROP_FOREIGN_KEY = structs.ZEOR_XSAC_ENTRY_TYPE_DROP_FOREIGN_KEY
	ZEOR_XSAC_ENTRY_TYPE_YEAR_PARTITION   = structs.ZEOR_XSAC_ENTRY_TYPE_YEAR_PARTITION
	ZEOR_XSAC_ENTRY_TYPE_MONTH_PARTITION  = structs.ZEOR_XSAC_ENTRY_TYPE_MONTH_PARTITION
	ZEOR_XSAC_ENTRY_TYPE_DAY_PARTITION    = structs.ZEOR_XSAC_ENTRY_TYPE_DAY_PARTITION
	ZEOR_XSAC_ENTRY_TYPE_CUSTOM_PARTITION = structs.ZEOR_XSAC_ENTRY_TYPE_CUSTOM_PARTITION
)

const XSAC_AUTO_PARSER_KEEPER = autosqlconf.XSAC_AUTO_PARSER_KEEPER

type ZeroXsacDeclares = structs.ZeroXsacDeclares
type ZeroXsacEntry = structs.ZeroXsacEntry
type ZeroXsacField = structs.ZeroXsacField

type ZeroXsacProcessor = autosqlconf.ZeroXsacProcessor
type ZeroXsacPostgresProcessor = autosqlconf.ZeroXsacPostgresProcessor
type ZeroXsacMysqlProcessor = autosqlconf.ZeroXsacMysqlProcessor
type ZeroXsacKeeper = autosqlconf.ZeroXsacKeeper
type ZeroXsacAutoParserKeeper = autosqlconf.ZeroXsacAutoParserKeeper

const (
	XSAC_BE_INSERT = processors.XSAC_BE_INSERT
	XSAC_BE_UPDATE = processors.XSAC_BE_UPDATE
	XSAC_BE_DELETE = processors.XSAC_BE_DELETE

	XSAC_AF_INSERT = processors.XSAC_AF_INSERT
	XSAC_AF_UPDATE = processors.XSAC_AF_UPDATE
	XSAC_AF_DELETE = processors.XSAC_AF_DELETE
)

type ZeroXsacTrigger = structs.ZeroXsacTrigger
type ZeroXsacAutoProcessor = processors.ZeroXsacAutoProcessor
type ZeroXsacPostgresAutoProcessor = processors.ZeroXsacPostgresAutoProcessor
type ZeroXsacMysqlAutoProcessor = processors.ZeroXsacMysqlAutoProcessor

var XautoProcessor = autohttpconf.XautoProcessor

const (
	XSAC_DML_ADD       = autohttpconf.XSAC_DML_ADD
	XSAC_DML_UP        = autohttpconf.XSAC_DML_UP
	XSAC_DML_RM        = autohttpconf.XSAC_DML_RM
	XSAC_DML_TOMBSTONE = autohttpconf.XSAC_DML_TOMBSTONE
	XSAC_DML_RESTORE   = autohttpconf.XSAC_DML_RESTORE

	XSAC_HTTPFETCH_READY    = autohttpconf.XSAC_HTTPFETCH_READY
	XSAC_HTTPFETCH_ROW      = autohttpconf.XSAC_HTTPFETCH_ROW
	XSAC_HTTPFETCH_COMPLETE = autohttpconf.XSAC_HTTPFETCH_COMPLETE
)

type ZeroXsacXhttpDeclares = autohttpconf.ZeroXsacXhttpDeclares
type ZeroXsacHttpDMLTrigger = autohttpconf.ZeroXsacHttpDMLTrigger
type ZeroXsacHttpFetchTrigger = autohttpconf.ZeroXsacHttpFetchTrigger
type ZeroXsacHttpSearchTrigger = autohttpconf.ZeroXsacHttpSearchTrigger
type ZeroXsacXhttpStructs = autohttpconf.ZeroXsacXhttpStructs
type ZeroXsacCustomPartChecker = autohttpconf.ZeroXsacCustomPartChecker
type ZeroXsacXhttpApi = autohttpconf.ZeroXsacXhttpApi

var XahttpOpt = structs.XahttpOpt
var XahttpOptNoS = structs.XahttpOptNoS
var XahttpOptNoU = structs.XahttpOptNoU
var XahttpOptNoR = structs.XahttpOptNoR
var XahttpOptNoUR = structs.XahttpOptNoUR
var XahttpOptIO = structs.XahttpOptIO
var XahttpOptFO = structs.XahttpOptFO
var XahttpOptAll = structs.XahttpOptAll
var XsacPhysically = structs.XsacPhysically
var XsacTombstone = structs.XsacTombstone
var XsacTombstoneAndHistory = structs.XsacTombstoneAndHistory
var XsacTombstoneAndForce = structs.XsacTombstoneAndForce
var XsacTombstoneAndRestore = structs.XsacTombstoneAndRestore
var XsacTombstoneWhole = structs.XsacTombstoneWhole

type ZeroSignature = signatures.ZeroSignature
type OssminiV2Keeper = ossminiv2.OssminiV2Keeper

var Xfexists = structs.Xfexists
var Xfmake = structs.Xfmake
var Xfwrite = structs.Xfwrite
var Xfread = structs.Xfread
var Xfmove = structs.Xfmove

var DocOnSuccess = structs.DocOnSuccess

type ZeroExceptionDef = errdef.ZeroExceptionDef

const (
	EXCEPTION_KEEPER    = errdef.EXCEPTION_KEEPER
	EXCEPTION_AUTO_PROC = errdef.EXCEPTION_AUTO_PROC
	EXCEPTION_OPERATION = errdef.EXCEPTION_OPERATION

	ES00500 = errdef.ES00500
)

type ZeroExceptionKeeper = errdef.ZeroExceptionKeeper
