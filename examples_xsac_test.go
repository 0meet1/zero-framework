package zeroframework_test

import (
	"reflect"

	"github.com/0meet1/zero-framework/autosqlconf"
	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/processors"
	"github.com/0meet1/zero-framework/server"

	x0meet1 "github.com/0meet1/zero-framework"
	x0meet1st "github.com/0meet1/zero-framework/structs"
)

type UgLiMediaArea struct {
	x0meet1.ZeroXsacXhttpStructs

	Name string `json:"name,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(128),NULL"`

	Type   string `json:"type,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(64),NULL"`
	Access string `json:"access,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(64),NULL"`
	// Polygon *UgLiPolygon `json:"polygon,omitempty"`

	// Media *UgLiMedia `json:"media,omitempty" xhttpopt:"OX" xsacname:"media_id" xsacfield:"-" xsacprop:"NO,UUID,NULL" xsackey:"foreign,zeroaikit_media,id"`
	Media *UgLiMedia `json:"media,omitempty" xhttpopt:"OX" xsacfield:"-" xsacref:"zeroaikit_media_area,area_id,media_id"`

	// Media *UgLiMedia `json:"media,omitempty" xhttpopt:"OX" xsacname:"media_id" xsacfield:"-" xsacprop:"NO,VARCHAR(36),NULL" xsackey:"foreign,zeroaikit_media,id"`
	// Media *UgLiMedia `json:"media,omitempty" xhttpopt:"OX" xsacfield:"-" xsacref:"zeroaikit_media_area,area_id,media_id"`

	// TargetDefs   []*UgLiTargetDef   `json:"targetDefs,omitempty"`
	// BehaviorDefs []*UgLiBehaviorDef `json:"behaviorDefs,omitempty"`
}

// func (_ *UgLiMediaArea) XsacPrimaryType() string { return "VARCHAR(36)" }
func (_ *UgLiMediaArea) XhttpPath() string { return "mediaarea" }

func (_ *UgLiMediaArea) XsacDbName() string { return "facestore" }

// func (_ *UgLiMediaArea) XsacDbName() string    { return "zeroframework" }
func (_ *UgLiMediaArea) XsacTableName() string { return "zeroaikit_area" }
func (_ *UgLiMediaArea) XsacDeleteOpt() byte   { return x0meet1.XsacTombstoneWhole() }

// func (_ *UgLiMediaArea) XhttpAutoProc() x0meet1.ZeroXsacAutoProcessor {
// 	return processors.NewXsacMysqlProcessor()
// }
// func (_ *UgLiMediaArea) XhttpQueryOperation() processors.ZeroQueryOperation {
// 	return &processors.ZeroMysqlQueryOperation{}
// }

func (_ *UgLiMediaArea) XhttpAutoProc() x0meet1.ZeroXsacAutoProcessor {
	return processors.NewXsacPostgresProcessor()
}
func (_ *UgLiMediaArea) XhttpQueryOperation() processors.ZeroQueryOperation {
	return &processors.ZeroPostgresQueryOperation{}
}

func (mediaArea *UgLiMediaArea) LoadRowData(rowmap map[string]interface{}) {
	mediaArea.ZeroXsacXhttpStructs.LoadRowData(rowmap)

	mediaArea.Name = x0meet1st.ParseStringField(rowmap, "name")
	mediaArea.Type = x0meet1st.ParseStringField(rowmap, "type")
	mediaArea.Access = x0meet1st.ParseStringField(rowmap, "access")
	x0meet1st.ParseIfExists(rowmap, "media_id", func(i interface{}) error {
		mediaArea.Media = &UgLiMedia{
			ZeroXsacXhttpStructs: x0meet1.ZeroXsacXhttpStructs{
				ZeroCoreStructs: x0meet1.ZeroCoreStructs{
					ID: x0meet1st.ParseStringField(rowmap, "media_id"),
				},
			},
		}
		return nil
	})
}

type UgLiMedia struct {
	x0meet1.ZeroXsacXhttpStructs

	Name string `json:"name,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(128),NULL" xsacindex:"unique"`

	Status string `json:"status,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(64),NULL"`
	Type   string `json:"type,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(64),NULL"`
	Fps    int    `json:"fps,omitempty" xhttpopt:"OO" xsacprop:"NO,INT,NULL"`
	// Areas  []*UgLiMediaArea `json:"areas,omitempty" xhttpopt:"OX" xsacchild:"Media"`
	Areas []*UgLiMediaArea `json:"areas,omitempty" xhttpopt:"OX" xsacchild:"Media" xsacref:"zeroaikit_media_area,media_id,area_id,inspect"`
}

// func (_ *UgLiMedia) XsacPrimaryType() string { return "VARCHAR(36)" }

func (_ *UgLiMedia) XhttpPath() string { return "media" }

func (_ *UgLiMedia) XsacDbName() string { return "facestore" }

