package filter

// SearchContactsQuery is a model for handling searching with filters and metadata passed in the query string
type SearchContactsQuery = SearchParams[ContactFilter]

// SearchTransactionsQuery is a model for handling searching with filters and metadata passed in the query string
type SearchTransactionsQuery = SearchParams[TransactionFilter]

// SearchUtxosQuery is a model for handling searching with filters and metadata passed in the query string
type SearchUtxosQuery = SearchParams[UtxoFilter]

// SearchAccessKeysQuery is a model for handling searching with filters and metadata passed in the query string
type SearchAccessKeysQuery = SearchParams[AccessKeyFilter]
