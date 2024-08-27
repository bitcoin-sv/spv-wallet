package swagger

type PageParams struct {
	Number int    `json:"page,omitempty" swaggertype:"integer" default:"1" example:"2"`
	Size   int    `json:"size,omitempty" swaggertype:"integer" default:"50" example:"5"`
	Order  string `json:"order,omitempty" swaggertype:"string" default:"desc" example:"desc"`
	SortBy string `json:"sortBy,omitempty" swaggertype:"string" default:"created_at" example:"id"`
}

type MetadataParams struct {
	Metadata []string `json:"metadata,omitempty" swaggertype:"array,string" example:"metadata[key]=value,metadata[key2]=value2"`
}
