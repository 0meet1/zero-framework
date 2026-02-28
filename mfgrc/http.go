package mfgrc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/0meet1/zero-framework/autohttpconf"
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
	"github.com/gofrs/uuid"
)

type MfgrcXhttpProcessor struct {
	processors.ZeroCoreProcessor

	executor *MfgrcXhttpExecutor
}

func (processor *MfgrcXhttpProcessor) mysqlGroupMonos() string {
	mono := processor.executor.newMono().(structs.ZeroXsacDeclares)
	group := processor.executor.newGroup()
	return fmt.Sprintf(`
		SELECT 
			t.* 
		FROM 
			(SELECT 
				* 
			FROM 
				%s 
			WHERE 
				create_time 
			BETWEEN 
				?
			AND 
				?) t, 
			%s b 
		WHERE 
			b.group_id = ?
		AND 
			t.id = b.mono_id
	`, mono.XsacTableName(), group.XLinkTable())
}

func (processor *MfgrcXhttpProcessor) postgresGroupMonos() string {
	mono := processor.executor.newMono().(structs.ZeroXsacDeclares)
	group := processor.executor.newGroup()
	return fmt.Sprintf(`
		SELECT 
			t.* 
		FROM 
			(SELECT 
				* 
			FROM 
				%s 
			WHERE 
				create_time 
			BETWEEN 
				$1
			AND 
				$2) t, 
			%s b 
		WHERE 
			b.group_id = $3
		AND 
			t.id = b.mono_id
	`, mono.XsacTableName(), group.XLinkTable())
}

func (processor *MfgrcXhttpProcessor) fetchGroupMonos(groupID string) ([]MfgrcMono, error) {
	FETCH_MONO_SQL := processor.postgresGroupMonos()
	if processor.executor.DataSource == database.DATABASE_MYSQL {
		FETCH_MONO_SQL = processor.mysqlGroupMonos()
	}
	startTime, endTime, err := structs.MonthDurationString(time.Now(), "2006-01-02 15:04:05")
	if err != nil {
		return nil, err
	}
	rows, err := processor.PreparedStmt(FETCH_MONO_SQL).Query(startTime, endTime, groupID)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	rowsmap := processor.Parser(rows)
	monos := make([]MfgrcMono, len(rowsmap))
	for i, row := range rowsmap {
		data, err := structs.XautoLoad(processor.executor.MonoType, row)
		if err != nil {
			panic(err)
		}
		data.(MfgrcMono).ThisDef(data)
		monos[i] = data.(MfgrcMono)
	}
	return monos, nil
}

func (processor *MfgrcXhttpProcessor) mysqlMonoByMonoID() string {
	mono := processor.executor.newMono().(structs.ZeroXsacDeclares)
	return fmt.Sprintf("SELECT * FROM %s WHERE create_time BETWEEN ? AND ? AND mono_id = ?", mono.XsacTableName())
}

func (processor *MfgrcXhttpProcessor) postgresMonoByMonoID() string {
	mono := processor.executor.newMono().(structs.ZeroXsacDeclares)
	return fmt.Sprintf("SELECT * FROM %s WHERE create_time BETWEEN $1 AND $2 AND mono_id = $3", mono.XsacTableName())
}

func (processor *MfgrcXhttpProcessor) fetchMonoByMonoID(startTime time.Time, endTime time.Time, monoID string) []interface{} {
	FETCH_MONO_SQL := processor.postgresMonoByMonoID()
	if processor.executor.DataSource == database.DATABASE_MYSQL {
		FETCH_MONO_SQL = processor.mysqlMonoByMonoID()
	}

	rows, err := processor.PreparedStmt(FETCH_MONO_SQL).Query(
		startTime.Format("2006-01-02 15:04:05"),
		endTime.Format("2006-01-02 15:04:05"),
		monoID,
	)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		panic(err)
	}
	rowsmap := processor.Parser(rows)
	monos := make([]interface{}, len(rowsmap))
	for i, row := range rowsmap {
		data, err := structs.XautoLoad(processor.executor.MonoType, row)
		if err != nil {
			panic(err)
		}
		monos[i] = data.(MfgrcMono)
	}
	return monos
}

