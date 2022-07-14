package domains

import (
	"starryProject/models"

	"github.com/volatiletech/null/v8"
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

func HouseholdModelsToHouseholdRespAgeFilter(household models.Household, minDate null.Time, maxDate null.Time, isOverlap null.Bool) *HouseholdResp {
	memberSlice := household.R.Members

	var members []Member
	var toAdd bool
	for _, member := range memberSlice {

		if !minDate.IsZero() && !maxDate.IsZero() && (!isOverlap.IsZero() && isOverlap.Bool == true) && (member.Dob.After(minDate.Time) || member.Dob.Before(maxDate.Time)) {
			toAdd = false
		} else if !minDate.IsZero() && !maxDate.IsZero() && (!isOverlap.IsZero() && isOverlap.Bool == false) && (member.Dob.Before(minDate.Time) && member.Dob.After(maxDate.Time)) {
			toAdd = false
		} else if !minDate.IsZero() && maxDate.IsZero() && member.Dob.Before(minDate.Time) {
			toAdd = false
		} else if minDate.IsZero() && !maxDate.IsZero() && member.Dob.After(maxDate.Time) {
			toAdd = false
		} else {
			toAdd = true
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
