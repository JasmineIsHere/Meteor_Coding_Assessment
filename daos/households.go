package daos

import (
	"starryProject/domains"
	"starryProject/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type HouseholdsDAO interface {
	AddHousehold(db boil.Executor, householdDomain *domains.NewHouseholdReq) error
}

type householdsDAO struct{}

func NewHouseholdsDAO() *householdsDAO {
	return &householdsDAO{}
}

func (dao *householdsDAO) AddHousehold(db boil.Executor, householdDomain *domains.NewHouseholdReq) error {
	household := &models.Household{
		Type: householdDomain.Type,
	}
	err := household.Insert(db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
