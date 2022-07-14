package households

import (
	"net/http"
	"starryProject/daos"
	"starryProject/domains"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Handler interface {
	RouterGroup(engine *gin.Engine)
}

type householdHandler struct {
	householdsDAO daos.HouseholdsDAO
}

func NewHandler(householdDAO daos.HouseholdsDAO) *householdHandler {
	return &householdHandler{
		householdDAO,
	}
}

func (h *householdHandler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/households")
	rg.POST("/", h.create)
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

	c.JSON(http.StatusOK, "Household successfully added")
}
