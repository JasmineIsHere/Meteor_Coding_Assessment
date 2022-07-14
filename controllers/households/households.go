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
	rg.POST("/", h.create)
	rg.POST("/:householdID", h.addMember)
}

func (h *householdHandler) create(c *gin.Context) {
	newHousehold := domains.NewHouseholdReq{}
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

	newMember := domains.NewMemberReq{}
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
		h.membersDAO.UpdateSpouse(boil.GetDB(), spouse, member.ID)
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Member successfully added to household with id = %v", householdID)})
}
