package domains

import (
	"time"

	"github.com/volatiletech/null/v8"
)

type NewMemberReq struct {
	Name           string    `json:"name"`
	Gender         string    `json:"gender"`
	MaritalStatus  string    `json:"marital_status"`
	SpouseID       null.Uint `json:"spouse_id"`
	OccupationType string    `json:"occupation_type"`
	AnnualIncome   float64   `json:"annual_income"`
	DOB            time.Time `json:"dob"`
}
