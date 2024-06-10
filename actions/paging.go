package actions

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func Paging(c *gin.Context) {
	pageable := extractPageableFromRequest(c)

	fmt.Printf("Pageable from request: %+v\n", pageable)

	// call service with pageable data

	c.JSON(http.StatusOK, gin.H{
		"pageable": "OK",
	})
}

// TODO: handle default parameters values
// TODO: handle errors
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
		Sort: *createSortFromQueryParam(sort),
	}
}

func createSortFromQueryParam(sort []string) *models.Sort {
	orders := make([]models.Order, 0)
	const indexOfProperty = 0
	const indexOfDirection = 1
	for _, s := range sort {
		tokens := strings.Split(s, ",")
		orders = append(orders, models.Order{
			Property:  tokens[indexOfProperty],
			Direction: tokens[indexOfDirection],
		})
	}

	return &models.Sort{Orders: orders}
}
