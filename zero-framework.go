package main

import (
	"zero-framework/database"
	"zero-framework/processors"
	"zero-framework/rocketmq"
	"zero-framework/structs"
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

type MQNotifyMessage = rocketmq.MQNotifyMessage
type MQMessageObserver = rocketmq.MQMessageObserver