type MfgrcXhttpExecutor struct {
	Prefix     string
	DataSource string

	GroupKeeper string
	GroupType   reflect.Type
	GroupStore  ZeroMfgrcGroupStore

	OnGroupReady   func(MfgrcGroup, map[string]any)
	OnGroupSuccess func(MfgrcGroup, map[string]any) any
	OnGroupFailed  func(MfgrcGroup, map[string]any) string

	MonoKeeper string
	MonoType   reflect.Type
	MonoStore  ZeroMfgrcMonoStore

	OnMonoReady   func(MfgrcMono, map[string]any)
	OnMonoSuccess func(MfgrcMono) any
	OnMonoFailed  func(MfgrcMono, map[string]any) string

	GenNmonoId func(MfgrcMono) error

	IncidentBeforeState func(map[string]any) error

	IncidentBeforeMonoPerform func(MfgrcMono) error
	IncidentBeforeMonoFailed  func(MfgrcMono) error

	IncidentBeforeGroupFailed  func(MfgrcGroup) error
	IncidentBeforeGroupPerform func(MfgrcGroup) error
}

func (xhttpExecutor *MfgrcXhttpExecutor) newGroup() MfgrcGroup {
	if xhttpExecutor.GroupType == nil {
		return nil
	}
	_meta := xhttpExecutor.GroupType
	if _meta.Kind() == reflect.Ptr {
		_meta = _meta.Elem()
	}
	return reflect.New(_meta).Interface().(MfgrcGroup)
}

func (xhttpExecutor *MfgrcXhttpExecutor) newMono() MfgrcMono {
	if xhttpExecutor.MonoType == nil {
		return nil
	}
	_meta := xhttpExecutor.MonoType
	if _meta.Kind() == reflect.Ptr {
		_meta = _meta.Elem()
	}
	return reflect.New(_meta).Interface().(MfgrcMono)
}

func (xhttpExecutor *MfgrcXhttpExecutor) uXmonoComplete(xRequest *structs.ZeroRequest, keeper *ZeroMfgrcKeeper, mono MfgrcMono) error {
	if len(xRequest.Querys) != 1 {
		return errors.New(" no support multiple tasks or task is empty ")
	}

	bytes, err := json.Marshal(xRequest.Querys[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, mono)
	if err != nil {
		return err
	}
	if strings.TrimSpace(mono.XuniqueCode()) == "" {
		return errors.New(" `uniqueCode` is empty ")
	}
	err = keeper.Check(mono)
	if err != nil {
		return err
	}
	mono.ThisDef(mono)
	return keeper.Check(mono)
}

func (xhttpExecutor *MfgrcXhttpExecutor) uXmonoPerformed(
	writer http.ResponseWriter,
	xRequest *structs.ZeroRequest,
	keeper *ZeroMfgrcKeeper,
	mono MfgrcMono,
	onReady func(MfgrcMono, map[string]any),
	onSuccess func(MfgrcMono) any,
	onFailed func(MfgrcMono, map[string]any) string) {

	expands := make(map[string]interface{})
	if onReady != nil {
		onReady(mono, expands)
	} else {
		expands["monoId"] = mono.XmonoId()
	}

	waittime, ok := xRequest.Expands["waittime"]
	if ok && waittime.(float64) > 0 {
		act := &ZeroMfgrcMonoActuator{Keeper: keeper}
		select {
		case err := <-act.Exec(mono):
			if err != nil {
				expands["state"] = "error"
				if onFailed != nil {
					_err := onFailed(mono, expands)
					if _err != "" {
						err = errors.New(_err)
					}
				}
				server.XhttpResponseMessages(writer, 500, err.Error())
				return
			} else {
				expands["state"] = "success"
				if onSuccess != nil {
					result := onSuccess(mono)
					if result != nil {
						expands["result"] = result
					}
				}

			}
		case <-time.After(time.Second * time.Duration(waittime.(float64))):
			expands["state"] = "timeout"
		}
	} else {
		go keeper.AddMono(mono)
		expands["state"] = "created"
	}
	server.XhttpResponseDatas(writer, 200, "success", make([]interface{}, 0), expands)
}

func (xhttpExecutor *MfgrcXhttpExecutor) uXgroupComplete(xRequest *structs.ZeroRequest, keeper *ZeroMfgrcGroupKeeper, group MfgrcGroup) error {
	if len(xRequest.Querys) != 1 {
		return errors.New(" no support multiple tasks or task is empty ")
	}
	_group := make(map[string]any)
	jsonbytes, err := json.Marshal(xRequest.Querys[0])
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonbytes, &_group)
	if err != nil {
		return err
	}

	monos := make([]MfgrcMono, 0)
	jsonmonos, ok := _group["monos"]
	if ok {
		delete(_group, "monos")
		for _, jsonmono := range jsonmonos.([]interface{}) {
			_jsonbytes, err := json.Marshal(jsonmono)
			if err != nil {
				return err
			}
			_mono := xhttpExecutor.newMono()
			err = json.Unmarshal(_jsonbytes, _mono)
			if err != nil {
				return err
			}
			monos = append(monos, _mono)
		}
	}

	jsonbytes, err = json.Marshal(_group)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonbytes, group)
	if err != nil {
		return err
	}

	reflect.ValueOf(group).Elem().FieldByName("Monos").Set(reflect.ValueOf(monos))

	if strings.TrimSpace(group.XuniqueCode()) == "" {
		return errors.New(" `uniqueCode` is empty ")
	}

	err = keeper.Check(group)
	if err != nil {
		return err
	}
	group.ThisDef(group)
	group.Xmonos()
	return nil
}

