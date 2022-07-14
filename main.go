package main

import (
	"net/http"
	"starryProject/controllers/households"
	"starryProject/daos"
	_ "starryProject/utils/db"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Service is Healthy!",
		})
	})

	householdsDAO := daos.NewHouseholdsDAO()
	membersDAO := daos.NewMembersDAO()

	households.NewHandler(householdsDAO, membersDAO).RouteGroup(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
