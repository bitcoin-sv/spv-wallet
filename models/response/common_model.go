package response

import "time"

// Model is a common model that contains common fields for all models.
type Model struct {
	// CreatedAt is a time when outer model was created.
	CreatedAt time.Time `json:"createdAt" example:"2024-02-26T11:00:28.069911Z"`
	// UpdatedAt is a time when outer model was updated.
	UpdatedAt time.Time `json:"updatedAt" example:"2024-02-26T11:01:28.069911Z"`
	// DeletedAt is a time when outer model was deleted.
	DeletedAt *time.Time `json:"deletedAt" example:"2024-02-26T11:02:28.069911Z"`
	// Metadata is a metadata map of outer model.
	Metadata map[string]interface{} `json:"metadata" swaggertype:"object,string" example:"key:value,key2:value2"`
}
