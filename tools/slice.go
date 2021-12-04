package tools

// // slice中是否存在指定项 or map中是否存在指定key
// func SliceContain(slice interface{}, item interface{}) bool {
// 	switch reflect.TypeOf(slice).Kind() {
// 	case reflect.Slice, reflect.Array:
// 		{
// 			s := reflect.ValueOf(slice)
// 			for i := 0; i < s.Len(); i++ {
// 				if reflect.DeepEqual(item, s.Index(i).Interface()) {
// 					return true
// 				}
// 			}
// 		}
// 	case reflect.Map:
// 		{
// 			if reflect.ValueOf(slice).MapIndex(reflect.ValueOf(item)).IsValid() {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// // map或者slice中是否存在指定Key
// func SliceContainKey(slice interface{}, key interface{}) bool {
// 	switch reflect.TypeOf(slice).Kind() {
// 	case reflect.Slice, reflect.Array:
// 		{
// 			fmt.Println(slice, key)
// 			fmt.Println("reflect.Slice, reflect.Array")

// 			s := reflect.ValueOf(slice)
// 			fmt.Println(s)
// 			for i := 0; i < s.Len(); i = i + 2 {
// 				if reflect.DeepEqual(key, s.Index(i).Interface()) {
// 					fmt.Println("True")
// 					return true
// 				}
// 			}
// 			return false
// 		}
// 	case reflect.Map:
// 		{
// 			if reflect.ValueOf(slice).MapIndex(reflect.ValueOf(key)).IsValid() {
// 				return true
// 			}
// 			return false
// 		}
// 	}
// 	return false
// }
func SliceContains(slice []string, key string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == key {
			return true
		}
	}
	return false
}

// map或者slice中是否存在指定Key
func SliceContainKey(slice []string, key string) bool {
	for i := 0; i < len(slice); i = i + 2 {
		if slice[i] == key {
			return true
		}
	}
	return false
}

func MapContainsKey(imap map[string]string, key string) bool {
	for k := range imap {
		if k == key {
			return true
		}
	}
	return false
}

func MapVArrayContainsKey(imap map[string][]string, key string) bool {
	for k := range imap {
		if k == key {
			return true
		}
	}
	return false
}

func MapContainsValue(imap map[string]string, value string) bool {
	for k := range imap {
		if imap[k] == value {
			return true
		}
	}
	return false
}

func MapGetFirstKeyByValue(imap map[string]string, value string) string {
	for k := range imap {
		if imap[k] == value {
			return k
		}
	}
	return ""
}

func StringSliceRemoveReplica(slc []string) []string {
	/*
	   slice(string类型)元素去重
	*/
	result := make([]string, 0)
	tempMap := make(map[string]bool, len(slc))
	for _, e := range slc {
		if !tempMap[e] {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

func IntSliceRemoveReplica(slc []int) []int {
	/*
	   slice(int类型)元素去重
	*/
	result := make([]int, 0)
	tempMap := make(map[int]bool, len(slc))
	for _, e := range slc {
		if !tempMap[e] {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}
