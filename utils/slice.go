package utils

// MergeSliceString 合并数组
func MergeSliceString(a, b []string) []string {
	var arr []string
	for _, i := range a {
		arr = append(arr, i)
	}
	for _, j := range b {
		arr = append(arr, j)
	}
	return arr
}

// UniqueSliceString 去重
func UniqueSliceString(m []string) []string {
	d := make([]string, 0)
	tempMap := make(map[string]bool, len(m))
	for _, v := range m { // 以值作为键名
		if tempMap[v] == false {
			tempMap[v] = true
			d = append(d, v)
		}
	}
	return d
}

// ExcludeSliceString 排除
func ExcludeSliceString(s, exclude []string) []string {
	d := make([]string, 0)
	for _, v := range s {
		if ok, _ := InArray(v, exclude); !ok {
			d = append(d, v)
		}
	}
	return d
}
