package daos

import (
	"starryProject/domains"
	"starryProject/enums/occupation_types"
	"starryProject/models"
	"time"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type HouseholdsDAO interface {
	AddHousehold(exec boil.Executor, householdDomain *domains.Household) error
	GetAll(exec boil.Executor) (*models.HouseholdSlice, error)
	GetByID(exec boil.Executor, householdID uint) (*models.Household, error)
	GetSEB(exec boil.Executor, youngerThan time.Time, cutoffIncome int) (*models.HouseholdSlice, error)
	GetMGS(exec boil.Executor, youngerThan time.Time, olderThan time.Time, cutoffIncome int) (*models.HouseholdSlice, error)
	GetEB(exec boil.Executor, olderThan time.Time, householdType string) (*models.HouseholdSlice, error)
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

func (dao *householdsDAO) GetSEB(exec boil.Executor, youngerThan time.Time, cutoffIncome int) (*models.HouseholdSlice, error) {
	households, err := models.Households(
		qm.Load(models.HouseholdRels.Members),
		qm.InnerJoin("member ON member.household_id = household.id"),
		models.MemberWhere.OccupationType.EQ(occupation_types.STUDENT.String()),
		models.MemberWhere.Dob.GT(youngerThan),
		qm.GroupBy("household.id"),
		qm.Having("SUM(member.annual_income) < ?", cutoffIncome),
	).All(exec)
	if err != nil {
		return nil, err
	}
	return &households, err
}

func (dao *householdsDAO) GetMGS(exec boil.Executor, youngerThan time.Time, olderThan time.Time, cutoffIncome int) (*models.HouseholdSlice, error) {
	households, err := models.Households(
		qm.Load(models.HouseholdRels.Members),
		qm.InnerJoin("member ON member.household_id = household.id"),
		models.MemberWhere.Dob.GT(youngerThan),
		qm.Or("member.dob < ?", olderThan),
		qm.GroupBy("household.id"),
		qm.Having("SUM(member.annual_income) < ?", cutoffIncome),
	).All(exec)
	if err != nil {
		return nil, err
	}
	return &households, err
}

func (dao *householdsDAO) GetEB(exec boil.Executor, olderThan time.Time, householdType string) (*models.HouseholdSlice, error) {
	households, err := models.Households(
		qm.Load(models.HouseholdRels.Members),
		qm.InnerJoin("member ON member.household_id = household.id"),
		models.MemberWhere.Dob.LTE(olderThan),
		models.HouseholdWhere.Type.EQ(householdType),
		qm.GroupBy("household.id"),
	).All(exec)
	if err != nil {
		return nil, err
	}
	return &households, err
}
