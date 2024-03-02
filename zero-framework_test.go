package zeroframework_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	x0meet1 "github.com/0meet1/zero-framework"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacTestPerson struct {
	x0meet1.ZeroXsacXhttpStructs

	IdCard  string                      `json:"idCard,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(32),NULL" xsackey:"unique"`
	Records []*ZeroXsacTestPersonRecord `json:"records,omitempty" xhttpopt:"OX" xsacchild:"Person"`
	// Records []*ZeroXsacTestPersonRecord `json:"records,omitempty" xhttpopt:"XXO" xsacref:"test_person_ref,person_id,record_id,inspect"`
}

func (person *ZeroXsacTestPerson) XsacDbName() string {
	return "test"

}

func (person *ZeroXsacTestPerson) XsacTableName() string {
	return "test_person"
}

type ZeroXsacTestPersonRecord struct {
	x0meet1.ZeroXsacXhttpStructs

	// Person *ZeroXsacTestPerson `json:"person,omitempty" xhttpopt:"OX" xsacname:"person_id" xsacfield:"ID" xsacprop:"NO,VARCHAR(32),NULL" xsackey:"foreign,test_person,id"`
	Person   *ZeroXsacTestPerson `json:"person,omitempty" xhttpopt:"XXO" xsacfield:"-" xsacref:"test_person_ref,record_id,person_id"`
	Contents string              `json:"contents,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(256),NULL"`
}

func (personr *ZeroXsacTestPersonRecord) XsacDbName() string {
	return "test"
}

func (personr *ZeroXsacTestPersonRecord) XsacTableName() string {
	return "test_person_record"
}

type Param struct {
	Start structs.Time `json:"start"`
	End   structs.Time `json:"end"`
}

func TestZeroCoreStructs(t *testing.T) {

	fmt.Println(structs.YearDurationString(time.Now(), "2006-01-02 15:04:05"))
	fmt.Println(structs.MonthDurationString(time.Now(), "2006-01-02 15:04:05"))
	fmt.Println(structs.DayDurationString(time.Now(), "2006-01-02 15:04:05"))

	xt := &ZeroXsacTestPersonRecord{}
	xt.InitDefault()

	jsonbytes, err := json.Marshal(xt)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(jsonbytes))

	err = json.Unmarshal(jsonbytes, xt)
	if err != nil {
		t.Error(err)
	}
	axt := x0meet1.Time(time.Now())
	axtp := &axt

	fmt.Println(xt.CreateTime.Time())
	fmt.Println(structs.FindMetaType(reflect.TypeOf(xt.CreateTime)).Kind())
	// fmt.Println(structs.FindMetaType())
	fmt.Println(structs.FindMetaType(reflect.TypeOf(time.Now())).Name())
	fmt.Println("------")
	fmt.Println(structs.FindMetaType(reflect.ValueOf(axtp).Type()).PkgPath())
	p := &Param{}
	str := `{"start":"2019-12-10T18:12:49","end":"2019-12-10T18:12:49"}`
	err = json.Unmarshal([]byte(str), p)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(p.Start.Time(), p.End.Time())
	fmt.Println(p, str)
	xat := structs.Time(time.Now())
	xatp := &xat
	fmt.Println(xatp.Time().Format("2006-01-02 15:04:05"))

	fmt.Println(xt.CreateTime)
}
