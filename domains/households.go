package domains

type Household struct {
	Type string `json:"type"`
}

type HouseholdResp struct {
	Type    string   `json:"type"`
	Members []Member `json:"family_members"`
}
