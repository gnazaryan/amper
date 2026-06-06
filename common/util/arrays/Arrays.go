package arrays

import (
	"container/list"
	"strconv"
)

func Partition[T any](input []T, size int) [][]T {
	var result [][]T
	for i := 0; i < len(input); i += size {
		end := i + size
		if end > len(input) {
			end = len(input)
		}
		result = append(result, input[i:end])
	}

	return result
}

func Remove[T comparable](s []T, e T) []T {
	result := make([]T, 0)
	for _, v := range s {
		if v != e {
			result = append(result, v)
		}
	}
	return result
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func ContainsS(array []string, element string) bool {
	for _, item := range array {
		if item == element {
			return true
		}
	}
	return false
}

func IntToString(input *[]int64) []string {
	result := make([]string, len(*input))
	for i, v := range *input {
		result[i] = strconv.FormatInt(v, 10)
	}

	return result
}

func InterfaceToString(input *[]interface{}) ([]string, bool) {
	result := make([]string, len(*input))
	for i, v := range *input {
		strV, okV := v.(string)
		if okV {
			result[i] = strV
		} else {
			return nil, false
		}
	}
	return result, true
}

func ToArray(input *list.List) (result *[]any) {
	temp := make([]interface{}, input.Len())
	index := 0
	for e := input.Front(); e != nil; e = e.Next() {
		temp[index] = e
	}
	result = &temp
	return
}
