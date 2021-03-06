package daos

import (
	"starryProject/domains"
	"starryProject/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type MembersDAO interface {
	AddMember(exec boil.Executor, memberDomain domains.Member, householdID uint) (*models.Member, error)
	GetByHouseholdID(exec boil.Executor, householdID uint) (*models.MemberSlice, error)
	GetByID(exec boil.Executor, memberID uint) (*models.Member, error)
	UpdateSpouse(exec boil.Executor, member *models.Member, spouseID uint) (int64, error)
}

type membersDAO struct{}

func NewMembersDAO() *membersDAO {
	return &membersDAO{}
}

func (dao *membersDAO) AddMember(exec boil.Executor, memberDomain domains.Member, householdID uint) (*models.Member, error) {
	member := &models.Member{
		Name:           memberDomain.Name,
		Gender:         memberDomain.Gender,
		MaritalStatus:  memberDomain.MaritalStatus,
		SpouseID:       memberDomain.SpouseID,
		OccupationType: memberDomain.OccupationType,
		AnnualIncome:   memberDomain.AnnualIncome,
		Dob:            memberDomain.DOB,
		HouseholdID:    householdID,
	}
	if err := member.Insert(exec, boil.Infer()); err != nil {
		return nil, err
	}
	return member, nil
}

func (dao *membersDAO) GetByHouseholdID(exec boil.Executor, householdID uint) (*models.MemberSlice, error) {
	members, err := models.Members(
		qm.Load(models.MemberRels.Household),
		qm.InnerJoin("household ON member.household_id = household.id"),
		models.HouseholdWhere.ID.EQ(householdID),
	).All(exec)

	if err != nil {
		return nil, err
	}

	return &members, err
}

func (dao *membersDAO) GetByID(exec boil.Executor, memberID uint) (*models.Member, error) {
	member, err := models.Members(models.MemberWhere.ID.EQ(memberID)).One(exec)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (dao *membersDAO) UpdateSpouse(exec boil.Executor, member *models.Member, spouseID uint) (int64, error) {
	member.SpouseID = null.UintFrom(spouseID)
	rowsAff, err := member.Update(exec, boil.Whitelist(models.MemberColumns.SpouseID))
	if err != nil {
		return 0, err
	}
	return rowsAff, err
}
