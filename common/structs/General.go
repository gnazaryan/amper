package structs

// Result is used to encapsulate the api controller response
type Result struct {
	Success       bool   `json:"success"`
	Error         string `json:"error"`
	Authenticated int    `json:"authenticated"` // 0 means authenticated, -1 means not authenticated
}

type ResultValue struct {
	Result
	Value *string `json:"value"`
}

// ValidationResult is used to send a boolean value for valid results
type ValidationResult struct {
	Result
	Valid bool `json:"valid"`
}

// Notification is used to hold generic properties of notification
type Notification struct {
	Headline       *string
	EmailContent   *string
	AmperResources *string
	AmperURL       *string
	Data           interface{}
}

type ObjectTypeField struct {
	ID                    *int64  `json:"id"`
	FieldId               *int64  `json:"fieldId"`
	ApiName               *string `json:"apiName"`
	Type                  *string `json:"type"`
	Label                 *string `json:"label"`
	CreatedBy             *int64  `json:"createdBy"`
	CreatedDate           *string `json:"createdDate"`
	ObjectTypeId          *int64  `json:"objectTypeId"`
	ObjectTypeApiName     *string `json:"objectTypeApiName"`
	ObjectTypeLabel       *string `json:"objectTypeLabel"`
	ObjectTypeExtendsTo   *int64  `json:"objectTypeExtendsTo"`
	ObjectTypeCreatedBy   *int64  `json:"objectTypeCreatedBy"`
	ObjectTypeCreatedDate *string `json:"objectTypecreatedDate"`
}

type ObjectType struct {
	Session
	ID               *int64            `json:"id"`
	ObjectId         *int64            `json:"objectId"`
	ApiName          *string           `json:"apiName"`
	Label            *string           `json:"label"`
	ExtendsTo        *int64            `json:"extendsTo"`
	ExtendsToLabel   *string           `json:"extendsToLabel"`
	CreatedBy        *int64            `json:"createdBy"`
	CreatedDate      *string           `json:"createdDate"`
	ObjectTypeFields []ObjectTypeField `json:"objectTypeFields"`
}

type ObjectTypeHierarchy struct {
	ID             *int64               `json:"id"`
	ObjectId       *int64               `json:"objectId"`
	ApiName        *string              `json:"apiName"`
	Label          *string              `json:"label"`
	ExtendsTo      *ObjectTypeHierarchy `json:"extendsTo"`
	ExtendsToLabel *string              `json:"extendsToLabel"`
	CreatedBy      *int64               `json:"createdBy"`
	CreatedDate    *string              `json:"createdDate"`
}

// Entities struct representing entities
type ObjectTypes struct {
	Result
	Data *[]ObjectType `json:"data"`
}

type Domain struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	ServerName string `json:"serverName"`
	Port       int    `json:"port"`
	Auth       string `json:"auth"`
}

type Imap struct {
	Domains []Domain `json:"domains"`
}

type Smtp struct {
	Domains []Domain `json:"domains"`
}

type Settings struct {
	RootDirectory   *string `json:"rootDirectory"`
	AdobeLicenseKey *string `json:"adobeLicenseKey"`
	Imap            Imap    `json:"imap"`
	Smtp            Smtp    `json:"smtp"`
}

type SettingsResult struct {
	Result
	Data *Settings `json:"data"`
}

// ProfileResult is used to send a boolean value for valid results
type ProfileConfigurationResult struct {
	Result
	Data map[string]interface{} `json:"data"`
}
