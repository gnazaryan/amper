package argument

import (
	"amper/common/util"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// ValidateInt a validation tool to validate int parameters
func ValidateInt(input *int, name *string) error {
	if input == nil {
		return fmt.Errorf("the provided parameter: %s is nil, it must be supplied with a value", *name)
	}
	return nil
}

// ValidateString a validation tool to validate string parameters
func ValidateString(input *string, name *string) error {
	if util.EmptyString(input) {
		return fmt.Errorf("the provided parameter: %s is nil or empty, it must be supplied with a value", *name)
	}
	return nil
}

func ValidateApiName(apiName *string) bool {
	if apiName == nil ||
		!strings.HasSuffix(*apiName, "_amp") {
		return false
	}
	return true
}

// Validate is used to validate a set of parameters
func Validate(input map[string]interface{}) error {
	var result error
	var missingParams []string
	var k string
	var v interface{}
	for k, v = range input {
		switch v := v.(type) {
		case *int:
			var intValue = v
			if intValue == nil {
				missingParams = append(missingParams, k)
			}
		case *int64:
			var intValue = v
			if intValue == nil {
				missingParams = append(missingParams, k)
			}
		case *[]int64:
			var intValue = v
			if intValue == nil {
				missingParams = append(missingParams, k)
			}
		case *[]byte:
			var byteAArrayValue = v
			if byteAArrayValue == nil || len(*byteAArrayValue) < 1 {
				missingParams = append(missingParams, k)
			}
		case *string:
			var stingValue = v
			if util.EmptyString(stingValue) {
				missingParams = append(missingParams, k)
			}
		case *[]interface{}:
			var arrayValue = v
			if arrayValue == nil || len(*arrayValue) == 0 {
				missingParams = append(missingParams, k)
			}
		case *[]string:
			var arrayValue = v
			if arrayValue == nil || len(*arrayValue) == 0 {
				missingParams = append(missingParams, k)
			}
		case *[]map[string]interface{}:
			var arrayValue = v
			if arrayValue == nil || len(*arrayValue) == 0 {
				missingParams = append(missingParams, k)
			}
		default:
			var arrayValue = v
			if arrayValue == nil {
				missingParams = append(missingParams, k)
			}
		}
	}

	if len(missingParams) > 0 {
		var errorMessage bytes.Buffer
		errorMessage.WriteString("the")
		plural := len(missingParams) > 1
		if plural {
			errorMessage.WriteString(" parameters ")
		} else {
			errorMessage.WriteString(" parameter ")
		}
		errorMessage.WriteString(strings.Join(missingParams[:], ", "))
		if plural {
			errorMessage.WriteString(" are not provided, they must be supplied with values")
		} else {
			errorMessage.WriteString(" is not provided, it must be supplied with a value")
		}
		result = errors.New(errorMessage.String())
	}
	return result
}
