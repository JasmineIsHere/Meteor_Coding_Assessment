package domains

import (
	"starryProject/models"
	"time"
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

func HouseholdModelsToHouseholdRespAgeFilter(household models.Household, cutoffDate time.Time, inequality string) *HouseholdResp {
	memberSlice := household.R.Members

	var members []Member
	var toAdd bool
	for _, member := range memberSlice {

		if inequality == "<" && member.Dob.After(cutoffDate) {
			toAdd = true
		} else if inequality == "<=" && (member.Dob.After(cutoffDate) || member.Dob.Equal(cutoffDate)) {
			toAdd = true
		} else if inequality == ">" && member.Dob.Before(cutoffDate) {
			toAdd = true
		} else if inequality == ">=" && (member.Dob.Before(cutoffDate) || member.Dob.Equal(cutoffDate)) {
			toAdd = true
		} else {
			toAdd = false
		}

		if toAdd == true {
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
	}

	return &HouseholdResp{
		Type:    household.Type,
		Members: members,
	}
}
