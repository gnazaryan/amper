package maps

func ContainsValue[K comparable, V comparable](input map[K]*V, v *V) bool {
	for _, value := range input {
		if value != nil && *value == *v {
			return true
		}
	}
	return false
}

func GetFirst(input *map[any]*any) *any {
	if input != nil {
		for _, value := range *input {
			return value
		}
	}
	return nil
}

func GetStringToInterfaceMap(input *map[string]string) *map[string]interface{} {
	if input != nil {
		result := make(map[string]interface{})
		for key, value := range *input {
			result[key] = value
		}
		return &result
	}
	return nil
}
