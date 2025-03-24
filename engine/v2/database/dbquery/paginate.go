package dbquery

import (
	"context"
	"strings"

	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"gorm.io/gorm"
)

// PaginatedQuery is a generic function for getting paginated results from a database.
func PaginatedQuery[T any](ctx context.Context, page filter.Page, db *gorm.DB, scopes ...func(tx *gorm.DB) *gorm.DB) (*models.PagedResult[T], error) {
	PageWithDefaults(&page)
	model := models.PagedResult[T]{}
	var modelType T
	var totalElements int64
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&modelType).Scopes(scopes...)

		if err := query.
			Scopes(Paginate(page)).
			Find(&model.Content).Error; err != nil {
			return err
		}

		if err := query.
			Count(&totalElements).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get paginated result")
	}

	model.PageDescription.Number = page.Number
	model.PageDescription.Size = len(model.Content)
	model.PageDescription.TotalElements, err = conv.Int64ToInt(totalElements)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to convert total elements")
	}
	model.PageDescription.TotalPages = model.PageDescription.TotalElements / page.Size
	if model.PageDescription.TotalElements%page.Size > 0 {
		model.PageDescription.TotalPages++
	}
	return &model, nil
}

// Paginate is a Scope function that returns a function that paginates a database query.
func Paginate(page filter.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page.Number - 1) * page.Size
		return db.Order(page.SortBy + " " + page.Sort).Offset(offset).Limit(page.Size)
	}
}

// PageWithDefaults sets default values for a Page object (in place).
func PageWithDefaults(page *filter.Page) {
	if page.Number <= 0 {
		page.Number = 1
	}

	switch {
	case page.Size > 100:
		page.Size = 100
	case page.Size <= 0:
		page.Size = 20
	}

	page.SortBy = strings.ToLower(page.SortBy)
	if page.SortBy == "" {
		page.SortBy = "created_at"
	}

	if strings.ToLower(page.Sort) == "asc" {
		page.Sort = "ASC"
	} else {
		page.Sort = "DESC"
	}
}
