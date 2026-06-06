package structs

type Dashboard struct {
	ID            *int64    `json:"id"`
	Label         *string   `json:"label"`
	Description   *string   `json:"description"`
	Configuration *string   `json:"configuration"`
	Widgets       *[]Widget `json:"widgets"`
	CreatedBy     *int64    `json:"createdBy"`
	CreatedDate   *string   `json:"createdDate"`
}

type DashboardsResult struct {
	Result
	Data *[]Dashboard `json:"data"`
}

type Widget struct {
	ID            *int64  `json:"id"`
	Label         *string `json:"label"`
	Description   *string `json:"description"`
	Configuration *string `json:"configuration"`
	CreatedBy     *int64  `json:"createdBy"`
	CreatedDate   *string `json:"createdDate"`
}

type DashboardResult struct {
	Result
	Data *Dashboard `json:"data"`
}
