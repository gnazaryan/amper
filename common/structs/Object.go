package structs

type Field struct {
	ID              *int64  `json:"id"`
	ApiName         *string `json:"apiName"`
	Label           *string `json:"label"`
	Type            *string `json:"type"`
	Status          *int8   `json:"status"`
	Required        *int8   `json:"required"`
	ObjectId        *int64  `json:"objectId"`
	CreatedBy       *int64  `json:"createdBy"`
	CreatedDate     *string `json:"createdDate"`
	TextLength      *int64  `json:"textLength"`
	ObjectReference *int64  `json:"objectReference"`
}

// Entities struct representing entities
type FieldsResult struct {
	Result
	Data *[]Field `json:"data"`
}

type Fields []Field

func (f Fields) ApiNameList() []string {
	var result []string
	for _, field := range f {
		result = append(result, *field.ApiName)
	}
	return result
}

type ObjectMetadata struct {
	Object      *Entity
	ObjectTypes *[]ObjectType
	Fields      *[]Field
}
