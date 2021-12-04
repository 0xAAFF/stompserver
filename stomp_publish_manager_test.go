package stompserver

import (
	"fmt"
	"strings"
	"testing"
)

func TestManagerGo(t *testing.T) {
	destination := "/root/main/item/"
	fmt.Println(destination[:1])
	fmt.Println(destination[1:])
	fmt.Println(destination[1:5])

	fmt.Println("........................................")

	index := strings.Index(destination[1:], "/")
	if index > -1 {

		fmt.Println(destination[:index+1])

	}
}
