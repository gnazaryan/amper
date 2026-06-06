package ampstrings

import (
	"amper/common/util"
	"container/list"
	"math/rand"
	"strconv"
	"strings"
)

var SEPERATOR = "_|~|_"

func HasValue(input *string) bool {
	if input != nil && len(*input) > 0 {
		return true
	}
	return false
}

func JoinInt64(input *[]int64, delimeter string) *string {
	var sb strings.Builder
	if input != nil {
		for index, value := range *input {
			if index > 0 {
				sb.WriteString(delimeter)
			}
			sb.WriteString(strconv.FormatInt(value, 10))
		}
	}
	return util.PointerString(sb.String())
}

func JoinListInt64(input *list.List, delimeter string) *string {
	var sb strings.Builder
	if input != nil {
		index := 0
		for item := input.Front(); item != nil; item = item.Next() {
			if index > 0 {
				sb.WriteString(delimeter)
			}
			sb.WriteString(strconv.FormatInt(*item.Value.(*int64), 10))
			index++
		}
	}
	return util.PointerString(sb.String())
}

func JoinStringBuilder(input *strings.Builder) string {
	return ""
}

func EmptyIfNil(input *string) string {
	if input != nil {
		return *input
	}
	return ""
}

func EmptyIfNilInt64(input *int64) string {
	if input != nil {
		return strconv.FormatInt(*input, 10)
	}
	return ""
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
