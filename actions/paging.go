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

// TODO: handle errors
func extractPageableFromRequest(c *gin.Context) *models.Pageable {
	pageFromQueryParam, err := strconv.Atoi(c.Query("page"))
	fmt.Println("Page from request", pageFromQueryParam)
	if err != nil {
		pageFromQueryParam = 0
	}

	sizeFromQueryParam, err := strconv.Atoi(c.Query("size"))
	fmt.Println("Size from request", sizeFromQueryParam)
	if err != nil {
		sizeFromQueryParam = 100
	}

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
