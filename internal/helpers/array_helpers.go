package helpers

func RemoveDuplicatesStringsArr(arr []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range arr {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicatesInt64Arr(arr []int64) []int64 {
	allKeys := make(map[int64]bool)
	var list []int64
	for _, item := range arr {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
