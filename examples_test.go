package zeroframework_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	x0meet1 "github.com/0meet1/zero-framework"
	"github.com/0meet1/zero-framework/structs"
)

type ZeroXsacTestPerson struct {
	x0meet1.ZeroXsacXhttpStructs

	IdCard  string                      `json:"idCard,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(32),NULL" xsackey:"unique" xapi:"身份证号,String"`
	Records []*ZeroXsacTestPersonRecord `json:"records,omitempty" xhttpopt:"OX" xsacchild:"Person" xapi:"记录列表,Array[ZeroXsacTestPersonRecord]"`
	// Records []*ZeroXsacTestPersonRecord `json:"records,omitempty" xhttpopt:"XXO" xsacref:"test_person_ref,person_id,record_id,inspect"`
}

func (person *ZeroXsacTestPerson) XhttpPath() string     { return "person" }
func (person *ZeroXsacTestPerson) XsacDbName() string    { return "test" }
func (person *ZeroXsacTestPerson) XsacTableName() string { return "test_person" }
func (person *ZeroXsacTestPerson) XsacApiName() string   { return "人员信息" }
func (person *ZeroXsacTestPerson) XsacApiEnums() []string {
	return structs.NewApiEnums("command指令类型", structs.ApiEnums(
		"49", "单基色",
		"50", "双基色",
		"51", "三基色"))
}

type ZeroXsacTestPersonRecord struct {
	x0meet1.ZeroXsacXhttpStructs

	// Person *ZeroXsacTestPerson `json:"person,omitempty" xhttpopt:"OX" xsacname:"person_id" xsacfield:"ID" xsacprop:"NO,VARCHAR(32),NULL" xsackey:"foreign,test_person,id"`
	Person   *ZeroXsacTestPerson `json:"person,omitempty" xhttpopt:"XXO" xsacfield:"-" xsacref:"test_person_ref,record_id,person_id" xapi:"关联人员,ZeroXsacTestPerson"`
	Contents string              `json:"contents,omitempty" xhttpopt:"OO" xsacprop:"NO,VARCHAR(256),NULL" xapi:"记录内容,String"`
}

func (personr *ZeroXsacTestPersonRecord) XhttpPath() string     { return "personrecord" }
func (personr *ZeroXsacTestPersonRecord) XsacDbName() string    { return "test" }
func (personr *ZeroXsacTestPersonRecord) XsacTableName() string { return "test_person_record" }
func (personr *ZeroXsacTestPersonRecord) XsacApiName() string   { return "人员记录" }

type Param struct {
	Start structs.Time `json:"start"`
	End   structs.Time `json:"end"`
}

func markdownTest(t *testing.T) {
	rows := make([]string, 0)
	rows = append(rows, structs.NewApiHeader("UgLi服务", "v202311"))

	rows = append(rows, structs.NewApiContentHeader("一、布控区域管理"))
	rows = append(rows, structs.NewApiDataMod("UgLiAIGround模型参数(布控区域)", structs.ApiDataMods(
		"id", "UUID", "唯一标识", "NO", "NO", "",
		"createTime", "DateTime", "创建时间", "NO", "NO", "`yyyy-MM-ddTHH:mm:ss`",
		"updateTime", "DateTime", "更新时间", "NO", "NO", "`yyyy-MM-ddTHH:mm:ss`",
		"features", "JSON", "特征", "YES", "YES", "",
		"deviceCode", "String", "设备编号", "YES", "NO", ""))...)
	rows = append(rows, structs.NewApiEnums("command指令类型", structs.ApiEnums(
		"49", "单基色",
		"50", "双基色",
		"51", "三基色"))...)

	req := &structs.ZeroCoreStructs{}
	req.InitDefault()
	reqbytes, _ := json.MarshalIndent(req, "", "\t")

	respmap := make(map[string]interface{})
	respmap["code"] = 200
	respmap["message"] = "success"
	jsonbytes, _ := json.MarshalIndent(respmap, "", "\t")

	rows = append(rows, structs.NewApiContentNOE("添加布控区域：aiground/add", "/zeroapi/v1/ugliserv/aiground/add",
		string(reqbytes), string(jsonbytes))...)
	rows = append(rows, structs.NewApiContentNOE("更新布控区域：aiground/up", "/zeroapi/v1/ugliserv/aiground/up",
		string(reqbytes), string(jsonbytes))...)
	rows = append(rows, structs.NewApiContentNOE("删除布控区域：aiground/rm", "/zeroapi/v1/ugliserv/aiground/rm",
		string(reqbytes), string(jsonbytes))...)
	rows = append(rows, structs.NewApiContentNOE("查询布控区域：aiground/fetch", "/zeroapi/v1/ugliserv/aiground/fetch",
		string(reqbytes), string(jsonbytes))...)

	mkd := structs.NewMarkdown(rows...)
	fmt.Println(mkd.String())

	file, err := os.Create("/Users/bourbon/Desktop/README.MD")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(mkd.String())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("成功写入文件！")

	file2, err := os.Create("/Users/bourbon/Desktop/READMEx.html")
	if err != nil {
		t.Fatal(err)
	}
	defer file2.Close()

	_, err = file2.WriteString(mkd.HTML())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("成功写入文件！")
}

