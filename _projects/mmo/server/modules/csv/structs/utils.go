package structs

func GetIndexFromStrings(strs []string, str string, defValue int) int {
	for i := 0; i < len(strs); i++ {
		if str == strs[i] {
			return i
		}
	}
	return defValue
}
