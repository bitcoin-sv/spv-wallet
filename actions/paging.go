package actions

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Paging(c *gin.Context) {
	pageable := extractPageableFromRequest(c)

	fmt.Printf("Pageable from request: %+v\n", pageable)

	// call service with pageable data

	c.JSON(http.StatusOK, gin.H{
		"pageable": "OK",
	})
}

func extractPageableFromRequest(c *gin.Context) *models.Pageable {
	pageFromQueryParam, _ := strconv.Atoi(c.Query("page"))
	fmt.Println("Page from request", pageFromQueryParam)

	sizeFromQueryParam, _ := strconv.Atoi(c.Query("size"))
	fmt.Println("Size from request", sizeFromQueryParam)

	sort := c.QueryArray("sort")
	fmt.Println("Sort from request", sort)
	return &models.Pageable{
		Page: pageFromQueryParam,
		Size: sizeFromQueryParam,
		Sort: models.Sort{},
	}
}
