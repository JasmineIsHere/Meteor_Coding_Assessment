package daos

import (
	"starryProject/domains"
	"starryProject/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type HouseholdsDAO interface {
	AddHousehold(exec boil.Executor, householdDomain *domains.Household) error
	GetAll(exec boil.Executor) (*models.HouseholdSlice, error)
	GetByID(exec boil.Executor, householdID uint) (*models.Household, error)
	// GetSEB(cutoffDate time.Time, cutoffIncome int) (*models.Household, error)
}

type householdsDAO struct{}

func NewHouseholdsDAO() *householdsDAO {
	return &householdsDAO{}
}

func (dao *householdsDAO) AddHousehold(exec boil.Executor, householdDomain *domains.Household) error {
	household := &models.Household{
		Type: householdDomain.Type,
	}
	if err := household.Insert(exec, boil.Infer()); err != nil {
		return err
	}
	return nil
}

func (dao *householdsDAO) GetAll(exec boil.Executor) (*models.HouseholdSlice, error) {
	householdSlice, err := models.Households(
		qm.Load(models.HouseholdRels.Members)).All(exec)
	if err != nil {
		return nil, err
	}
	return &householdSlice, nil
}

func (dao *householdsDAO) GetByID(exec boil.Executor, householdID uint) (*models.Household, error) {
	household, err := models.Households(
		qm.Load(models.HouseholdRels.Members),
		models.HouseholdWhere.ID.EQ(householdID)).One(exec)
	if err != nil {
		return nil, err
	}
	return household, nil
}