func (xhttpExecutor *MfgrcXhttpExecutor) uXgroupPerformed(
	writer http.ResponseWriter,
	xRequest *structs.ZeroRequest,
	keeper *ZeroMfgrcGroupKeeper,
	group MfgrcGroup,
	onReady func(MfgrcGroup, map[string]any),
	onSuccess func(MfgrcGroup, map[string]any) any,
	onFailed func(MfgrcGroup, map[string]any) string) {

	expands := make(map[string]interface{})
	if onReady != nil {
		onReady(group, expands)
	} else {
		expands["groupId"] = group.XgroupId()
	}

	waittime, ok := xRequest.Expands["waittime"]
	if ok && waittime.(float64) > 0 {
		act := &ZeroMfgrcGroupActuator{Keeper: keeper}
		select {
		case err := <-act.Exec(group):
			if err != nil {
				expands["state"] = "error"
				if onFailed != nil {
					_err := onFailed(group, expands)
					if _err != "" {
						err = errors.New(_err)
					}
				}
				server.XhttpResponseDatas(writer, 500, err.Error(), make([]interface{}, 0), expands)
				return
			} else {
				expands["state"] = "success"
				if onSuccess != nil {
					result := onSuccess(group, expands)
					if result != nil {
						expands["result"] = result
					}
				}

			}
		case <-time.After(time.Second * time.Duration(waittime.(float64))):
			expands["state"] = "timeout"
		}
	} else {
		go keeper.AddGroup(group)
		expands["state"] = "created"
	}
	server.XhttpResponseDatas(writer, 200, "success", make([]interface{}, 0), expands)
}

