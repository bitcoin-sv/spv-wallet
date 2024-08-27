package swagger

import "time"

type AccessKeyParams struct {
	IncludeDeleted *bool `json:"includeDeleted,omitempty" swaggertype:"boolean" default:"false" example:"true"`

	CreatedRangeFrom *time.Time `json:"createdRange[from],omitempty" example:"2024-02-26T11:01:28Z"`
	CreatedRangeTo   *time.Time `json:"createdRange[to],omitempty" example:"2024-02-26T11:01:28Z"`

	UpdatedRangeFrom *time.Time `json:"updatedRange[from],omitempty" example:"2024-02-26T11:01:28Z"`
	UpdatedRangeTo   *time.Time `json:"updatedRange[to],omitempty" example:"2024-02-26T11:01:28Z"`

	RevokedRangeFrom *time.Time `json:"revokedRange[from],omitempty" example:"2024-02-26T11:01:28Z"`
	RevokedRangeTo   *time.Time `json:"revokedRange[to],omitempty" example:"2024-02-26T11:01:28Z"`
}
