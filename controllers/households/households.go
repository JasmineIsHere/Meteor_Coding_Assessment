package households

import (
	"errors"
	"fmt"
	"net/http"
	"starryProject/daos"
	"starryProject/domains"
	"starryProject/enums/household_types"
	"starryProject/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
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

	rgGrants := r.Group("/grants")
	rgGrants.GET("/seb", h.seb)   // Student Encouragement Bonus
	rgGrants.GET("/mgs", h.mgs)   // Multi-Generation Scheme
	rgGrants.GET("/eb", h.eb)     // Elder Bonus
	rgGrants.GET("/bsg", h.bsg)   // Baby Sunshine Grant
	rgGrants.GET("/yolo", h.yolo) // YOLO GST Grant
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
		households = append(households, *domains.HouseholdModelsToHouseholdResp(*household))
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

	householdResp := domains.HouseholdModelsToHouseholdResp(*household)
	c.JSON(http.StatusOK, householdResp)
}

func (h *householdHandler) seb(c *gin.Context) {
	// ASSUMPTION: ELIGIBILITY = (at least one member whose occupationType = "Student" AND age > 16 years) AND total household income < 200,000
	// ASSUMPTION: a person's age depends on whether a person's birthday has passed
	// ASSUMPTION: Households incomes of less than $200,000 refers to the family's total annual income

	year, month, day := time.Now().Date()
	cutoffDate := time.Date(year-16, month, day, 0, 0, 0, 0, time.Local)

	cutoffIncome := 200000
	householdSlice, err := h.householdsDAO.GetSEB(boil.GetDB(), cutoffDate, cutoffIncome)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var households []domains.HouseholdResp
	for _, household := range *householdSlice {
		households = append(households, *domains.HouseholdModelsToHouseholdRespAgeFilter(*household, null.TimeFrom(cutoffDate), null.Time{}, null.Bool{}))
	}
	c.JSON(http.StatusOK, households)
}

func (h *householdHandler) mgs(c *gin.Context) {
	// ASSUMPTION: ELIGIBILITY = at least ONE member whose age is < 18 or > 55 which will make everyone in the household qualified
	// ASSUMPTION: Households incomes of less than $150,000 refers to the family's total annual income
	year, month, day := time.Now().Date()
	minDate := time.Date(year-18, month, day, 0, 0, 0, 0, time.Local)
	maxDate := time.Date(year-55, month, day, 0, 0, 0, 0, time.Local)

	cutoffIncome := 150000
	householdSlice, err := h.householdsDAO.GetMGS(boil.GetDB(), minDate, maxDate, cutoffIncome)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var households []domains.HouseholdResp
	for _, household := range *householdSlice {
		households = append(households, *domains.HouseholdModelsToHouseholdResp(*household))
	}
	c.JSON(http.StatusOK, households)
}

func (h *householdHandler) eb(c *gin.Context) {
	year, month, day := time.Now().Date()
	maxDate := time.Date(year-55, month, day, 0, 0, 0, 0, time.Local)
	householdSlice, err := h.householdsDAO.GetEB(boil.GetDB(), maxDate, household_types.HDB.String())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var households []domains.HouseholdResp
	for _, household := range *householdSlice {
		households = append(households, *domains.HouseholdModelsToHouseholdRespAgeFilter(*household, null.Time{}, null.TimeFrom(maxDate), null.Bool{}))
	}
	c.JSON(http.StatusOK, households)
}

func (h *householdHandler) bsg(c *gin.Context) {
	year, month, day := time.Now().Date()
	minDate := time.Date(year, month-8, day, 0, 0, 0, 0, time.Local)
	householdSlice, err := h.householdsDAO.GetBSG(boil.GetDB(), minDate)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var households []domains.HouseholdResp
	for _, household := range *householdSlice {
		households = append(households, *domains.HouseholdModelsToHouseholdRespAgeFilter(*household, null.TimeFrom(minDate), null.Time{}, null.Bool{}))
	}
	c.JSON(http.StatusOK, households)
}

func (h *householdHandler) yolo(c *gin.Context) {
	// ASSUMPTION: Households incomes of less than $100,000 refers to the family's total annual income
	cutoffIncome := 100000
	householdSlice, err := h.householdsDAO.GetYOLO(boil.GetDB(), cutoffIncome)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, c.Errors.Last())
		return
	}

	var households []domains.HouseholdResp
	for _, household := range *householdSlice {
		households = append(households, *domains.HouseholdModelsToHouseholdResp(*household))
	}
	c.JSON(http.StatusOK, households)
}