func (xhttpExecutor *MfgrcXhttpExecutor) groupc(writer http.ResponseWriter, req *http.Request) {
	group := xhttpExecutor.newGroup()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			if xhttpExecutor.IncidentBeforeMonoFailed != nil {
				_err := xhttpExecutor.IncidentBeforeGroupFailed(group)
				if _err != nil {
					global.Logger().ErrorS(_err)
				}
			}
			server.XhttpResponseMessages(writer, 500, err.(error).Error())
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	k := global.Value(xhttpExecutor.GroupKeeper).(*ZeroMfgrcGroupKeeper)
	err = xhttpExecutor.uXgroupComplete(xRequest, k, group)
	if err != nil {
		panic(err)
	}

	err = group.Ready(xhttpExecutor.GroupStore)
	if err != nil {
		panic(err)
	}

	if xhttpExecutor.IncidentBeforeGroupPerform != nil {
		err = xhttpExecutor.IncidentBeforeGroupPerform(group)
		if err != nil {
			panic(err)
		}
	}

	xhttpExecutor.uXgroupPerformed(writer, xRequest, k, group, xhttpExecutor.OnGroupReady, xhttpExecutor.OnGroupSuccess, xhttpExecutor.OnGroupFailed)
}

func (xhttpExecutor *MfgrcXhttpExecutor) defaultNmonoId(mono MfgrcMono) {
	uid, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	reflect.ValueOf(mono).Elem().FieldByName("MonoID").SetString(fmt.Sprintf("%s-%s-%s-%s", mono.Xoperator(), mono.Xoption(), mono.XuniqueCode(), uid.String()))
}

func (xhttpExecutor *MfgrcXhttpExecutor) monoc(writer http.ResponseWriter, req *http.Request) {
	mono := xhttpExecutor.newMono()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			if xhttpExecutor.IncidentBeforeMonoFailed != nil {
				_err := xhttpExecutor.IncidentBeforeMonoFailed(mono)
				if _err != nil {
					global.Logger().ErrorS(_err)
				}
			}
			_err := mono.Revoke()
			if _err != nil {
				global.Logger().ErrorS(_err)
				err = _err
			}
			_err = mono.Delete()
			if _err != nil {
				global.Logger().ErrorS(_err)
				err = _err
			}
			server.XhttpResponseMessages(writer, 500, err.(error).Error())
		}
	}()
	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}
	k := global.Value(xhttpExecutor.MonoKeeper).(*ZeroMfgrcKeeper)
	err = xhttpExecutor.uXmonoComplete(xRequest, k, mono)
	if err != nil {
		panic(err)
	}

	if xhttpExecutor.GenNmonoId != nil {
		err = xhttpExecutor.GenNmonoId(mono)
		if err != nil {
			panic(err)
		}
	} else {
		xhttpExecutor.defaultNmonoId(mono)
	}

	err = mono.Ready(k, xhttpExecutor.MonoStore)
	if err != nil {
		panic(err)
	}

	if xhttpExecutor.IncidentBeforeMonoPerform != nil {
		err = xhttpExecutor.IncidentBeforeMonoPerform(mono)
		if err != nil {
			panic(err)
		}
	}

	xhttpExecutor.uXmonoPerformed(writer, xRequest, k, mono, xhttpExecutor.OnMonoReady, xhttpExecutor.OnMonoSuccess, xhttpExecutor.OnMonoFailed)
}

func (xhttpExecutor *MfgrcXhttpExecutor) revoke(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttpExecutor.DataSource).(database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, err.(error).Error())
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	processor := &MfgrcXhttpProcessor{}
	processor.Build(transaction)

	expands := make(map[string]interface{})
	for i := 0; i < len(xRequest.Querys); i++ {
		bytes, err := json.Marshal(xRequest.Querys[i])
		if err != nil {
			panic(err)
		}
		mono := xhttpExecutor.newMono()
		err = json.Unmarshal(bytes, mono)
		if err != nil {
			panic(err)
		}

		refCreateTime := reflect.ValueOf(mono).Elem().FieldByName("CreateTime")
		if refCreateTime.IsNil() {
			panic(fmt.Errorf(" `mono.createTime` not found "))
		}
		startTime, endTime, err := structs.MonthDuration(refCreateTime.Interface().(*structs.Time).Time())
		if err != nil {
			panic(err)
		}

		monos := processor.fetchMonoByMonoID(startTime, endTime, mono.XmonoId())
		if len(monos) <= 0 {
			expands[mono.XmonoId()] = fmt.Sprintf("mono %s not found", mono.XmonoId())
		} else {
			_mono := monos[0].(MfgrcMono)
			err = global.Value(xhttpExecutor.MonoKeeper).(*ZeroMfgrcKeeper).RevokeMono(_mono)
			if err != nil {
				expands[_mono.XmonoId()] = err.Error()
			} else {
				expands[_mono.XmonoId()] = "success"
			}
		}
	}

	server.XhttpResponseDatas(writer, 200, "success", make([]interface{}, 0), expands)
}

