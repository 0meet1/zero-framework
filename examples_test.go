package zeroframework_test

import (
	"fmt"
	"testing"
	"time"
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

func x2(c chan string) bool {
	select {
	case s := <-c:
		fmt.Println(s)
		return true

	case <-time.After(time.Duration(500) * time.Millisecond):
		fmt.Println("timeout ")
		return false
	}
}

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

	c := make(chan string)

	go func() {
		for i := 0; i < 100; i++ {
			go func() {
				c <- fmt.Sprintf("%d", i)
			}()
		}
	}()

	for {
		if !x2(c) {
			break
		}
	}

}
