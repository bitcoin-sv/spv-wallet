package filter

// SearchDestinations is a model for handling searching with filters and metadata
type SearchDestinations = SearchModel[DestinationFilter]

// CountDestinations is a model for handling counting filtered destinations
type CountDestinations = ConditionsModel[DestinationFilter]

// SearchContacts is a model for handling searching with filters and metadata
type SearchContacts = SearchModel[ContactFilter]

// SearchPaymails is a model for handling searching with filters and metadata
type SearchPaymails = SearchModel[AdminPaymailFilter]

// CountPaymails is a model for handling counting filtered paymails
type CountPaymails = ConditionsModel[AdminPaymailFilter]
