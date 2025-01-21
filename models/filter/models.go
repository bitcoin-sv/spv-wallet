package filter

// SearchContacts is a model for handling searching with filters and metadata
type SearchContacts = SearchModel[ContactFilter]

// AdminSearchContacts is a model for handling searching with filters and metadata
type AdminSearchContacts = SearchModel[AdminContactFilter]

// SearchTransactions is a model for handling searching with filters and metadata
type SearchTransactions = SearchModel[TransactionFilter]

// CountTransactions is a model for handling counting filtered transactions
type CountTransactions = ConditionsModel[TransactionFilter]

// SearchXpubs is a model for handling searching with filters and metadata
type SearchXpubs = SearchModel[XpubFilter]

// CountXpubs is a model for handling counting filtered xPubs
type CountXpubs = ConditionsModel[XpubFilter]

// AdminSearchAccessKeys is a model for handling searching with filters and metadata
type AdminSearchAccessKeys = SearchModel[AdminAccessKeyFilter]

// AdminCountAccessKeys is a model for handling counting filtered transactions
type AdminCountAccessKeys = ConditionsModel[AdminAccessKeyFilter]

// SearchAccessKeys is a model for handling searching with filters and metadata
type SearchAccessKeys = SearchModel[AccessKeyFilter]

// CountAccessKeys is a model for handling counting filtered transactions
type CountAccessKeys = ConditionsModel[AccessKeyFilter]
