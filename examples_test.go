package zeroframework_test

import (
	"container/list"
	"container/ring"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/0meet1/zero-framework/structs"
)

// func TestOSSMiniV2(t *testing.T) {

// 	stagingAppId := "staging"
// 	stagingkey := "038AD75CD48825CA66EE8F10D561C4F9"

// 	keeper := ossminiv2.NewKeeper("www.qdeasydo.com", "frecstore").UseSSL().StagingSecret(stagingAppId, stagingkey)
// 	// imagebytes, err := os.ReadFile("/Users/bourbon/Downloads/12.jpeg")
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// ticket, err := keeper.Staging("12.jpeg", imagebytes)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println(ticket)

// 	// xP, err := keeper.Exchange(ticket)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println(xP)

// 	xuri, err := keeper.Complete("2024-03/018e63da-5d6f-7dd5-bb1d-e05c168e05b1")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(xuri)
// }

// func x2(c chan string) bool {
// 	select {
// 	case s := <-c:
// 		fmt.Println(s)
// 		return true

// 	case <-time.After(time.Duration(500) * time.Millisecond):
// 		fmt.Println("timeout ")
// 		return false
// 	}
// }

func TestConnectOracle(t *testing.T) {
	// database, err := gorm.Open(ora.Open("kangni/BSZnvPgL@158.158.5.57:1521/sapbmsprddb"), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	// dbPool, err := database.DB()
	// if err != nil {
	// 	panic(err)
	// }

	// dbPool.SetMaxIdleConns(10)
	// dbPool.SetMaxOpenConns(30)
	// dbPool.SetConnMaxLifetime(time.Second * time.Duration(100))

	// do somethings

	// fmt.Println(processors.ParseJSONColumnName("Xc1Feature.abc"))
	// var err error = &errdef.ZeroExceptionDef{}
	// fmt.Println(reflect.TypeOf(err) == reflect.TypeOf(errdef.ZeroExceptionDef{}))

	// c := make(chan string)

	// go func() {
	// 	for i := 0; i < 100; i++ {
	// 		go func() {
	// 			c <- fmt.Sprintf("%d", i)
	// 		}()
	// 	}
	// }()

	// for {
	// 	if !x2(c) {
	// 		close(c)
	// 		break
	// 	}
	// }

}

func TestXX(t *testing.T) {

	// jsonbytes, err := json.Marshal(1)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(x0meet1st.Md5Bytes(jsonbytes))
	// var i byte = 1.0
	v := reflect.ValueOf(structs.Time{})
	fmt.Println(v.Type().String())
	fmt.Println(reflect.Float64.String())
}

type TesT struct {
	Terc1 string
	Terc2 string
	Terc3 byte
	Terc4 float32
}

func TestRT(t *testing.T) {
	te := &TesT{
		Terc1: "xxx1",
		Terc2: "ccx",
	}

	fmt.Println(te)
	ptr1 := reflect.ValueOf(te)

	fmt.Println(ptr1.Kind())
	fmt.Println(ptr1.Kind() == reflect.Pointer)

	Terc1rf := ptr1.Elem().FieldByName("Terc1")
	fmt.Println(Terc1rf.Addr())
	fmt.Println(Terc1rf.String())
	Terc1rf.SetString("1231")

	Terc3rf := ptr1.Elem().FieldByName("Terc3")
	fmt.Println(Terc3rf.Addr())
	Terc3rf.SetInt(33)

	Terc4rf := ptr1.Elem().FieldByName("Terc4")
	fmt.Println(Terc4rf.Addr())
	Terc4rf.SetFloat(33.3333)

	fmt.Println(te)
}

func TestRT2(t *testing.T) {
	ti := structs.Time(time.Now())
	jsonbytes, err := json.Marshal(&ti)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonbytes))

	jt := `"2025-04-01T17:16:50"`
	fmt.Println([]byte(jt))
	ti2 := structs.Time{}
	err = json.Unmarshal([]byte(jt), &ti2)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.ValueOf(ti2).Type().String())
	fmt.Println(reflect.ValueOf(ti2).Type().String() == "structs.Time")
	fmt.Println(ti2.Time().Format("2006-01-02 15:04:05"))

	jt2 := `"2025-04-01T17:16:50Z"`
	fmt.Println([]byte(jt2))
	dat := time.Now()
	err = json.Unmarshal([]byte(jt2), &dat)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.ValueOf(dat).Type().String())
	fmt.Println(reflect.ValueOf(dat).Type().String() == "time.Time")
	fmt.Println(dat.Format("2006-01-02 15:04:05"))
}

