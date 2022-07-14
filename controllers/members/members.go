package members

import (
	"starryProject/daos"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RouterGroup(engine *gin.Engine)
}

type memberHandler struct {
	membersDAO daos.MembersDAO
}

func NewHandler(memberDAO daos.MembersDAO) *memberHandler {
	return &memberHandler{
		memberDAO,
	}
}

func (h *memberHandler) RouteGroup(r *gin.Engine) {
}
