package request

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type AccessKeyConditions struct {
	// IncludeDeleted is a flag whether or not to include deleted items in the search results
	IncludeDeleted *bool `form:"includeDeleted,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	// CreatedRange specifies the time range when a record was created.
	CreatedRangeFrom *time.Time `form:"createdRangeFrom,omitempty" example:"2024-02-26T11:01:28Z"`
	CreatedRangeTo   *time.Time `form:"createdRangeTo,omitempty" example:"2024-02-26T11:01:28Z"`

	// UpdatedRange specifies the time range when a record was updated.
	UpdatedRangeFrom *time.Time `form:"updatedRangeFrom,omitempty" example:"2024-02-26T11:01:28Z"`
	UpdatedRangeTo   *time.Time `form:"updatedRangeTo,omitempty" example:"2024-02-26T11:01:28Z"`

	// RevokedRange specifies the time range when a record was revoked.
	RevokedRangeFrom *time.Time `form:"revokedRangeFrom,omitempty" example:"2024-02-26T11:01:28Z"`
	RevokedRangeTo   *time.Time `form:"revokedRangeTo,omitempty" example:"2024-02-26T11:01:28Z"`
}

func (conditions *AccessKeyConditions) MapToAccessKeyFilter() filter.AccessKeyFilter {
	createdRange := &filter.TimeRange{From: conditions.CreatedRangeFrom, To: conditions.CreatedRangeTo}
	updatedRange := &filter.TimeRange{From: conditions.UpdatedRangeFrom, To: conditions.UpdatedRangeTo}
	revokedRange := &filter.TimeRange{From: conditions.RevokedRangeFrom, To: conditions.RevokedRangeTo}

	modelFilter := filter.ModelFilter{
		IncludeDeleted: conditions.IncludeDeleted,
		CreatedRange:   createdRange,
		UpdatedRange:   updatedRange,
	}

	return filter.AccessKeyFilter{
		ModelFilter:  modelFilter,
		RevokedRange: revokedRange,
	}
}