func TestMap(t *testing.T) {

	te := &TesT{
		Terc1: "xxx1",
		Terc2: "ccx",
	}

	jsonbytes, err := json.Marshal(te)
	if err != nil {
		panic(err)
	}

	m := make(map[string]any)

	mc := reflect.New(reflect.ValueOf(m).Type()).Interface()
	fmt.Println(mc)
	err = json.Unmarshal(jsonbytes, &mc)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
	fmt.Println(reflect.ValueOf(mc).Elem())
}

func TestStruct(t *testing.T) {

	te := &TesT{
		Terc1: "xxx1",
		Terc2: "ccx",
	}

	jsonbytes, err := json.Marshal(te)
	if err != nil {
		panic(err)
	}

	mc := reflect.New(reflect.ValueOf(te).Elem().Type()).Interface()
	fmt.Println(mc)
	err = json.Unmarshal(jsonbytes, mc)
	if err != nil {
		panic(err)
	}
	fmt.Println(te)
	fmt.Println(reflect.ValueOf(mc).Elem())
}

type TesTChild struct {
	TesT

	TercCC3 string
}

func TestRTChild(t *testing.T) {
	te := &TesTChild{
		TesT: TesT{
			Terc1: "xxx1",
			Terc2: "ccx",
		},
	}
	fmt.Println(reflect.New(reflect.TypeOf(&TesTChild{})).Interface())
	fmt.Println(te)
	// ptr1 := reflect.ValueOf(te)
	ftype := reflect.TypeOf(TesTChild{})
	ptr1 := reflect.New(ftype)
	fmt.Println(ptr1.Interface())
	// for i := 0; i < ptr1.Elem().NumField(); i++ {
	// 	fmt.Println("111")
	// }

	fmt.Println(ptr1.Kind())
	fmt.Println(ptr1.Kind() == reflect.Pointer)

	Terc1rf := reflect.ValueOf(ptr1.Interface()).Elem().FieldByName("Terc1")
	fmt.Println(Terc1rf.Addr())
	fmt.Println(Terc1rf.String())
	Terc1rf.SetString("1231")

	// Terc3rf := ptr1.Elem().FieldByName("Terc3")
	// fmt.Println(Terc3rf.Addr())
	// Terc3rf.SetInt(33)

	Terc4rf := reflect.ValueOf(ptr1.Interface()).Elem().FieldByName("Terc4")
	fmt.Println(Terc4rf.Addr())
	Terc4rf.SetFloat(33.3333)

	ptr2 := reflect.ValueOf(ptr1.Interface())
	TercCC33rf := ptr2.Elem().FieldByName("TercCC3")
	fmt.Println(TercCC33rf.Addr())
	TercCC33rf.SetString("1231xxxx")

	fmt.Println(ptr1.Interface())
	fmt.Println(ptr2.Interface())
	fmt.Println(reflect.ValueOf(ptr1.Interface()).Interface())

	jsonbytes, _ := json.Marshal(reflect.ValueOf(ptr1.Interface()).Interface())
	fmt.Println(string(jsonbytes))

	fmt.Println(reflect.ValueOf(nil))
}

func TestRing(t *testing.T) {

	list.New()

	R_LEN := 10
	r := ring.New(10)
	for i := 0; i < R_LEN; i++ {
		r.Value = 1 + i
		r = r.Next()
	}

	exitch := make(chan int)

	go func() {
		for {
			fmt.Println("*******")
			fmt.Println(r.Prev().Value)
			fmt.Println(r.Value)
			r = r.Next()
			<-time.After(time.Duration(1) * time.Second)
		}
	}()

	select {
	case <-exitch:
		fmt.Println(" stoped ")
	case <-time.After(time.Duration(60) * time.Second):
		fmt.Println(" timeout ")
	}

	// l := list.New()
}
