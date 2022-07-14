package domains

import (
	"starryProject/models"
)

type Household struct {
	Type string `json:"type"`
}

type HouseholdResp struct {
	Type    string   `json:"type"`
	Members []Member `json:"family_members"`
}

func HouseholdModelsToHouseholdResp(household models.Household) *HouseholdResp {
	memberSlice := household.R.Members

	var members []Member
	for _, member := range memberSlice {
		members = append(members, Member{
			Name:           member.Name,
			Gender:         member.Gender,
			MaritalStatus:  member.MaritalStatus,
			SpouseID:       member.SpouseID,
			OccupationType: member.OccupationType,
			AnnualIncome:   member.AnnualIncome,
			DOB:            member.Dob,
		})
	}

	return &HouseholdResp{
		Type:    household.Type,
		Members: members,
	}
}
