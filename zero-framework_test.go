package zeroframework_test

import (
	"fmt"
	"reflect"
	"testing"

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

	person := &ZeroXsacTestPerson{}
	person.InitDefault()
	person.ThisDef(person)
	fmt.Println(person.XsacDeclares())
	fmt.Println(person.XsacRefDeclares())

	personr := &ZeroXsacTestPersonRecord{
		Contents: "test",
	}
	personr.ThisDef(personr)
	// personr.XsacAutoConfig("test", "test_person_record")
	// fmt.Println(fmt.Sprintf("%b", 0b1111&0b100))
	// fmt.Println(reflect.ValueOf(personr).MethodByName("InitDefault").Call([]reflect.Value{}))
	fmt.Println(personr.XsacDeclares())
	fmt.Println(personr.XsacRefDeclares())

	fmt.Println(reflect.ValueOf(person).Elem().FieldByName("ID").Interface())

	fmt.Println(person.XsacFields())
	fmt.Println()
	fmt.Println(personr.XsacFields())

	fmt.Println(structs.FindMetaType(reflect.TypeOf(ZeroXsacTestPersonRecord{})))

	inx := reflect.New(structs.FindMetaType(reflect.TypeOf(&ZeroXsacTestPersonRecord{}))).Interface().(structs.ZeroXsacDeclares)
	fmt.Println(reflect.ValueOf(inx).Type())
}