// func (_ *UgLiMedia) XsacDbName() string    { return "zeroframework" }
func (_ *UgLiMedia) XsacTableName() string { return "zeroaikit_media" }
func (_ *UgLiMedia) XsacDeleteOpt() byte   { return x0meet1.XsacTombstoneWhole() }

// func (_ *UgLiMedia) XhttpAutoProc() x0meet1.ZeroXsacAutoProcessor {
// 	return processors.NewXsacMysqlProcessor()
// }
// func (_ *UgLiMedia) XhttpQueryOperation() processors.ZeroQueryOperation {
// 	return &processors.ZeroMysqlQueryOperation{}
// }

func (_ *UgLiMedia) XhttpAutoProc() x0meet1.ZeroXsacAutoProcessor {
	return processors.NewXsacPostgresProcessor()
}
func (_ *UgLiMedia) XhttpQueryOperation() processors.ZeroQueryOperation {
	return &processors.ZeroPostgresQueryOperation{}
}

func (media *UgLiMedia) LoadRowData(rowmap map[string]interface{}) {
	media.ZeroXsacXhttpStructs.LoadRowData(rowmap)

	media.Name = x0meet1st.ParseStringField(rowmap, "name")
	media.Status = x0meet1st.ParseStringField(rowmap, "status")
	media.Type = x0meet1st.ParseStringField(rowmap, "type")
	media.Fps = x0meet1st.ParseIntField(rowmap, "fps")
}

type UgLiEvent struct {
	x0meet1.ZeroXsacXhttpStructs

	OccurTime *x0meet1.Time `json:"occurTime,omitempty" xhttpopt:"OO" xsacprop:"NO,TIMESTAMPTZ,NULL"`
	// OccurTime *x0meet1.Time `json:"occurTime,omitempty" xhttpopt:"OO" xsacprop:"NO,DATETIME,NULL" xsackey:"unique"`
}

// func (_ *UgLiEvent) XsacPrimaryType() string { return "VARCHAR(36)" }
func (_ *UgLiEvent) XhttpPath() string { return "event" }

func (_ *UgLiEvent) XsacDbName() string { return "facestore" }

// func (_ *UgLiEvent) XsacDbName() string    { return "zeroframework" }
func (_ *UgLiEvent) XsacTableName() string { return "zeroaikit_event" }

// func (_ *UgLiEvent) XsacPartition() string { return x0meet1st.XSAC_PARTITION_MONTH }
// func (_ *UgLiEvent) XhttpAutoProc() x0meet1.ZeroXsacAutoProcessor {
// 	return processors.NewXsacMysqlProcessor()
// }
// func (_ *UgLiEvent) XhttpQueryOperation() processors.ZeroQueryOperation {
// 	return &processors.ZeroMysqlQueryOperation{}
// }

func (_ *UgLiEvent) XhttpAutoProc() x0meet1.ZeroXsacAutoProcessor {
	return processors.NewXsacPostgresProcessor()
}
func (_ *UgLiEvent) XhttpQueryOperation() processors.ZeroQueryOperation {
	return &processors.ZeroPostgresQueryOperation{}
}

// func (_ *UgLiEvent) XhttpDMLTrigger() x0meet1.ZeroXsacHttpDMLTrigger { return &UgLiEventDMLTrigger{} }

func (eve *UgLiEvent) LoadRowData(rowmap map[string]interface{}) {
	eve.ZeroXsacXhttpStructs.LoadRowData(rowmap)
	eve.OccurTime = x0meet1st.ParseDateField(rowmap, "occur_time")
}

// type UgLiEventDMLTrigger struct{}

// func (_ *UgLiEventDMLTrigger) On(opt string, timing string, xRequest *x0meet1.ZeroRequest, datas ...interface{}) error {
// 	fmt.Println("UgLiEventDMLTrigger.On", opt, timing, xRequest, datas)

// 	if len(datas) > 0 {
// 		evet := datas[0].(*UgLiEvent)
// 		fmt.Println(evet.OccurTime.Time())
// 	}
// 	fmt.Println("-000")
// 	return nil
// }

func test() {
	global.Run("ugliserv", func() {
		// database.InitMYSQLDatabase()
		database.InitPostgresDatabase()

		xscaKeeper := autosqlconf.NewKeeper(
			// reflect.TypeOf(x0meet1.ZeroXsacMysqlProcessor{}),
			reflect.TypeOf(x0meet1.ZeroXsacPostgresProcessor{}),
			reflect.TypeOf(&UgLiMedia{}),
			reflect.TypeOf(&UgLiMediaArea{}),
			reflect.TypeOf(&UgLiEvent{})).
			// DataSource(x0meet1.DATABASE_MYSQL)
			DataSource(x0meet1.DATABASE_POSTGRES)
		xscaKeeper.RunKeeper()

		executors := make([]*server.XhttpExecutor, 0)
		executors = append(executors, xscaKeeper.Exports()...)
		go server.RunHttpServer(executors...)
	})
}
