package helper

func FilterStruct[T any](data []T, f func(T) bool) []T {
	// CHECK ARRAY
	if len(data) < 1 {
		return data
	}

	var newData []T
	for _, d := range data {
		if f(d) {
			newData = append(newData, d)
		}
	}

	return newData
}