func TestZeroCoreStructs(t *testing.T) {

	// fmt.Println(structs.YearDurationString(time.Now(), "2006-01-02 15:04:05"))
	// fmt.Println(structs.MonthDurationString(time.Now(), "2006-01-02 15:04:05"))
	// fmt.Println(structs.DayDurationString(time.Now(), "2006-01-02 15:04:05"))

	xt := &ZeroXsacTestPersonRecord{}
	xt.InitDefault()

	// jsonbytes, err := json.Marshal(xt)
	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Println(string(jsonbytes))

	// err = json.Unmarshal(jsonbytes, xt)
	// if err != nil {
	// 	t.Error(err)
	// }
	// axt := x0meet1.Time(time.Now())
	// axtp := &axt

	// fmt.Println(xt.CreateTime.Time())
	// fmt.Println(structs.FindMetaType(reflect.TypeOf(xt.CreateTime)).Kind())
	// fmt.Println(structs.FindMetaType())
	// fmt.Println(structs.FindMetaType(reflect.TypeOf(time.Now())).Name())
	// fmt.Println("------")
	// fmt.Println(structs.FindMetaType(reflect.ValueOf(axtp).Type()).PkgPath())
	// p := &Param{}
	// str := `{"start":"2019-12-10T18:12:49","end":"2019-12-10T18:12:49"}`
	// err = json.Unmarshal([]byte(str), p)
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(p.Start.Time(), p.End.Time())
	// fmt.Println(p, str)
	// xat := structs.Time(time.Now())
	// xatp := &xat
	// fmt.Println(xatp.Time().Format("2006-01-02 15:04:05"))
	// fmt.Println(xt.CreateTime)

	// rows := make([]string, 0)
	// rows = append(rows, structs.NewApiHeader("UgLi服务", "v202311"))
	// persion := &ZeroXsacTestPerson{}
	// persion.ThisDef(persion)
	// precord := &ZeroXsacTestPersonRecord{}
	// precord.ThisDef(precord)

	// rows = append(rows, persion.XsacApiExports("一、", "/zeroapi/v1/ugliserv")...)
	// rows = append(rows, precord.XsacApiExports("二、", "/zeroapi/v1/ugliserv")...)

	// mkd := structs.NewMarkdown(rows...)
	// fmt.Println(mkd.String())

	// file, err := os.Create("/Users/bourbon/Desktop/README.MD")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer file.Close()

	// _, err = file.WriteString(mkd.String())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println("成功写入文件！")

	// file2, err := os.Create("/Users/bourbon/Desktop/READMEx.html")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer file2.Close()

	// _, err = file2.WriteString(mkd.HTML())
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println("成功写入文件！")

	// fmt.Println(structs.NumberToChinese(892000843))

	met := reflect.ValueOf(reflect.ValueOf(xt).Elem().MethodByName("XhttpPath111"))
	fmt.Println(met.TryRecv())
}
