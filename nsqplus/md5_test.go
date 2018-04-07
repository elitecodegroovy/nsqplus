package nsqplus

import (
	"testing"
	"fmt"
	"encoding/json"
)

/**

    {
    "platform":2,
    "point_code":"2110100",
    "created_time":1492341118994,
    "mobile_info":"iphone 6 plus",
    "app_version":"5.6.0",
    "source_from":"妈妈100App",
    "sign":"91571d62b8697b8b1c6ae19e6b03a97b"
    }
 */
func TestGetMD5Digest(t *testing.T) {
	hash_code := GetMD5Digest("1234567890")
	if hash_code == "e807f1fcf82d132f9bb018ca6738a19f" {
		fmt.Println("Success")
	} else {
		fmt.Println("result: ", hash_code)
		t.Fail()
	}
	events := Events{}
	event := &Event{}
	event.Platform = 2
	event.PointCode = "2110100"
	event.CreatedTime = 1492341118994
	event.AppVersion = "5.6.0"
	event.MobileInfo = "iphone 6 plus"
	event.SourceFrom = "妈妈100App"
	event.Sign = "91571d62b8697b8b1c6ae19e6b03a97b"
	events = append(events, event)
	js ,err := json.Marshal(events)
	if err != nil {
		t.Logf("%s", err.Error())
	}
	fmt.Println("json : \n ", fmt.Sprintf("%s", js))
}
