package main

func referenceify[T any](x T) *T {
	return &x
}

func dereferenceify[T any](x *T) T {
	return *x
}

func arrayToReference[T any](x []T) []*T {
	var result []*T
	for _, v := range x {
		result = append(result, &v)
	}
	return result
}

func referenceToArray[T any](x []*T) []T {
	var result []T
	for _, v := range x {
		result = append(result, *v)
	}
	return result
}