package tools

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {

	s := []string{"Na", "FF", "AG", "Eight", "Nm", "666", "id", "ccvvX"}

	if SliceContainKey(s, "AG") {
		fmt.Println("Key:AG OK")
	}

	if SliceContainKey(s, "SX") {
		fmt.Println("Key:SX 是不应该出现的")
	}

	if SliceContainKey(s, "FF") {
		fmt.Println("Key:FF 不是Key,不应该出现")
	}

	if SliceContains(s, "AG") {
		fmt.Println("---- :AG OK")
	}

	if SliceContains(s, "FF") {
		fmt.Println("---- :FF 是里面的一个项")
	}

	fmt.Println("............................")
	imap := make(map[string]string)

	imap["a"] = "Tian"
	imap["b"] = "Di"
	imap["c"] = "Xuan"
	imap["d"] = "Huang"
	imap["e"] = "Yu"
	imap["f"] = "Zhou"
	imap["g"] = "Hong"
	imap["h"] = "Huang"

	if MapContainsKey(imap, "a") {
		fmt.Println("---- :a OK")
	}

	if MapContainsKey(imap, "Yu") {
		fmt.Println("---- :Yu OK")
	}

	for k, v := range imap {
		fmt.Println("k:", k, " v:", v)
	}

	// if MapContainsKey(imap, "a") {
	// 	fmt.Println("Key:a存在")
	// }

	// if MapContainsKey(imap, "i") {
	// 	fmt.Println("Key:i存在")
	// }

	// if MapContainsValue(imap, "Xuan") {
	// 	fmt.Println("Value:Xuan存在")
	// }

	// if MapContainsValue(imap, "Han") {
	// 	fmt.Println("Value:Han存在")
	// }

}
