package stompserver

import (
	"fmt"
	"strings"
	"testing"
)

func TestStompUnit(t *testing.T) {
	text := "this is a text with \n .and it will \x00 contain \x00 more and more \x00"

	textArray := strings.Split(text, "\x00")
	for i := 0; i < len(textArray); i++ {
		fmt.Println(textArray[i])
		fmt.Printf("%d %X \n", i, textArray[i])

	}
}