func (xhttpExecutor *MfgrcXhttpExecutor) state(writer http.ResponseWriter, _ *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			server.XhttpResponseMessages(writer, 500, err.(error).Error())
		}
	}()

	expands := make(map[string]interface{})

	if xhttpExecutor.IncidentBeforeState != nil {
		err := xhttpExecutor.IncidentBeforeState(expands)
		if err != nil {
			panic(err)
		}
	}

	monoKeeper := global.Value(xhttpExecutor.MonoKeeper)
	if monoKeeper != nil {
		fluxexp, err := monoKeeper.(*ZeroMfgrcKeeper).Export()
		if err != nil {
			panic(err)
		}
		expands["fluxs"] = fluxexp
	}

	groupKeeper := global.Value(xhttpExecutor.GroupKeeper)
	if groupKeeper != nil {
		groupexp, err := groupKeeper.(*ZeroMfgrcGroupKeeper).Export()
		if err != nil {
			panic(err)
		}
		expands["groups"] = groupexp
	}

	server.XhttpResponseDatas(writer, 200, "success", make([]interface{}, 0), expands)
}

func (xhttpExecutor *MfgrcXhttpExecutor) checkzone(xRequest *structs.ZeroRequest) (string, string) {
	if len(xRequest.Querys) <= 0 {
		panic(errors.New("missing necessary parameter `query[0]`"))
	}

	if xRequest.Expands == nil {
		panic(errors.New("missing necessary parameter `expands.zone`"))
	}

	zone, ok := xRequest.Expands["zone"]
	if !ok {
		panic(errors.New("missing necessary parameter `expands.zone`"))
	}

	date, err := time.Parse("2006-01-02", zone.(string))
	if err != nil {
		panic(err)
	}
	startTime, endTime, err := structs.MonthDurationString(date, "2006-01-02 15:04:05")
	if err != nil {
		panic(err)
	}
	return startTime, endTime
}

var httpCompleteQueryOperation = func(xRequest *structs.ZeroRequest, xProcessor processors.ZeroQueryOperation, tableName string) (processors.ZeroQueryOperation, *processors.ZeroQuery, error) {
	xQuery, err := server.XhttpZeroQuery(xRequest)
	if err != nil {
		return nil, nil, err

	}
	xProcessor.AddQuery(xQuery)
	xProcessor.AddTableName(tableName)
	return xProcessor, xQuery, nil
}

func (xhttpExecutor *MfgrcXhttpExecutor) grouphistory(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttpExecutor.DataSource).(database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, err.(error).Error())
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	startTime, endTime := xhttpExecutor.checkzone(xRequest)
	declares := xhttpExecutor.newGroup()

	xOperation, _, err := httpCompleteQueryOperation(xRequest, declares.(autohttpconf.ZeroXsacXhttpDeclares).XhttpQueryOperation(), declares.(structs.ZeroXsacDeclares).XsacTableName())
	if err != nil {
		panic(err)
	}
	xOperation.Build(transaction)
	xOperation.AppendCondition(fmt.Sprintf("create_time BETWEEN '%s' AND '%s'", startTime, endTime))
	datas, expands := xOperation.Exec()

	processor := &MfgrcXhttpProcessor{executor: xhttpExecutor}
	processor.Build(transaction)

	groups := make([]map[string]interface{}, len(datas))
	for i, data := range datas {
		group, err := structs.XautoLoad(xhttpExecutor.GroupType, data)
		if err != nil {
			panic(err)
		}
		group.(MfgrcGroup).ThisDef(group)
		groupexp, err := group.(MfgrcGroup).Export()
		if err != nil {
			panic(err)
		}

		if server.XhttpContainsOptions(xRequest, "mono") {
			monos, err := processor.fetchGroupMonos(reflect.ValueOf(group).Elem().FieldByName("ID").String())
			if err != nil {
				panic(err)
			}

			_monos := make([]map[string]any, 0)
			for _, m := range monos {
				_expm, err := m.Export()
				if err != nil {
					panic(err)
				}
				_monos = append(_monos, _expm)
			}

			groupexp["monos"] = _monos
		}
		groups[i] = groupexp
	}
	server.XhttpResponseMaps(writer, 200, "success", groups, expands)
}

