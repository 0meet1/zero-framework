package zeroframework_test

import (
	"encoding/json"
	"fmt"
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

	fmt.Println(xt.CreateTime.Time())
}
