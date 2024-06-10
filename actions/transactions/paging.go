package transactions

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (a *Action) paging(c *gin.Context) {
	pageFromQueryParam, _ := strconv.Atoi(c.Query("pageFromQueryParam"))
	fmt.Println("Page from request", pageFromQueryParam)

	sizeFromQueryParam, _ := strconv.Atoi(c.Query("sizeFromQueryParam"))
	fmt.Println("Size from request", sizeFromQueryParam)

	sort := c.QueryArray("sort")
	fmt.Println("Sort from request", sort)

	pageable := models.Pageable{
		Page: pageFromQueryParam,
		Size: sizeFromQueryParam,
		Sort: models.Sort{},
	}

	fmt.Printf("%+v\n", pageable)

	// call service with pageable data
}
