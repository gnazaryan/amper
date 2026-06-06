package structs

type Amper struct {
	Id              *int64  `json:"id"`
	Identifier      *int64  `json:"identifier"`
	Name            *string `json:"name"`
	Type            *string `json:"type"`
	Address         *string `json:"address"`
	Port            *string `json:"port"`
	State           *int    `json:"state"`
	StateUpdateDate *string `json:"stateUpdateDate"`
	Usage           *int64  `json:"usage"`
	Limit           *int64  `json:"limit"`
	Directory       *string `json:"directory"`
	Key             *string `json:"key"`
}

type AmperResults struct {
	Result
	Data []Amper `json:"data"`
}

type AmperResult struct {
	Result
	Data *Amper `json:"data"`
}
