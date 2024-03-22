package zeroframework_test

import (
	"fmt"
	"testing"

	"github.com/0meet1/zero-framework/ossminiv2"
)

func TestOSSMiniV2(t *testing.T) {

	stagingAppId := "staging"
	stagingkey := "038AD75CD48825CA66EE8F10D561C4F9"

	keeper := ossminiv2.NewKeeper("www.qdeasydo.com", "frecstore").UseSSL().StagingSecret(stagingAppId, stagingkey)
	// imagebytes, err := os.ReadFile("/Users/bourbon/Downloads/12.jpeg")
	// if err != nil {
	// 	panic(err)
	// }

	// ticket, err := keeper.Staging("12.jpeg", imagebytes)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(ticket)

	// xP, err := keeper.Exchange(ticket)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(xP)

	xuri, err := keeper.Complete("2024-03/018e63da-5d6f-7dd5-bb1d-e05c168e05b1")
	if err != nil {
		panic(err)
	}
	fmt.Println(xuri)
}
