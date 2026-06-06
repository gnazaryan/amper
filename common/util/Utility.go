package util

import (
	"amper/common/constants"
	"amper/common/crypto"
	"errors"
	"log"
	"strconv"
	"strings"
)

// EmptyString checks if the provided input value is nil or empty string
func EmptyString(input *string) (result bool) {
	if input == nil || len(*input) == 0 {
		result = true
	}
	return
}

// NonEmpty checks and returns non nil or epty of 2 input strings priority to input1
func NonEmpty(input1 *string, input2 *string) (result *string) {
	if EmptyString(input1) {
		result = input2
	} else {
		result = input1
	}
	return
}

// UUID generate UUID as a unique identifier
func UUID() (result *string) {
	return PointerString(crypto.UUID())
}

func Loggify(err error) {
	if err != nil {
		log.Println(err.Error(), err)
	}
}

// PointerInt returns the address of the input int
func PointerInt(input int) (result *int) {
	return &input
}

// PointerInt returns the address of the input int
func PointerInt64(input int64) (result *int64) {
	return &input
}

// PointerString returns the address of the input string
func PointerBoolean(input bool) (result *bool) {
	return &input
}

// PointerString returns the address of the input string
func PointerString(input string) (result *string) {
	return &input
}

// StringPointer returns the value of the input string address
func StringPointer(input *string) (result string) {
	return *input
}

// PSA returns the address of the input array string
func PSA(input []string) (result *[]string) {
	return &input
}

// Apply is replacing all placeholders in the input string with the key values in the map
func Apply(input *string, values *map[string]string) error {
	if len(*values) > 0 {
		for k, v := range *values {
			placeholder := "{" + k + "}"
			beginIndex := strings.Index(*input, placeholder)
			if beginIndex > 0 {
				temp := strings.Replace(*input, placeholder, v, 100)
				input = &temp
			}
		}
	}
	return nil
}

// FlatArray convert a map to a flat array of key values
func FlatArray(values *map[string]string) (result []string) {
	if values != nil {
		for k, v := range *values {
			result = append(result, k)
			result = append(result, v)
		}
	}
	return
}

func ComputeIfApsent(input *map[interface{}][]interface{}, key *interface{}, value *interface{}) {
	temp := (*input)[key]
	if temp == nil {
		var initialize []interface{}
		initialize = append(initialize, value)
		(*input)[key] = initialize
	} else {
		temp = append(temp, value)
		(*input)[key] = temp
	}
}

func IfElse(condition bool, input1 interface{}, input2 interface{}) interface{} {
	if condition {
		return input1
	} else {
		return input2
	}
}

func IfElsePointer(condition bool, input1 *interface{}, input2 *interface{}) interface{} {
	if condition {
		return *input1
	} else {
		return *input2
	}
}

func IntToStr(input *int64) string {
	var result string
	if input != nil {
		result = strconv.FormatInt(*input, 10)
	}
	return result
}

// Append a suffix to Api name to make it unique
func AppendApiName(apiName string, append string) string {
	index := strings.LastIndex(apiName, constants.API_SUFFIX)
	return apiName[:index] + string(append) + apiName[index:]
}

func ParseInt64(input *string) *int64 {
	res, err := strconv.ParseInt(*input, 10, 64)
	if err != nil {
		Loggify(err)
		return nil
	}
	return &res
}

func I2Num(i interface{}) (int64, error) {

	switch v := i.(type) {
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	default:
		return 0, errors.New("type error")
	}
}
