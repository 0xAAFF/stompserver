package stompserver

import (
	"fmt"
	"testing"
)

//
// go test -v frame_test.go frame.go header.go reader.go media_type_names.go command.go encode.go writer.go

func TestFrame(t *testing.T) {

	// fmt.Println(-0xff)
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(TimeUUID().String())
	// }
	//TestingT(t)

	messageFrameWhoami()
}

func messageFrameWhoami() {
	body := `{"statue":1,"Name":"stompserver"}`
	messageFrame, _ := NewMessageFrame("/broadcast/whoami", "90246c88-70a0-45dd-9cbf-cf948f56639d-1", "sub-0")
	messageFrame.SetBody(body)
	fmt.Println(messageFrame.Serialize())
}
