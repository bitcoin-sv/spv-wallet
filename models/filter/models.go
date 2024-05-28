package filter

// SearchDestinations is a model for handling searching with filters and metadata
type SearchDestinations = SearchModel[DestinationFilter]

// CountDestinations is a model for handling counting filtered destinations
type CountDestinations = ConditionsModel[DestinationFilter]