func (xhttpExecutor *MfgrcXhttpExecutor) monohistory(writer http.ResponseWriter, req *http.Request) {
	transaction := global.Value(xhttpExecutor.DataSource).(database.DataSource).Transaction()
	defer func() {
		err := recover()
		if err != nil {
			global.Logger().ErrorS(err.(error))
			transaction.Rollback()
			server.XhttpResponseMessages(writer, 500, err.(error).Error())
		} else {
			transaction.Commit()
		}
	}()

	xRequest, err := server.XhttpZeroRequest(req)
	if err != nil {
		panic(err)
	}

	startTime, endTime := xhttpExecutor.checkzone(xRequest)
	declares := xhttpExecutor.newMono()
	xOperation, _, err := httpCompleteQueryOperation(xRequest, declares.(autohttpconf.ZeroXsacXhttpDeclares).XhttpQueryOperation(), declares.(structs.ZeroXsacDeclares).XsacTableName())
	if err != nil {
		panic(err)
	}

	xOperation.Build(transaction)
	xOperation.AppendCondition(fmt.Sprintf("create_time BETWEEN '%s' AND '%s'", startTime, endTime))
	datas, expands := xOperation.Exec()

	monos := make([]map[string]interface{}, len(datas))
	for i, data := range datas {
		mono, err := structs.XautoLoad(xhttpExecutor.MonoType, data)
		if err != nil {
			panic(err)
		}
		mono.(MfgrcMono).ThisDef(mono)
		monos[i], err = mono.(MfgrcMono).Export()
		if err != nil {
			panic(err)
		}
	}
	server.XhttpResponseMaps(writer, 200, "success", monos, expands)
}

func (xhttpExecutor *MfgrcXhttpExecutor) Exports() []*server.XhttpExecutor {
	executors := make([]*server.XhttpExecutor, 0)
	prefix := ""
	if strings.TrimSpace(xhttpExecutor.Prefix) != "" {
		prefix = strings.TrimSpace(xhttpExecutor.Prefix)
		if prefix[:1] == "/" {
			prefix = prefix[1:]
		}
		if prefix[len(prefix)-1:] != "/" {
			prefix = fmt.Sprintf("%s/", prefix)
		}
	}

	if xhttpExecutor.GroupType != nil {
		executors = append(executors, server.XhttpFuncHandle(xhttpExecutor.groupc, fmt.Sprintf("%sworker/group", prefix)))
		executors = append(executors, server.XhttpFuncHandle(xhttpExecutor.grouphistory, fmt.Sprintf("%shistory/group", prefix)))
	}

	if xhttpExecutor.MonoType != nil {
		executors = append(executors, server.XhttpFuncHandle(xhttpExecutor.monoc, fmt.Sprintf("%sworker/push", prefix)))
		executors = append(executors, server.XhttpFuncHandle(xhttpExecutor.monohistory, fmt.Sprintf("%shistory/mono", prefix)))
		executors = append(executors, server.XhttpFuncHandle(xhttpExecutor.revoke, fmt.Sprintf("%sworker/revoke", prefix)))
	}

	executors = append(executors, server.XhttpFuncHandle(xhttpExecutor.state, fmt.Sprintf("%sworker/state", prefix)))
	return executors
}
