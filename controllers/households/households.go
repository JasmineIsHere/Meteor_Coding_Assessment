package households

import (
	"errors"
	"fmt"
	"net/http"
	"starryProject/daos"
	"starryProject/domains"
	"starryProject/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Handler interface {
	RouterGroup(engine *gin.Engine)
}

type householdHandler struct {
	householdsDAO daos.HouseholdsDAO
	membersDAO    daos.MembersDAO
}

func NewHandler(householdDAO daos.HouseholdsDAO, membersDAO daos.MembersDAO) *householdHandler {
	return &householdHandler{
		householdDAO,
		membersDAO,
	}
}

func (h *householdHandler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/households")
	rg.GET("/all", h.getAll)
	rg.GET("/:householdID", h.getByID)
	rg.POST("/", h.create)
	rg.POST("/:householdID", h.addMember)
}

func (h *householdHandler) create(c *gin.Context) {
	newHousehold := domains.Household{}
	if err := c.BindJSON(&newHousehold); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.Last())
		return
	}

	if err := h.householdsDAO.AddHousehold(boil.GetDB(), &newHousehold); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.Last())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Household successfully added"})
}

func (h *householdHandler) addMember(c *gin.Context) {
	householdIDUint64, _ := strconv.ParseUint(c.Param("householdID"), 10, 64)
	householdID := uint(householdIDUint64)

	newMember := domains.Member{}
	if err := c.BindJSON(&newMember); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.Last())
		return
	}

	household, err := h.householdsDAO.GetByID(boil.GetDB(), householdID)
	if err != nil {
		// no household with such householdID exists
		c.Error(errors.New(fmt.Sprintf("No household with householdID = %v exists", householdID)))
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var spouse *models.Member
	if !newMember.SpouseID.IsZero() {
		// check if spouseID exists
		if spouse, err = h.membersDAO.GetByID(boil.GetDB(), newMember.SpouseID.Uint); err != nil {
			c.Error(errors.New(fmt.Sprintf("No spouse with id = %v exists", newMember.SpouseID.Uint)))
			c.JSON(http.StatusNotFound, c.Errors.Last())
			return
		}
	}

	member, err := h.membersDAO.AddMember(boil.GetDB(), newMember, household.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.Last())
		return
	}

	// update spouse
	if spouse != nil {
		_, err := h.membersDAO.UpdateSpouse(boil.GetDB(), spouse, member.ID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Member successfully added to household with id = %v", householdID)})
}

func (h *householdHandler) getAll(c *gin.Context) {
	householdSlice, err := h.householdsDAO.GetAll(boil.GetDB())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var households []domains.HouseholdResp
	for _, household := range *householdSlice {
		var members []domains.Member
		memberSlice, err := h.membersDAO.GetByHouseholdID(boil.GetDB(), household.ID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusNotFound, c.Errors.Last())
			return
		}
		for _, member := range *memberSlice {
			members = append(members, domains.Member{
				Name:           member.Name,
				Gender:         member.Gender,
				MaritalStatus:  member.MaritalStatus,
				SpouseID:       member.SpouseID,
				OccupationType: member.OccupationType,
				AnnualIncome:   member.AnnualIncome,
				DOB:            member.Dob,
			})
		}
		households = append(households, domains.HouseholdResp{
			Type:    household.Type,
			Members: members,
		})
	}
	c.JSON(http.StatusOK, households)
}

func (h *householdHandler) getByID(c *gin.Context) {
	householdIDUint64, _ := strconv.ParseUint(c.Param("householdID"), 10, 64)
	householdID := uint(householdIDUint64)

	household, err := h.householdsDAO.GetByID(boil.GetDB(), householdID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	memberSlice, err := h.membersDAO.GetByHouseholdID(boil.GetDB(), householdID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}
	var members []domains.Member
	for _, member := range *memberSlice {
		members = append(members, domains.Member{
			Name:           member.Name,
			Gender:         member.Gender,
			MaritalStatus:  member.MaritalStatus,
			SpouseID:       member.SpouseID,
			OccupationType: member.OccupationType,
			AnnualIncome:   member.AnnualIncome,
			DOB:            member.Dob,
		})
	}

	householdResp := domains.HouseholdResp{
		Type:    household.Type,
		Members: members,
	}
	c.JSON(http.StatusOK, householdResp)
}
