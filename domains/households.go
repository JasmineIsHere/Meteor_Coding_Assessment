package domains

type Household struct {
	Type string `json:"type"`
}

type HouseholdsResp struct {
	Type    string   `json:"type"`
	Members []Member `json:"family_members"`
}
