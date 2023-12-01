package zeroframework

import (
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/rocketmq"
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
