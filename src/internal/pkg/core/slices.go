package core

import (
	"errors"
	"reflect"
)

func GetValueAt(slice []string, index int) string {
	if len(slice) >= index+1 {
		return slice[index]
	}

	return ""
}

func Segment(a []interface{}, chunkSize int) ([][]interface{}, error) {
	if chunkSize < 1 {
		return nil, errors.New("chunkSize must be greater that zero")
	}
	chunks := make([][]interface{}, 0, (len(a)+chunkSize-1)/chunkSize)

	for chunkSize < len(a) {
		a, chunks = a[chunkSize:], append(chunks, a[0:chunkSize:chunkSize])
	}
	chunks = append(chunks, a)
	return chunks, nil
}

func SegmentByJsonByteLength(a []interface{}, maximumByteLength int, maximumSliceLength int) ([][]interface{}, error) {
	if a == nil || len(a) == 0 {
		return make([][]interface{}, 0, 0), nil
	}

	if maximumByteLength < 1 {
		return nil, errors.New("maximum byte length must be greater that zero")
	}
	if maximumSliceLength < 1 {
		return nil, errors.New("maximum slice length must be greater that zero")
	}

	result := make([][]interface{}, 0, 0)
	current := make([]interface{}, 0)
	currentSize := 0
	for _, item := range a {
		itemAsJson := MapToJson(item)
		itemLength := len(itemAsJson)

		if currentSize+itemLength < maximumByteLength && len(current) < maximumSliceLength {
			current = append(current, item)
			currentSize += itemLength
		} else {
			result = append(result, current)
			current = make([]interface{}, 0)

			current = append(current, item)
			currentSize += itemLength
		}
	}
	if len(current) != 0 {
		result = append(result, current)
	}

	return result, nil
}

func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
